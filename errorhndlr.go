package gohttp

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/mekramy/gologger"
)

// ErrorCallback is a function type that handles custom error responses.
type ErrorCallback func(ctx *fiber.Ctx, err HttpError) error

// NewFiberErrorHandler creates a new Fiber error handler with logging and custom error response capabilities.
// It takes a logger, an optional error callback, and a list of status codes to log.
// If the error matches one of the provided status codes, it will be logged using the provided logger.
// If an error callback is provided, it will be used to handle the error response; otherwise, a default plain text response will be sent.
// For relative file name in log use os.Setenv("APP_ROOT", "your/project/root") to define your project root.
func NewFiberErrorHandler(l gologger.Logger, cb ErrorCallback, codes ...int) fiber.ErrorHandler {
	relative := func(path string) string {
		root := filepath.ToSlash(os.Getenv("APP_ROOT"))
		path = filepath.ToSlash(path)
		if root != "" {
			if p, err := filepath.Rel(root, path); err == nil {
				return p
			}
		}

		return filepath.ToSlash(path)
	}

	return func(ctx *fiber.Ctx, err error) error {
		// Parse error
		file := ""
		line := 0
		var body map[string]any
		status := fiber.StatusInternalServerError
		message := "Internal Server Error"

		if fe, ok := err.(*fiber.Error); ok {
			status = fe.Code
			message = fe.Error()
		}

		if he, ok := err.(HttpError); ok {
			line = he.Line
			file = he.File
			message = he.Error()
			status = he.Status
			body = he.Body
		}

		// Log
		if l != nil && (len(codes) == 0 || slices.Contains(codes, status)) {
			params := make([]gologger.LogOptions, 0)
			params = append(params, gologger.With("file", relative(file)))
			params = append(params, gologger.With("line", line))
			params = append(params, gologger.With("status", status))
			params = append(params, gologger.With("ip", ctx.IP()))
			params = append(params, gologger.With("path", ctx.Path()))
			params = append(params, gologger.With("method", ctx.Method()))
			params = append(params, gologger.WithMessage(message))
			for k, v := range body {
				params = append(params, gologger.With(k, v))

			}
			l.Error(params...)
		}

		// Return error
		if cb != nil {
			return cb(ctx, HttpError{
				Line:    line,
				File:    file,
				Body:    body,
				Status:  status,
				Message: message,
			})
		} else {
			ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return ctx.Status(status).SendString(message)
		}
	}
}
