package request

import (
	"github.com/gofiber/fiber/v2"
)

// XMLOnly is a middleware that ensures the request's Content-Type is "application/xml" or "text/xml".
// If the Content-Type is not "application/xml" or "text/xml", it will execute the optional onFail handler
// if provided, or return a 406 Not Acceptable status by default.
func XMLOnly(onFail ...fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentType := c.Get("Content-Type")
		if contentType != "application/xml" && contentType != "text/xml" {
			if len(onFail) > 0 && onFail[0] != nil {
				return onFail[0](c)
			}
			return c.SendStatus(fiber.StatusNotAcceptable)
		}
		return c.Next()
	}
}
