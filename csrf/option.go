package csrf

import "github.com/gofiber/fiber/v2"

// Options defines a function type for configuring CSRF Option.
type Options func(*Option)

// Option holds the configuration options for CSRF middleware.
type Option struct {
	header bool
	key    string
	fail   fiber.Handler
	next   func(*fiber.Ctx) bool
}

// WithFail sets a custom failure handler for CSRF validation.
func WithFail(handler fiber.Handler) Options {
	return func(c *Option) {
		c.fail = handler
	}
}

// WithNext sets a custom function can be used to skip CSRF validation for certain requests.
func WithNext(handler func(*fiber.Ctx) bool) Options {
	return func(c *Option) {
		c.next = handler
	}
}

// WithHeader configures the CSRF middleware to check CSRF token from header.
func WithHeader(name string) Options {
	return func(c *Option) {
		if name != "" {
			c.header = true
			c.key = name
		}
	}
}

// WithForm configures the CSRF middleware to check CSRF token from form field.
func WithForm(name string) Options {
	return func(c *Option) {
		if name != "" {
			c.header = false
			c.key = name
		}
	}
}
