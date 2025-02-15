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

// RefreshCSRF generates a new CSRF token and saves it to the session.
// It returns the generated token or an error if the session not passed
// and cannot be resolved.
func RefreshCSRF(c *fiber.Ctx, s Session) (string, error) {
	// Parse session if not passed
	if s == nil {
		s = Parse(c)
		if s == nil {
			return "", errors.New("failed to resolve session")
		}
	}

	// Save to session
	token := uuid.NewString()
	s.Set("csrf", token)
	return token, nil
}
