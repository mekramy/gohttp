package session

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mekramy/gocache"
	"github.com/mekramy/gocast"
)

// session represents a user session with associated data and metadata.
type session struct {
	id   string         // Unique identifier for the session.
	opt  Option         // Configuration options for the session.
	data map[string]any // Key-value store for session data.

	ttl      time.Duration // Additional time-to-live for the session.
	fresh    bool          // Flag indicating if session is fresh.
	modified bool          // Flag indicating if session data has been modified.

	ctx   *fiber.Ctx    // Fiber context associated with the session.
	cache gocache.Cache // Cache for storing session data.
	mutex sync.RWMutex  // Mutex for synchronizing access to session data.
}

func (s *session) Id() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.id
}

func (s *session) Context() *fiber.Ctx {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.ctx
}

func (s *session) Set(k string, v any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if k = strings.TrimSpace(k); k != "" {
		s.data[k] = v
		s.modified = true
	}
}

func (s *session) Get(k string) any {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.data[k]
}

func (s *session) Delete(k string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, k)
	s.modified = true
}

func (s *session) Exists(k string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.data[k]
	return ok
}

func (s *session) Cast(k string) gocast.Caster {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return gocast.NewCaster(s.data[k])
}

func (s *session) CreatedAt() *time.Time {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	raw, ok := s.data["created_at"].(string)
	if !ok {
		return nil
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil
	}

	return &t
}

func (s *session) AddTTL(t time.Duration) {
	// Skip empty ttl
	if t <= 0 {
		return
	}

	// Safe race condition
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Schedule update
	s.ttl = t
	s.modified = true
	s.sync()
}

func (s *session) SetTTL(t time.Duration) {
	// Skip empty ttl
	if t <= 0 {
		return
	}

	// Safe race condition
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Schedule update
	s.ttl = -t
	s.modified = true
	s.sync()
}

func (s *session) Destroy() error {
	// Skip empty session
	if s.id == "" {
		return nil
	}

	// Safe race condition
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Delete from cache
	err := s.cache.Forget(s.k())
	if err != nil {
		return err
	}

	// Clear data
	s.id = ""
	s.data = make(map[string]any)
	return nil
}

func (s *session) Save() error {
	// Skip un-initialized or unchanged or destroyed session
	if s.id == "" || (!s.fresh && !s.modified) {
		return nil
	}

	// Safe race condition
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Encode data
	encoded, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	// Store New
	if s.fresh {
		return s.cache.Put(s.k(), encoded, &s.opt.ttl)
	}

	// Add ttl
	if s.ttl > 0 {
		ttl, err := s.cache.TTL(s.k())
		if err != nil {
			return err
		} else if ttl <= 0 {
			ttl = s.ttl
		} else {
			ttl += s.ttl
		}
		return s.cache.Put(s.k(), encoded, &ttl)
	}

	// Set ttl
	if s.ttl < 0 {
		ttl := -s.ttl
		return s.cache.Put(s.k(), encoded, &ttl)
	}

	// Save data
	_, err = s.cache.Set(s.k(), encoded)
	return err
}

func (s *session) Fresh() error {
	// Safe race condition
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Destroy old session
	if s.id != "" {
		err := s.cache.Forget(s.k())
		if err != nil {
			return err
		}
	}

	// Set identifier and created at
	s.id = s.opt.generator()
	s.ttl = s.opt.ttl
	s.data = make(map[string]any)
	s.ttl = 0
	s.fresh = true
	s.modified = true
	s.data["created_at"] = time.Now().Format(time.RFC3339)
	s.sync()

	return nil
}

func (s *session) Load() (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Not generated or empty id
	if s.id == "" {
		return false, nil
	}

	// Check if session exists
	exists, err := s.cache.Exists(s.k())
	if err != nil {
		return false, err
	} else if !exists {
		return false, nil
	}

	// Parse data and decode data
	caster, err := s.cache.Cast(s.k())
	if err != nil {
		return false, err
	}

	encoded, err := caster.String()
	if err != nil {
		return false, err
	}

	s.data = make(map[string]any)
	err = json.Unmarshal([]byte(encoded), &s.data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *session) isHeader() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.opt.header
}

func (s *session) getName() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.opt.name
}

func (s *session) sync() {
	// Ignore empty or destroyed
	if s.id == "" {
		return
	}

	// Send header data
	if s.opt.header {
		s.ctx.Set(s.opt.name, s.id)
		return
	}

	// Send cookie
	s.ctx.Cookie(&fiber.Cookie{
		Name:        s.opt.name,
		Value:       s.id,
		Expires:     time.Now().Add(s.ttl),
		Secure:      s.opt.cookie.Secure,
		Domain:      s.opt.cookie.Domain,
		SameSite:    s.opt.cookie.SameSite,
		Path:        s.opt.cookie.Path,
		MaxAge:      s.opt.cookie.MaxAge,
		HTTPOnly:    s.opt.cookie.HTTPOnly,
		SessionOnly: s.opt.cookie.SessionOnly,
	})
}

func (s *session) k() string {
	return "ses-" + s.id
}
