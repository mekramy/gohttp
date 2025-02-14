package gohttp

import (
	"fmt"
	"mime/multipart"
	"runtime"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/inhies/go-bytesize"
)

// HttpError represents an HTTP error with additional context information.
type HttpError struct {
	Line    int            // Line number where the error occurred.
	File    string         // File name where the error occurred.
	Body    map[string]any // Request body data.
	Status  int            // HTTP status code.
	Message string         // Error message.
}

// Error returns the error message.
func (he HttpError) Error() string {
	return he.Message
}

// NewError creates a new HttpError with the provided error message and optional status code.
// If no status code is provided, it defaults to 500.
// It also captures the file and line number where the error occurred.
func NewError(e string, status ...int) error {
	code := 500
	if len(status) > 0 {
		code = status[0]
	}

	file, line, _ := realCaller()
	return HttpError{
		Line:    line,
		File:    file,
		Body:    nil,
		Status:  code,
		Message: e,
	}
}

// NewFormError creates a new HttpError with the provided error message, request context, and optional status code.
// It captures the file and line number where the error occurred and includes request body data if available.
func NewFormError(e string, ctx *fiber.Ctx, status ...int) error {
	code := 500
	if len(status) > 0 {
		code = status[0]
	}

	var body map[string]any
	if ctx != nil {
		body = make(map[string]any)
		if form, err := ctx.MultipartForm(); err == nil && form != nil {
			for k, v := range form.Value {
				if len(v) == 1 {
					body["form."+k] = v[0]
				} else if len(v) > 1 {
					body["form."+k] = v
				} else {
					body["form."+k] = nil
				}
			}

			for k, files := range form.File {
				values := make([]string, 0)
				for _, file := range files {
					size := bytesize.New(float64(file.Size))
					mime := detectMime(file)
					values = append(
						values,
						fmt.Sprintf("%s [%s] (%s)", file.Filename, size, mime),
					)
				}

				if len(values) == 0 {
					body["file."+k] = nil
				} else {
					body["file."+k] = values
				}
			}
		} else {
			var form map[string]any
			if err := ctx.BodyParser(&form); err != nil {
				body["form"] = err.Error()
			} else if len(form) == 0 {
				body["form"] = nil
			} else {
				for k, v := range form {
					body["form."+k] = v
				}
			}
		}
	}

	file, line, _ := realCaller()
	return HttpError{
		Line:    line,
		File:    file,
		Body:    body,
		Status:  code,
		Message: e,
	}
}

// detectMime detects the MIME type of the provided file header.
// It opens the file, reads its content, and returns the MIME type as a string.
// If the MIME type cannot be determined, it returns "?".
func detectMime(file *multipart.FileHeader) string {
	f, err := file.Open()
	if err == nil {
		defer f.Close()
		if mime, err := mimetype.DetectReader(f); err == nil && mime != nil {
			return mime.String()
		}
	}

	return "?"
}

// realCaller returns the file name and line number of error caller func.
func realCaller() (string, int, bool) {
	if _, f, l, ok := runtime.Caller(2); ok {
		return f, l, true
	}
	return "", 0, false
}
