
# fiberopenapi

[![CI](https://github.com/oaswrap/fiberopenapi/actions/workflows/ci.yml/badge.svg)](https://github.com/oaswrap/fiberopenapi/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/oaswrap/fiberopenapi/graph/badge.svg?token=FBKZ3VZBMJ)](https://codecov.io/gh/oaswrap/fiberopenapi)
[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/fiberopenapi.svg)](https://pkg.go.dev/github.com/oaswrap/fiberopenapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaswrap/fiberopenapi)](https://goreportcard.com/report/github.com/oaswrap/fiberopenapi)
[![License](https://img.shields.io/github/license/oaswrap/fiberopenapi)](https://github.com/oaswrap/fiberopenapi/blob/main/LICENSE)

**`fiberopenapi`** is a minimal adapter for the [Fiber](https://gofiber.io) web framework that connects your routes to an OpenAPI 3.x specification using [`oaswrap/spec`](https://github.com/oaswrap/spec).

This package lets you define your Fiber routes *and* generate OpenAPI docs automatically — with simple, chainable options.

The underlying spec builder uses [`swaggest/openapi-go`](https://github.com/swaggest/openapi-go) — a robust OpenAPI 3.0/3.1 generator for Go with support for struct tags.

## ✨ Features

- ✅ Integrates Fiber routes with an OpenAPI spec.
- ✅ Uses [`oaswrap/spec`](https://github.com/oaswrap/spec) powered by [`swaggest/openapi-go`](https://github.com/swaggest/openapi-go).
- ✅ Supports struct tags for request/response models — automatically maps field properties, examples, and validation.
- ✅ Define common route metadata: summary, description, tags, security, request/response models.
- ✅ Built-in validation and schema generation helpers.

## 📦 Installation

```bash
go get github.com/oaswrap/fiberopenapi
```

## ⚡️ Quick Start

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/fiberopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	app := fiber.New()

	// Setup OpenAPI router
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("Example API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("Sample Fiber + OpenAPI"),
		option.WithDocsPath("/docs"),
	)

	// Define a simple group and route
	api := r.Group("/api")
	api.Post("/hello", helloHandler).With(
		option.Summary("Say Hello"),
		option.Request(new(HelloRequest)),
		option.Response(200, new(HelloResponse)),
	)

	// Validate and run
	if err := r.Validate(); err != nil {
		log.Fatal(err)
	}

	app.Listen(":3000")
}

// helloHandler handles the /api/hello route
func helloHandler(c *fiber.Ctx) error {
	var req HelloRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	return c.JSON(HelloResponse{
		Message: "Hello, " + req.Name,
	})
}

// Request and response types
type HelloRequest struct {
	Name string `json:"name" example:"World" validate:"required"`
}

type HelloResponse struct {
	Message string `json:"message" example:"Hello, World"`
}
```

It will generate an OpenAPI spec for the defined routes

For more examples, check out the [examples directory](https://github.com/oaswrap/fiberopenapi/tree/main/examples).

**This example demonstrates:**  
- Nested groups (`/api/v1/auth`, `/api/v1/profile`)  
- Multiple HTTP methods  
- `With` options for detailed spec generation  
- Use of struct tags to enrich models automatically.

## 📚 Documentation

- 📦 [oaswrap/spec](https://github.com/oaswrap/spec) — the core OpenAPI builder.
- 🧩 [swaggest/openapi-go](https://github.com/swaggest/openapi-go) — the underlying OpenAPI engine.
- 📖 [Fiber](https://gofiber.io) — the web framework.
- 📚 [pkg.go.dev/github.com/oaswrap/fiberopenapi](https://pkg.go.dev/github.com/oaswrap/fiberopenapi)

## 🤝 Contributing

PRs and issues welcome!  
If you find bugs or want to improve adapters or helpers, please open an issue.

## 📄 License

[MIT](./LICENSE)

Made with ❤️ by [oaswrap](https://github.com/oaswrap)