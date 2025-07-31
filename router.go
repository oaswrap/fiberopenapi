package fiberopenapi

import (
	stdpath "path"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/fiberopenapi/internal/constant"
	"github.com/oaswrap/fiberopenapi/internal/handler"
	"github.com/oaswrap/fiberopenapi/internal/util"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

// Router defines the interface for an OpenAPI router.
type Router interface {
	// Use applies middleware to the router.
	Use(args ...any) Router

	// Get registers a GET route.
	Get(path string, handler ...fiber.Handler) Route
	// Head registers a HEAD route.
	Head(path string, handler ...fiber.Handler) Route
	// Post registers a POST route.
	Post(path string, handler ...fiber.Handler) Route
	// Put registers a PUT route.
	Put(path string, handler ...fiber.Handler) Route
	// Patch registers a PATCH route.
	Patch(path string, handler ...fiber.Handler) Route
	// Delete registers a DELETE route.
	Delete(path string, handler ...fiber.Handler) Route
	// Connect registers a CONNECT route.
	Connect(path string, handler ...fiber.Handler) Route
	// Options registers an OPTIONS route.
	Options(path string, handler ...fiber.Handler) Route
	// Trace registers a TRACE route.
	Trace(path string, handler ...fiber.Handler) Route

	// Add registers a route with the specified method and path.
	Add(method, path string, handler ...fiber.Handler) Route
	// Static serves static files from the specified root directory.
	Static(prefix, root string, config ...fiber.Static) Router

	// Group creates a new sub-router with the specified prefix and handlers.
	// The prefix is prepended to all routes in the sub-router.
	Group(prefix string, handlers ...fiber.Handler) Router

	// Route creates a new sub-router with the specified prefix and applies options.
	Route(prefix string, fn func(router Router)) Router

	// With applies options to the router.
	// This allows you to configure tags, security, and visibility for the routes.
	With(opts ...option.GroupOption) Router

	// Validate checks for errors at OpenAPI router initialization.
	//
	// It returns an error if there are issues with the OpenAPI configuration.
	Validate() error

	// GenerateOpenAPISchema generates the OpenAPI schema in the specified format.
	// Supported formats are "json" and "yaml".
	// If no format is specified, "yaml" is used by default.
	GenerateOpenAPISchema(format ...string) ([]byte, error)

	WriteSchemaTo(filePath string) error
}

type router struct {
	fiberRouter fiber.Router
	specRouter  spec.Router
	generator   *spec.Generator
}

func NewRouter(r fiber.Router, opts ...option.OpenAPIOption) Router {
	defaultOpts := []option.OpenAPIOption{
		option.WithTitle(constant.DefaultTitle),
		option.WithDescription(constant.DefaultDescription),
		option.WithVersion(constant.DefaultVersion),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	}
	opts = append(defaultOpts, opts...)
	generator := spec.NewGenerator(opts...)
	cfg := generator.Config()

	rr := &router{
		fiberRouter: r,
		specRouter:  generator,
		generator:   generator,
	}

	// If OpenAPI is disabled, return the router without any OpenAPI functionality.
	// This allows the application to run without OpenAPI if desired.
	if cfg.DisableOpenAPI {
		return rr
	}

	handler := handler.NewOpenAPIHandler(cfg, generator)
	openapiPath := stdpath.Join(cfg.DocsPath, constant.OpenAPIFileName)

	r.Get(cfg.DocsPath, handler.Docs)
	r.Get(openapiPath, handler.OpenAPIYaml)

	return rr
}

func (r *router) Use(args ...any) Router {
	r.fiberRouter.Use(args...)
	return r
}

func (r *router) Get(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodGet, path, handler...)
}

func (r *router) Head(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodHead, path, handler...)
}

func (r *router) Post(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodPost, path, handler...)
}

func (r *router) Put(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodPut, path, handler...)
}

func (r *router) Patch(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodPatch, path, handler...)
}

func (r *router) Delete(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodDelete, path, handler...)
}

func (r *router) Connect(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodConnect, path, handler...)
}

func (r *router) Options(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodOptions, path, handler...)
}

func (r *router) Trace(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodTrace, path, handler...)
}

func (r *router) Add(method, path string, handler ...fiber.Handler) Route {
	fr := r.fiberRouter.Add(method, path, handler...)
	sr := r.specRouter.Add(method, util.ConvertPath(path))

	route := &route{
		fr: fr,
		sr: sr,
	}

	return route
}

func (r *router) Static(prefix, root string, config ...fiber.Static) Router {
	r.fiberRouter.Static(prefix, root, config...)
	return r
}

func (r *router) Group(prefix string, handlers ...fiber.Handler) Router {
	rr := r.fiberRouter.Group(prefix, handlers...)
	sr := r.specRouter.Group(prefix)

	return &router{
		fiberRouter: rr,
		specRouter:  sr,
	}
}

func (r *router) Route(prefix string, fn func(router Router)) Router {
	fr := r.fiberRouter.Group(prefix)
	sr := r.specRouter.Group(prefix)

	subRouter := &router{
		fiberRouter: fr,
		specRouter:  sr,
	}

	fn(subRouter)

	return subRouter
}

func (r *router) With(opts ...option.GroupOption) Router {
	r.specRouter.Use(opts...)
	return r
}

func (r *router) Validate() error {
	if err := r.generator.Validate(); err != nil {
		return err
	}
	return nil
}

func (r *router) GenerateOpenAPISchema(formats ...string) ([]byte, error) {
	return r.generator.GenerateSchema(formats...)
}

func (r *router) WriteSchemaTo(path string) error {
	return r.generator.WriteSchemaTo(path)
}
