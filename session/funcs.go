package session

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Parse extracts the Session object from the fiber.Ctx context.
// If the session data is found and is of the correct type, it returns the Session object.
// Otherwise, it returns nil.
func Parse(ctx *fiber.Ctx) Session {
	session, ok := ctx.Locals("SESSION").(Session)
	if ok {
		return session
	}

	return nil
}

// CSRF retrieves the CSRF token from the session associated with the given Fiber context.
// Returns an empty string if the session is nil or the "csrf" value is not found.
func CSRF(ctx *fiber.Ctx) string {
	session := Parse(ctx)
	if session != nil {
		session.Cast("csrf").StringSafe("")
	}
	return ""
}

// RefreshCSRF generates a new CSRF token, sets it in the response headers and cookies,
// and saves it to the session. It returns the generated token or an error if the session
// not passed and cannot be resolved.
func RefreshCSRF(secure bool, c *fiber.Ctx, s Session) (string, error) {
	token := uuid.NewString()

	// Parse session if not passed
	if s == nil {
		s = Parse(c)
		if s == nil {
			return "", errors.New("failed to resolve session")
		}
	}

	// Set headers
	c.Append("Access-Control-Expose-Headers", "X-CSRF-TOKEN")
	c.Append("Access-Control-Allow-Headers", "X-CSRF-TOKEN")
	c.Set("X-CSRF-TOKEN", token)

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "csrf_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	// Save to session
	s.Set("csrf", token)

	return token, nil
}
