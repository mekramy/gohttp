package session

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Options is a function type that modifies an Option.
type Options func(*Option)

// Option represents configuration options for a session.
type Option struct {
	ttl       time.Duration // ttl specifies the time-to-live duration for the session.
	name      string        // name is the name of the session.
	header    bool          // header indicates whether the session should be stored in the header.
	cookie    *fiber.Cookie // cookie represents the session cookie settings.
	generator IdGenerator   // generator is the function used to generate session IDs.
}

// WithTTL returns an Options function that sets the TTL of an Option.
func WithTTL(ttl time.Duration) Options {
	return func(o *Option) {
		if ttl > 0 {
			o.ttl = ttl
		}
	}
}

// WithHeader sets the header name for the Option if the provided name is not empty.
// to indicate that a header is being used. It also clears any existing cookie settings.
func WithHeader(name string) Options {
	return func(o *Option) {
		name := strings.TrimSpace(name)
		if name != "" {
			o.name = name
			o.header = true
			o.cookie = nil
		}
	}
}

// WithCookie sets the cookie name for the Option if the provided name is not empty.
// to indicate that a cookie is being used. It also clears any existing header settings.
func WithCookie(name string, cookie fiber.Cookie) Options {
	return func(o *Option) {
		name := strings.TrimSpace(name)
		if name != "" {
			o.name = name
			o.cookie = &cookie
			o.header = false
		}
	}
}

// WithGenerator returns an Options function that sets the Generator of an Option.
func WithGenerator(generator IdGenerator) Options {
	return func(o *Option) {
		if generator != nil {
			o.generator = generator
		}
	}
}
