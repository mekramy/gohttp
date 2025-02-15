package limiter

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Options defines a function type for configuring Rate Limiter Option.
type Options func(*Option)

// Option holds the configuration options for Rate Limiter middleware.
type Option struct {
	key      string // Cache key
	attempts uint   // Maximum number of attempts allowed
	ttl      time.Duration
	skipFail bool
	fail     func(time.Duration) fiber.Handler // Custom failure handler
	next     func(*fiber.Ctx) bool             // Function to skip rate limiter validation for certain requests
	keys     func(*fiber.Ctx) []string         // Function to generate extra key based on request
}

// WithMaxAttempts sets the maximum number of attempts allowed.
func WithMaxAttempts(attempts uint) Options {
	return func(c *Option) {
		if attempts > 0 {
			c.attempts = attempts
		}
	}
}

// WithTTl sets the time-to-live for the rate limiter.
func WithTTl(ttl time.Duration) Options {
	return func(c *Option) {
		if ttl > 0 {
			c.ttl = ttl
		}
	}
}

// WithSkipFail sets the option to skip limiter if request has error.
func WithSkipFail(skipFail bool) Options {
	return func(c *Option) {
		c.skipFail = skipFail
	}
}

// WithFail sets a custom failure handler for Rate Limiter validation.
func WithFail(handler func(until time.Duration) fiber.Handler) Options {
	return func(c *Option) {
		c.fail = handler
	}
}

// WithNext sets a custom function to skip Rate Limiter validation for certain requests.
func WithNext(handler func(*fiber.Ctx) bool) Options {
	return func(c *Option) {
		c.next = handler
	}
}

// WithKeys sets a custom function to generate extra keys based on the request.
func WithKeys(handler func(*fiber.Ctx) []string) Options {
	return func(c *Option) {
		c.keys = handler
	}
}
