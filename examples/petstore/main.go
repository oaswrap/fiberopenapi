package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/fiberopenapi"
	"github.com/oaswrap/fiberopenapi/examples/petstore/handler"
	"github.com/oaswrap/fiberopenapi/examples/petstore/model"
	"github.com/oaswrap/fiberopenapi/examples/petstore/repository"
	"github.com/oaswrap/spec/option"
)

func main() {
	repo := repository.NewDummyPetRepository()
	handler := handler.NewPetHandler(repo)

	app := fiber.New()
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("Petstore API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("Sample Petstore API using Fiber and OpenAPI"),
	)

	r.Route("/pets", func(r fiberopenapi.Router) {
		r.Get("", handler.GetAllPets).With(
			option.Summary("Get all pets"),
			option.Description("Returns a list of all pets in the store"),
			option.Response(200, new([]model.Pet)),
		)
		r.Post("", handler.CreatePet).With(
			option.Summary("Create a new pet"),
			option.Description("Creates a new pet in the store"),
			option.Request(new(model.CreatePetRequest)),
			option.Response(201, new(model.Pet)),
		)
		r.Put("", handler.UpdatePet).With(
			option.Summary("Update an existing pet"),
			option.Description("Updates an existing pet in the store"),
			option.Request(new(model.UpdatePetRequest)),
			option.Response(200, new(model.Pet)),
		)
		r.Get("/findByStatus", handler.FindPetsByStatus).With(
			option.Summary("Find pets by status"),
			option.Description("Returns a list of pets based on their status"),
			option.Request(new(model.FindPetsByStatusRequest)),
			option.Response(200, new([]model.Pet)),
		)
		r.Get("/findByTags", handler.FindPetsByTags).With(
			option.Summary("Find pets by tags"),
			option.Description("Returns a list of pets based on their tags"),
			option.Request(new(model.FindPetsByTagsRequest)),
			option.Response(200, new([]model.Pet)),
		)
		r.Get("/:petId", handler.GetPetByID).With(
			option.Summary("Get pet by ID"),
			option.Description("Returns a single pet by its ID"),
			option.Request(new(struct {
				ID int64 `path:"petId"`
			})),
			option.Response(200, new(model.Pet)),
		)
		r.Post("/:petId", handler.UpdatePetFormData).With(
			option.Summary("Update pet by form data"),
			option.Description("Updates a pet using form data"),
			option.Request(new(model.UpdatePetFormData)),
			option.Response(200, new(model.Pet)),
		)
		r.Delete("/:petId", handler.DeletePet).With(
			option.Summary("Delete a pet"),
			option.Description("Deletes a pet from the store"),
			option.Request(new(struct {
				ID int64 `path:"petId"`
			})),
			option.Response(204, nil),
		)
	}).With(option.GroupTags("Pets"))

	// Validate the OpenAPI configuration
	if err := r.Validate(); err != nil {
		log.Fatalf("OpenAPI validation failed: %v", err)
	}

	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatalf("Failed to write OpenAPI schema: %v", err)
	}

	log.Println("✅ OpenAPI schema written to: openapi.yaml")

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	app.Listen(":3000")
}
