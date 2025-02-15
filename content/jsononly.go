package content

import (
	"github.com/gofiber/fiber/v2"
)

// JsonOnly is a middleware that ensures the request's Content-Type is "application/json".
// If the Content-Type is not "application/json", it will execute the optional onFail handler
// if provided, or return a 406 Not Acceptable status by default.
func JsonOnly(onFail ...fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("Content-Type") != "application/json" {
			if len(onFail) > 0 && onFail[0] != nil {
				return onFail[0](c)
			}
			return c.SendStatus(fiber.StatusNotAcceptable)
		}
		return c.Next()
	}
}
