package mw

import (
	"github.com/go-nio/nio"
	"github.com/go-nio/nio/internal/random"
)

type (
	// RequestIDConfig defines the config for RequestID middleware.
	RequestIDConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper nio.Skipper

		// Generator defines a function to generate an ID.
		// Optional. Default value random.String(32).
		Generator func() string
	}
)

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = RequestIDConfig{
		Skipper:   nio.DefaultSkipper,
		Generator: generator,
	}
)

// RequestID returns a X-Request-ID middleware.
func RequestID() nio.MiddlewareFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfig(config RequestIDConfig) nio.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = generator
	}

	return func(next nio.HandlerFunc) nio.HandlerFunc {
		return func(c nio.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(nio.HeaderXRequestID)
			if rid == "" {
				rid = config.Generator()
			}
			res.Header().Set(nio.HeaderXRequestID, rid)

			return next(c)
		}
	}
}

func generator() string {
	return random.String(32)
}
