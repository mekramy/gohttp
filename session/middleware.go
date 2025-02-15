package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mekramy/gocache"
)

// NewMiddleware creates a new session middleware for the Fiber framework.
// It initializes a session using the provided cache and options, sets the necessary headers,
// stores the session in the context, and ensures the session is saved after the request is processed.
func NewMiddleware(cache gocache.Cache, options ...Options) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create session
		s, err := New(c, cache, options...)
		if err != nil {
			return err
		}

		// Set Allowed header
		if s.isHeader() {
			c.Append("Access-Control-Expose-Headers", s.getName())
			c.Append("Access-Control-Allow-Headers", s.getName())
		}

		// Store to context
		c.Locals("SESSION", s)

		// Continue and save session
		err = c.Next()
		if err == nil {
			err = s.Save()
		}
		return err
	}
}
