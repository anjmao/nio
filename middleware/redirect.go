package middleware

import (
	"net/http"

	"github.com/dostack/dapi"
)

// RedirectConfig defines the config for Redirect middleware.
type RedirectConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper

	// Status code to be used when redirecting the request.
	// Optional. Default value http.StatusMovedPermanently.
	Code int `yaml:"code"`
}

// redirectLogic represents a function that given a scheme, host and uri
// can both: 1) determine if redirect is needed (will set ok accordingly) and
// 2) return the appropriate redirect url.
type redirectLogic func(scheme, host, uri string) (ok bool, url string)

const www = "www"

// DefaultRedirectConfig is the default Redirect middleware config.
var DefaultRedirectConfig = RedirectConfig{
	Skipper: DefaultSkipper,
	Code:    http.StatusMovedPermanently,
}

// HTTPSRedirect redirects http requests to https.
// For example, http://dostack.com will be redirect to https://dostack.com.
//
// Usage `Dapi#Pre(HTTPSRedirect())`
func HTTPSRedirect() dapi.MiddlewareFunc {
	return HTTPSRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSRedirect()`.
func HTTPSRedirectWithConfig(config RedirectConfig) dapi.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = scheme != "https"; ok {
			url = "https://" + host + uri
		}
		return
	})
}

// HTTPSWWWRedirect redirects http requests to https www.
// For example, http://dostack.com will be redirect to https://www.dostack.com.
//
// Usage `Dapi#Pre(HTTPSWWWRedirect())`
func HTTPSWWWRedirect() dapi.MiddlewareFunc {
	return HTTPSWWWRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSWWWRedirect()`.
func HTTPSWWWRedirectWithConfig(config RedirectConfig) dapi.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = scheme != "https" && host[:3] != www; ok {
			url = "https://www." + host + uri
		}
		return
	})
}

// HTTPSNonWWWRedirect redirects http requests to https non www.
// For example, http://www.dostack.com will be redirect to https://dostack.com.
//
// Usage `Dapi#Pre(HTTPSNonWWWRedirect())`
func HTTPSNonWWWRedirect() dapi.MiddlewareFunc {
	return HTTPSNonWWWRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSNonWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSNonWWWRedirect()`.
func HTTPSNonWWWRedirectWithConfig(config RedirectConfig) dapi.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = scheme != "https"; ok {
			if host[:3] == www {
				host = host[4:]
			}
			url = "https://" + host + uri
		}
		return
	})
}

// WWWRedirect redirects non www requests to www.
// For example, http://dostack.com will be redirect to http://www.dostack.com.
//
// Usage `Dapi#Pre(WWWRedirect())`
func WWWRedirect() dapi.MiddlewareFunc {
	return WWWRedirectWithConfig(DefaultRedirectConfig)
}

// WWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `WWWRedirect()`.
func WWWRedirectWithConfig(config RedirectConfig) dapi.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = host[:3] != www; ok {
			url = scheme + "://www." + host + uri
		}
		return
	})
}

// NonWWWRedirect redirects www requests to non www.
// For example, http://www.dostack.com will be redirect to http://dostack.com.
//
// Usage `Dapi#Pre(NonWWWRedirect())`
func NonWWWRedirect() dapi.MiddlewareFunc {
	return NonWWWRedirectWithConfig(DefaultRedirectConfig)
}

// NonWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `NonWWWRedirect()`.
func NonWWWRedirectWithConfig(config RedirectConfig) dapi.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = host[:3] == www; ok {
			url = scheme + "://" + host[4:] + uri
		}
		return
	})
}

func redirect(config RedirectConfig, cb redirectLogic) dapi.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next dapi.HandlerFunc) dapi.HandlerFunc {
		return func(c dapi.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req, scheme := c.Request(), c.Scheme()
			host := req.Host
			if ok, url := cb(scheme, host, req.RequestURI); ok {
				return c.Redirect(config.Code, url)
			}

			return next(c)
		}
	}
}
