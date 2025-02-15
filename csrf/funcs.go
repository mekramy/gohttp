package csrf

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mekramy/gohttp/session"
)

// GetToken retrieves the token from the session associated with the given Fiber context.
// Returns an empty string if the session is nil or the "csrf" value is not found.
func GetToken(c *fiber.Ctx) string {
	s := session.Parse(c)
	if s == nil {
		return ""
	}
	return s.Cast("csrf").StringSafe("")
}

// RefreshToken generates a new CSRF token and saves it to the session.
// It returns the generated token or an error if the session cannot be resolved.
func RefreshToken(c *fiber.Ctx) (string, error) {
	// Parse session
	s := session.Parse(c)
	if s == nil {
		return "", errors.New("failed to resolve session")
	}

	// Save to session
	return refresh(s), nil
}

// refresh csrf on session
func refresh(s session.Session) string {
	token := uuid.NewString()
	s.Set("csrf", token)
	return token
}
