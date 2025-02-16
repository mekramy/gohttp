package uploader

import (
	"strings"

	"github.com/mekramy/gocache"
)

// Options defines a function type that modifies an Option.
type Options func(*Option)

// Option holds configuration settings for the uploader.
type Option struct {
	queue    gocache.Queue // Queue for handling files that failed to delete.
	numbered bool          // Generate numeric file name.
	prefix   string        // Add path prefix to exclude from file url.
}

// WithQueue sets the queue for handling files that failed to delete.
// Files in queue must delete manually later.
func WithQueue(queue gocache.Queue) Options {
	return func(c *Option) {
		c.queue = queue
	}
}

// WithNumbered enables generating numeric file names.
func WithNumbered() Options {
	return func(c *Option) {
		c.numbered = true
	}
}

// WithTimestamped enable generating timestamped file names.
func WithTimestamped() Options {
	return func(c *Option) {
		c.numbered = false
	}
}

// WithPrefix sets the prefix to exclude from the file URL.
func WithPrefix(prefix string) Options {
	prefix = strings.TrimSpace(prefix)
	return func(c *Option) {
		c.prefix = prefix
	}
}
