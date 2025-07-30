package fiberopenapi

import (
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

func newConfig(opts ...option.OpenAPIOption) *spec.Config {
	cfg := &spec.Config{
		OpenAPIVersion:  "3.1.0",
		Title:           "Fiber OpenAPI",
		Description:     nil,
		DisableOpenAPI:  false,
		DocsPath:        "/docs",
		SwaggerConfig:   &spec.SwaggerConfig{},
		SecuritySchemes: make(map[string]*spec.SecurityScheme),
		Logger:          &noopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

type noopLogger struct{}

func (noopLogger) Printf(format string, v ...any) {}
