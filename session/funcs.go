package session

import (
	"github.com/gofiber/fiber/v2"
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
