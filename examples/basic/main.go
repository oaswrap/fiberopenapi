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

	// Write schema to file (optional)
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatalf("Failed to write OpenAPI schema: %v", err)
	}
	log.Println("✅ OpenAPI schema written to: openapi.yaml")

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

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
