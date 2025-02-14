package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mekramy/gocache"
)

// NewMiddleware creates a new session middleware for the Fiber framework.
// It initializes a session using the provided cache and options, sets the necessary headers,
// stores the session in the context, and ensures the session is saved after the request is processed.
func NewMiddleware(cache gocache.Cache, options ...Options) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Create session
		s, err := New(ctx, cache, options...)
		if err != nil {
			return err
		}

		// Set Allowed header
		if s.isHeader() {
			ctx.Append("Access-Control-Expose-Headers", s.getName())
			ctx.Append("Access-Control-Allow-Headers", s.getName())
		}

		// Store to context
		ctx.Locals("SESSION", s)

		// Continue and save session
		err = ctx.Next()
		if err == nil {
			err = s.Save()
		}
		return err
	}
}
