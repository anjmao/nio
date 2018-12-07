package middleware

import (
	"net/http"

	"github.com/dostack/dapi"
)

type (
	// MethodOverrideConfig defines the config for MethodOverride middleware.
	MethodOverrideConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Getter is a function that gets overridden method from the request.
		// Optional. Default values MethodFromHeader(dapi.HeaderXHTTPMethodOverride).
		Getter MethodOverrideGetter
	}

	// MethodOverrideGetter is a function that gets overridden method from the request
	MethodOverrideGetter func(dapi.Context) string
)

var (
	// DefaultMethodOverrideConfig is the default MethodOverride middleware config.
	DefaultMethodOverrideConfig = MethodOverrideConfig{
		Skipper: DefaultSkipper,
		Getter:  MethodFromHeader(dapi.HeaderXHTTPMethodOverride),
	}
)

// MethodOverride returns a MethodOverride middleware.
// MethodOverride  middleware checks for the overridden method from the request and
// uses it instead of the original method.
//
// For security reasons, only `POST` method can be overridden.
func MethodOverride() dapi.MiddlewareFunc {
	return MethodOverrideWithConfig(DefaultMethodOverrideConfig)
}

// MethodOverrideWithConfig returns a MethodOverride middleware with config.
// See: `MethodOverride()`.
func MethodOverrideWithConfig(config MethodOverrideConfig) dapi.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultMethodOverrideConfig.Skipper
	}
	if config.Getter == nil {
		config.Getter = DefaultMethodOverrideConfig.Getter
	}

	return func(next dapi.HandlerFunc) dapi.HandlerFunc {
		return func(c dapi.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			if req.Method == http.MethodPost {
				m := config.Getter(c)
				if m != "" {
					req.Method = m
				}
			}
			return next(c)
		}
	}
}

// MethodFromHeader is a `MethodOverrideGetter` that gets overridden method from
// the request header.
func MethodFromHeader(header string) MethodOverrideGetter {
	return func(c dapi.Context) string {
		return c.Request().Header.Get(header)
	}
}

// MethodFromForm is a `MethodOverrideGetter` that gets overridden method from the
// form parameter.
func MethodFromForm(param string) MethodOverrideGetter {
	return func(c dapi.Context) string {
		return c.FormValue(param)
	}
}

// MethodFromQuery is a `MethodOverrideGetter` that gets overridden method from
// the query parameter.
func MethodFromQuery(param string) MethodOverrideGetter {
	return func(c dapi.Context) string {
		return c.QueryParam(param)
	}
}
