package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// FormOrURLEncodedOnly is a middleware that ensures the request's Content-Type is either
// "application/x-www-form-urlencoded" or "multipart/form-data".
// If the Content-Type is neither of these, it will execute the optional onFail handler
// if provided, or return a 406 Not Acceptable status by default.
func FormOrURLEncodedOnly(onFail ...fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentType := c.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" && contentType != "multipart/form-data" {
			if len(onFail) > 0 && onFail[0] != nil {
				return onFail[0](c)
			}
			return c.SendStatus(fiber.StatusNotAcceptable)
		}
		return c.Next()
	}
}
