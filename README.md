# gohttp

`gohttp` is a Go package that provides a collection of utilities and middleware for building HTTP servers using the Fiber framework. It includes functionalities for content type validation, CSRF protection, error handling, rate limiting, session management, and file uploading.

## Features

- **Content Type Middleware**: Validate request content types such as JSON, XML, multipart form data, etc.
- **CSRF Protection**: Middleware for protecting against Cross-Site Request Forgery attacks.
- **Error Handling**: Custom error handling with logging and detailed error responses.
- **Rate Limiting**: Middleware for limiting the number of requests a client can make within a specified time period.
- **Session Management**: Middleware for managing user sessions with support for cookies and headers.
- **File Uploading**: Utilities for handling file uploads, including size and MIME type validation.

## Installation

To install the package, run:

```sh
go get github.com/mekramy/gohttp
```

## Usage

### Error Handling

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/mekramy/gohttp"
    "github.com/mekramy/gologger"
)

func main() {
    app := fiber.New(fiber.Config{
        ErrorHandler: gohttp.NewFiberErrorHandler(
            nil,
            func(ctx *fiber.Ctx, err gohttp.HttpError) error {
                return ctx.Status(500).SendString(err.Message)
            },
        ),
    })

    app.Listen(":3000")
}
```

### Session Management

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/mekramy/gocache"
    "github.com/mekramy/gohttp/session"
)

func main() {
    app := fiber.New()
    cache := gocache.NewMemoryCache()
    app.Use(session.NewMiddleware(cache))

    app.Listen(":3000")
}
```

### CSRF Protection

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/mekramy/gocache"
    "github.com/mekramy/gohttp/session"
    "github.com/mekramy/gohttp/csrf"
)

func main() {
    app := fiber.New()
    cache := gocache.NewMemoryCache()
    app.Use(session.NewMiddleware(cache))
    app.Use(csrf.NewMiddleware())

    app.Get("/csrf_token", func(c *fiber.Ctx) error{
        return c.JSON(csrf.GetToken(c))
    })

    app.Post("/refresh_csrf", func(c *fiber.Ctx) error{
        newToken, err := csrf.RefreshToken(c)
        if err != nil{
            return err
        }
        return c.JSON(newToken)
    })

    app.Listen(":3000")
}
```

### Rate Limiting

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/mekramy/gohttp/limiter"
    "github.com/mekramy/gocache"
)

func main() {
    app := fiber.New()
    cache := gocache.NewMemoryCache()
    app.Use(limiter.NewMiddleware(cache))

    app.Listen(":3000")
}
```

### Content Type Middleware

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/mekramy/gohttp/content"
)

func main() {
    app := fiber.New()
    app.Use(content.JsonOnly())

    app.Listen(":3000")
}
```

### File Uploading

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/mekramy/gohttp/uploader"
)

func main() {
    app := fiber.New()

    app.Post("/upload", func(c *fiber.Ctx) error {
        file, err := uploader.NewFiberUploader("./uploads", c, "file")
        if err != nil {
            return err
        }

        if err := file.Save(); err != nil {
            return err
        }

        return c.JSON(fiber.Map{
            "url": file.URL(),
        })
    })

    app.Listen(":3000")
}
```
