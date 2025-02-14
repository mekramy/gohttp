package session

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mekramy/gocache"
	"github.com/mekramy/gocast"
)

// Session represents a user session interface with methods to manage session data.
type Session interface {
	// Id returns the session identifier.
	Id() string

	// Context returns the associated Fiber context.
	Context() *fiber.Ctx

	// Set stores a value in the session for the given key.
	Set(key string, value any)

	// Get retrieves a value from the session for the given key.
	Get(key string) any

	// Delete removes a value from the session for the given key.
	Delete(key string)

	// Exists checks if a key exists in the session.
	Exists(key string) bool

	// Cast returns a Caster for the value associated with the given key.
	Cast(key string) gocast.Caster

	// CreatedAt retrieves session creation date.
	CreatedAt() *time.Time

	// AddTTL extends the session's time-to-live.
	AddTTL(ttl time.Duration)

	// SetTTL set session's time-to-live.
	SetTTL(ttl time.Duration)

	// Destroy terminates the session.
	Destroy() error

	// Save persists the session data to storage if changed.
	// Must be called at the end of middleware.
	Save() error

	// Fresh generates a new session.
	Fresh() error

	// Load retrieves session data from storage.
	// Returns false if the session does not exist.
	Load() (bool, error)

	isHeader() bool
	getName() string
}

// New create or parse session driver.
func New(ctx *fiber.Ctx, cache gocache.Cache, options ...Options) (Session, error) {
	// Generate option
	option := &Option{
		ttl:       24 * time.Hour,
		name:      "session",
		header:    false,
		cookie:    &fiber.Cookie{},
		generator: UUIDGenerator,
	}
	for _, opt := range options {
		opt(option)
	}

	// Get session id
	var id string
	if option.header {
		id = ctx.Get(option.name)
	} else {
		id = ctx.Cookies(option.name)
	}

	// Generate session
	session := &session{
		id:    id,
		opt:   *option,
		ttl:   0,
		ctx:   ctx,
		cache: cache,
		data:  make(map[string]any),
	}

	ok, err := session.Load()
	if err != nil {
		return nil, err
	}

	if !ok {
		err := session.Fresh()
		if err != nil {
			return nil, err
		}
	}

	return session, nil
}
