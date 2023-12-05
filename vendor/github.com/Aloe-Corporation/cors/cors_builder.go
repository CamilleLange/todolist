package cors

import (
	"time"

	"github.com/gin-contrib/cors"
)

// Builder is a builder pattern for creating CORS (Cross-Origin Resource
// Sharing) configurations. It provides methods to set various CORS options and
// ultimately build a cors.Config.
type Builder struct {
	*cors.Config
}

// New creates a new Builder instance with default CORS configuration values.
func (builder *Builder) New() *Builder {
	builder.Config = &cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return builder
}

// NewFromConfig creates a new Builder instance with a custom configuration
// provided by the Conf parameter.
func (builder *Builder) NewFromConfig(c *Conf) *Builder {
	builder.Config = &cors.Config{
		AllowOrigins:     c.AllowOrigins,
		AllowMethods:     c.AllowMethods,
		AllowHeaders:     c.AllowHeaders,
		ExposeHeaders:    c.ExposeHeaders,
		AllowCredentials: c.AllowCredentials,
		MaxAge:           c.MaxAge,
	}
	return builder
}

// WithOrigins sets the allowed origins for CORS.
func (builder *Builder) WithOrigins(origins ...string) *Builder {
	builder.Config.AllowOrigins = origins
	return builder
}

// WithMethods sets the allowed HTTP methods for CORS.
func (builder *Builder) WithMethods(methods ...string) *Builder {
	builder.Config.AllowMethods = methods
	return builder
}

// WithHeaders sets the allowed headers for CORS.
func (builder *Builder) WithHeaders(headers ...string) *Builder {
	builder.Config.AllowHeaders = headers
	return builder
}

// WithCredentials sets whether credentials (including cookies) can be sent
// with the CORS request.
func (builder *Builder) WithCredentials(allowCredentials bool) *Builder {
	builder.Config.AllowCredentials = allowCredentials
	return builder
}

// Build returns the finalized cors.Config after applying the configured options.
func (builder *Builder) Build() *cors.Config {
	return builder.Config
}
