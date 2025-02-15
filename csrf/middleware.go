package csrf

import (
	"errors"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/mekramy/gohttp/session"
)

// NewMiddleware creates a new CSRF middleware handler with the provided options.
// It validates the CSRF token for incoming requests and generates a new token if needed.
// By default, this middleware generates a 419 HTTP response if CSRF validation fails.
//
// This middleware must be called after the session middleware.
func NewMiddleware(options ...Options) fiber.Handler {
	// Generate option
	option := &Option{
		header: false,
		key:    "csrf_token",
		fail:   nil,
		next:   nil,
	}
	for _, opt := range options {
		opt(option)
	}

	return func(c *fiber.Ctx) error {
		// Skip
		if option.next != nil && option.next(c) {
			return c.Next()
		}

		// Parse and generate token
		session := session.Parse(c)
		if session == nil {
			return errors.New("failed to resolve session")
		}
		token := session.Cast("csrf").StringSafe("")
		if token == "" { // Generate or refresh token if needed
			token = refresh(session)
		}

		// Proccess request
		// RFC9110 (GET, HEAD, OPTIONS, and TRACE) methods
		shouldCheck := slices.Contains(
			[]string{fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete},
			c.Method(),
		)
		if option.header {
			c.Append("Access-Control-Allow-Headers", "X-CSRF-TOKEN")
			if shouldCheck {
				input := c.Get("X-CSRF-Token")
				if token == "" || input != token {
					if option.fail != nil {
						return option.fail(c)
					}
					return c.Status(419).SendString("invalid csrf token")
				}
			}
		} else {
			if shouldCheck {
				type Form struct {
					Token string `json:"csrf_token" form:"csrf_token"`
				}
				var inp Form
				c.BodyParser(&inp)

				if token == "" || inp.Token != token {
					if option.fail != nil {
						return option.fail(c)
					}
					return c.Status(419).SendString("invalid csrf token")
				}
			}
		}

		return c.Next()
	}
}
