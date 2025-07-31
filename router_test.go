package fiberopenapi_test

import (
	"flag"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/fiberopenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func PingHandler(c *fiber.Ctx) error {
	return c.SendString("pong")
}

type AllBasicDataTypes struct {
	Int     int     `json:"int"`
	Int8    int8    `json:"int8"`
	Int16   int16   `json:"int16"`
	Int32   int32   `json:"int32"`
	Int64   int64   `json:"int64"`
	Uint    uint    `json:"uint"`
	Uint8   uint8   `json:"uint8"`
	Uint16  uint16  `json:"uint16"`
	Uint32  uint32  `json:"uint32"`
	Uint64  uint64  `json:"uint64"`
	Float32 float32 `json:"float32"`
	Float64 float64 `json:"float64"`
	Byte    byte    `json:"byte"`
	Rune    rune    `json:"rune"`
	String  string  `json:"string"`
	Bool    bool    `json:"bool"`
}

type AllBasicDataTypesPointers struct {
	Int     *int     `json:"int"`
	Int8    *int8    `json:"int8"`
	Int16   *int16   `json:"int16"`
	Int32   *int32   `json:"int32"`
	Int64   *int64   `json:"int64"`
	Uint    *uint    `json:"uint"`
	Uint8   *uint8   `json:"uint8"`
	Uint16  *uint16  `json:"uint16"`
	Uint32  *uint32  `json:"uint32"`
	Uint64  *uint64  `json:"uint64"`
	Float32 *float32 `json:"float32"`
	Float64 *float64 `json:"float64"`
	Byte    *byte    `json:"byte"`
	Rune    *rune    `json:"rune"`
	String  *string  `json:"string"`
	Bool    *bool    `json:"bool"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ErrorResponse struct {
	Status int    `json:"status" example:"400"`
	Title  string `json:"title" example:"Bad Request"`
	Detail string `json:"detail,omitempty" example:"Invalid input provided"`
}

type ValidationResponse struct {
	Status int          `json:"status" example:"422"`
	Title  string       `json:"title" example:"Validation Error"`
	Detail string       `json:"detail,omitempty" example:"Input validation failed"`
	Errors []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
	Field   string `json:"field" example:"username"`
	Message string `json:"message" example:"Username is required"`
}

type Response[T any] struct {
	Status int `json:"status" example:"200"`
	Data   T   `json:"data"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfile struct {
	ID              string    `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	EmailVerifiedAt NullTime  `json:"email_verified_at"`
	FullName        string    `json:"full_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type NullTime struct {
	Time  time.Time `json:"time"`
	Valid bool      `json:"valid"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Pet struct {
	ID        int64     `json:"id"`
	Category  *Category `json:"category,omitempty"`
	Name      string    `json:"name" validate:"required"`
	PhotoURLs []string  `json:"photoUrls" validate:"required"`
	Tags      []Tag     `json:"tags,omitempty"`
	Status    string    `json:"status,omitempty" enum:"available,pending,sold"`
}

type ApiResponse struct {
	Code    int32  `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type FindPetByIdRequest struct {
	ID int64 `params:"petId" path:"petId"`
}

type FindPetsByStatusRequest struct {
	Status string `query:"status" validate:"required" enum:"available,pending,sold"`
}

type FindPetsByTagsRequest struct {
	Tags []string `query:"tags" required:"false"`
}

type DeletePetRequest struct {
	ApiKey string `header:"api_key"`
	ID     int64  `params:"petId" path:"petId"`
}

type UpdatePetFormDataRequest struct {
	ID     int64  `params:"petId" path:"petId"`
	Name   string `formData:"name" validate:"required"`
	Status string `formData:"status" enum:"available,pending,sold"`
}

type UploadImageRequest struct {
	ID                 int64           `params:"petId" path:"petId"`
	AdditionalMetaData string          `query:"additionalMetadata"`
	_                  *multipart.File `contentType:"application/octet-stream"`
}

func TestRouter(t *testing.T) {
	tests := []struct {
		name    string
		golden  string
		options []option.OpenAPIOption
		setup   func(r fiberopenapi.Router)
	}{
		{
			name:   "Basic Data Types",
			golden: "basic_data_types.yaml",
			setup: func(r fiberopenapi.Router) {
				r.Post("/data-types", PingHandler).With(
					option.Summary("All Basic Data Types"),
					option.Description("Endpoint to test all basic data types"),
					option.Request(new(AllBasicDataTypes)),
					option.Response(200, new(AllBasicDataTypes)),
				)
			},
		},
		{
			name:   "Basic Data Types Pointers",
			golden: "basic_data_types_pointers.yaml",
			setup: func(r fiberopenapi.Router) {
				r.Put("/data-types-pointers", PingHandler).With(
					option.Summary("All Basic Data Types Pointers"),
					option.Description("Endpoint to test all basic data types with pointers"),
					option.Request(new(AllBasicDataTypesPointers)),
					option.Response(200, new(AllBasicDataTypesPointers)),
				)
			},
		},
		{
			name:   "Generic Response",
			golden: "generic_response.yaml",
			setup: func(r fiberopenapi.Router) {
				r.Post("/auth/login", PingHandler).With(
					option.Summary("User Login"),
					option.Description("Endpoint for user login"),
					option.Request(new(LoginRequest)),
					option.Response(200, new(Response[Token])),
					option.Response(400, new(ErrorResponse)),
					option.Response(422, new(ValidationResponse)),
				)
			},
		},
		{
			name:   "Custom Type Mapping",
			golden: "type_mapping.yaml",
			options: []option.OpenAPIOption{
				option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
				option.WithReflectorConfig(
					option.TypeMapping(NullTime{}, new(time.Time)),
				),
			},
			setup: func(r fiberopenapi.Router) {
				r.Get("/auth/me", PingHandler).With(
					option.Summary("Get User Profile"),
					option.Description("Endpoint to get the authenticated user's profile"),
					option.Security("bearerAuth"),
					option.Response(200, new(Response[UserProfile])),
					option.Response(401, new(ErrorResponse)),
				)
			},
		},
		{
			name:   "Pet Store API",
			golden: "petstore.yaml",
			options: []option.OpenAPIOption{
				option.WithTitle("Pet Store API - OpenAPI 3.1"),
				option.WithVersion("1.0.0"),
				option.WithDescription("This is a sample Pet Store API using OpenAPI 3.1"),
				option.WithDocsPath("/docs"),
				option.WithServer("https://petstore3.swagger.io", option.ServerDescription("Pet Store Server")),
				option.WithSecurity("petstore_auth", option.SecurityOAuth2(
					openapi.OAuthFlows{
						Implicit: &openapi.OAuthFlowsDefsImplicit{
							AuthorizationURL: "https://petstore3.swagger.io/oauth/authorize",
							Scopes: map[string]string{
								"write:pets": "modify pets in your account",
								"read:pets":  "read your pets",
							},
						},
					},
				)),
			},
			setup: func(r fiberopenapi.Router) {
				api := r.Group("/api")
				v3 := api.Group("/v3")
				v3.Route("/pet", func(r fiberopenapi.Router) {
					r.Get("/findByStatus", nil).With(
						option.Summary("Finds Pets by status."),
						option.Description("Multiple status values can be provided with comma separated strings"),
						option.Request(new(FindPetsByStatusRequest)),
						option.Response(200, new([]Pet)),
					)
					r.Get("/findByTags", nil).With(
						option.Summary("Finds Pets by tags."),
						option.Description("Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing."),
						option.Request(new(FindPetsByTagsRequest)),
						option.Response(200, new([]Pet)),
					)
					r.Get("/:petId", nil).With(
						option.Summary("Find a pet by ID."),
						option.Description("Returns a single pet."),
						option.Request(new(FindPetByIdRequest)),
						option.Response(200, new(Pet)),
					)
					r.Post("/:petId", nil).With(
						option.Summary("Updates a pet in the store with form data."),
						option.Description("Update a pet resource based on form data."),
						option.Request(new(UpdatePetFormDataRequest)),
						option.Response(200, new(Pet)),
					)
					r.Delete("/:petId", nil).With(
						option.Summary("Deletes a pet."),
						option.Request(new(DeletePetRequest)),
					)
					r.Post("/:petId/uploadImage", nil).With(
						option.Summary("Uploads an image."),
						option.Description("Uploads image of the pet."),
						option.Request(new(UploadImageRequest)),
						option.Response(200, new(ApiResponse)),
					)
					r.Post("/", nil).With(
						option.Summary("Add a new pet to the store."),
						option.Request(new(Pet)),
						option.Response(200, new(Pet)),
					)
					r.Put("/", nil).With(
						option.Summary("Update an existing pet."),
						option.Description("Update an existing pet by Id."),
						option.Request(new(Pet)),
						option.Response(200, new(Pet)),
					)
				}).With(option.GroupTags("pet"), option.GroupSecurity("petstore_auth", "write:pets", "read:pets"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			opts := []option.OpenAPIOption{
				option.WithTitle("Test API " + tt.name),
				option.WithVersion("1.0.0"),
				option.WithDescription("This is a test API for " + tt.name),
				option.WithReflectorConfig(
					option.RequiredPropByValidateTag(),
				),
			}
			if len(tt.options) > 0 {
				opts = append(opts, tt.options...)
			}
			r := fiberopenapi.NewRouter(app, opts...)

			tt.setup(r)

			err := r.Validate()
			assert.NoError(t, err, "failed to validate OpenAPI configuration")

			// Test the OpenAPI schema generation
			schema, err := r.GenerateOpenAPISchema()

			require.NoError(t, err, "failed to generate OpenAPI schema")
			goldenFile := filepath.Join("testdata", tt.golden)

			if *update {
				err = r.WriteSchemaTo(goldenFile)
				require.NoError(t, err, "failed to write golden file")
				t.Logf("Updated golden file: %s", goldenFile)
			}

			want, err := os.ReadFile(goldenFile)
			require.NoError(t, err, "failed to read golden file %s", goldenFile)

			testutil.EqualYAML(t, want, schema)
		})
	}
}

func TestRouterWriteSchemaTo(t *testing.T) {
	app := fiber.New()
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("Test API Write Schema"),
		option.WithVersion("1.0.0"),
		option.WithDescription("This is a test API for writing OpenAPI schema to file"),
	)

	r.Get("/ping", PingHandler).With(
		option.Summary("Ping Endpoint"),
		option.Description("Endpoint to test ping functionality"),
	)

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	tempFile, err := os.CreateTemp("", "openapi-schema-*.yaml")
	require.NoError(t, err, "failed to create temporary file for OpenAPI schema")
	defer func() {
		err := os.Remove(tempFile.Name())
		require.NoError(t, err, "failed to remove temporary file")
	}()

	err = r.WriteSchemaTo(tempFile.Name())
	require.NoError(t, err, "failed to write OpenAPI schema to file")

	schema, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err, "failed to read OpenAPI schema from file")
	assert.NotEmpty(t, schema, "expected non-empty OpenAPI schema")
}
