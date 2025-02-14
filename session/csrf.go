package session

import (
	"errors"
	"slices"

	"github.com/gofiber/fiber/v2"
)

// NewCSRF creates a new CSRF middleware handler for Fiber. It validates the CSRF token
// for incoming requests and generates a new token if needed. The secure argument specifies
// whether the CSRF cookie should be generated securely. The onFail function is called if the
// CSRF validation fails, and the next function can be used to skip CSRF validation for certain
// requests. By default, this middleware generates a 419 HTTP response if CSRF validation fails.
// This middleware must be called after the session middleware.
func NewCSRF(secure bool, onFail fiber.Handler, next func(*fiber.Ctx) bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip
		if next != nil && next(c) {
			return c.Next()
		}

		// Set allowed headers
		c.Append("Access-Control-Expose-Headers", "X-CSRF-TOKEN")
		c.Append("Access-Control-Allow-Headers", "X-CSRF-TOKEN")

		// Parse session
		session := Parse(c)
		if session == nil {
			return errors.New("failed to resolve session")
		}

		// Get CSRF token
		token := session.Cast("csrf").StringSafe("")

		// Generate or refresh token if needed
		if token == "" {
			token, _ = RefreshCSRF(secure, c, session)
		}

		// Ignore RFC9110 (GET, HEAD, OPTIONS, and TRACE) methods
		if slices.Contains([]string{fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete}, c.Method()) {
			// Parse request csrf token from X-CSRF-Token header
			input := c.Get("X-CSRF-Token")

			// Parse request csrf token from csrf_token cookie
			if input == "" {
				input = c.Cookies("csrf_token")
			}

			// Parse request csrf token from csrf_token form
			if input == "" {
				type Form struct {
					Token string `json:"csrf_token" form:"csrf_token"`
				}
				var inp Form
				c.BodyParser(&inp)
				input = inp.Token
			}

			// Validate CSRF
			if token == "" || input != token {
				if onFail != nil {
					return onFail(c)
				}
				return c.Status(419).SendString("invalid csrf token")
			}
		}

		return c.Next()
	}
}
