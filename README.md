
# fiberopenapi

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
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/fiberopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	app := fiber.New()

	// Initialize OpenAPI router with configuration
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("This is a sample API"),
		option.WithDocsPath("/docs"),
		option.WithServer("http://localhost:3000", option.ServerDescription("Local Server")),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer()),
		option.WithDebug(true),
	)

	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.Route("/auth", func(r fiberopenapi.Router) {
		r.Post("/login", dummyHandler).With(
			option.Summary("User Login"),
			option.Description("Endpoint for user login"),
			option.Request(new(LoginRequest)),
			option.Response(200, new(Response[Token])),
			option.Response(400, new(ErrorResponse)),
			option.Response(422, new(ValidationResponse)),
		)
		r.Get("/me", dummyHandler).With(
			option.Summary("Get User Profile"),
			option.Description("Endpoint to get the authenticated user's profile"),
			option.Security("bearerAuth"),
			option.Response(200, new(Response[User])),
			option.Response(401, new(ErrorResponse)),
		)
	}).With(
		option.RouteTags("Authentication"),
	)

	v1.Route("/profile", func(r fiberopenapi.Router) {
		r.Put("/update", dummyHandler).With(
			option.Summary("Update User Profile"),
			option.Description("Endpoint to update the user's profile"),
			option.Request(new(UpdateProfileRequest)),
			option.Response(200, new(Response[User])),
			option.Response(400, new(ErrorResponse)),
			option.Response(422, new(ValidationResponse)),
		)
		r.Put("/password", dummyHandler).With(
			option.Summary("Update Password"),
			option.Description("Endpoint to update the user's password"),
			option.Request(new(UpdatePasswordRequest)),
			option.Response(200, new(MessageResponse)),
			option.Response(400, new(ErrorResponse)),
			option.Response(422, new(ValidationResponse)),
		)
	}).With(
		option.RouteTags("Profile"),
		option.RouteSecurity("bearerAuth"),
	)

	// Validate the OpenAPI configuration
	if err := r.Validate(); err != nil {
		log.Fatalf("OpenAPI validation failed: %v", err)
	}

	// Write the OpenAPI schema to files (Optional)
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatalf("Failed to write OpenAPI schema: %v", err)
	}

	fmt.Println("Open http://localhost:3000/docs to view the OpenAPI documentation")

	app.Listen(":3000")
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