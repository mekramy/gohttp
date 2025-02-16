// Package uploader provides functionality for handling file uploads.
package uploader

import (
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mekramy/goutils"
	"github.com/valyala/fasthttp"
)

// Uploader is an interface that defines methods for handling file uploads.
type Uploader interface {
	// IsNil checks if the uploader is nil.
	IsNil() bool

	// ValidateSize checks if the file size is within the specified limit.
	// Use B, KB, MB, GB for size string
	ValidateSize(min, max string) (bool, error)

	// ValidateMime checks if the file MIME type is among the allowed types.
	ValidateMime(mimes ...string) (bool, error)

	// Path returns the file path where the uploaded file is stored.
	Path() string

	// URL returns the URL where the uploaded file can be accessed.
	URL() string

	// Save stores the uploaded file.
	Save() error

	// Delete removes the uploaded file.
	Delete() error

	// SafeDelete removes the uploaded file safely, queueing the file name on failure.
	SafeDelete()
}

// NewUploader creates a new Uploader instance with the given root directory and file header.
func NewUploader(root string, file *multipart.FileHeader, options ...Options) (Uploader, error) {
	// Initialize and normalize
	var name string
	root = strings.TrimSpace(root)

	// Create option with default values.
	option := &Option{
		queue:    nil,
		numbered: false,
		prefix:   "",
	}
	for _, opt := range options {
		opt(option)
	}

	// Generate file name
	if file != nil {
		if option.numbered {
			n, err := goutils.NumberedFile(root, file.Filename)
			if err != nil {
				return nil, err
			}
			name = n
		} else {
			name = goutils.TimestampedFile(file.Filename)
		}
	}

	// Create and return the uploader instance.
	u := &uploader{
		opt:  *option,
		file: file,
		name: name,
		root: root,
	}
	return u, nil
}

// NewFiberUploader creates a new Uploader instance for a Fiber context.
func NewFiberUploader(root string, c *fiber.Ctx, name string, options ...Options) (Uploader, error) {
	file, err := c.FormFile(name)
	if err == fasthttp.ErrMissingFile {
		return NewUploader(root, nil, options...)
	}

	if err != nil {
		return nil, err
	}

	return NewUploader(root, file, options...)
}

// FiberFile retrieves a file from a Fiber context by its form field name.
// If the file is not found, it returns nil without an error.
// If another error occurs, it returns the error.
func FiberFile(c *fiber.Ctx, name string) (*multipart.FileHeader, error) {
	f, err := c.FormFile(name)
	if err == fasthttp.ErrMissingFile {
		return nil, nil
	}
	return f, err
}
