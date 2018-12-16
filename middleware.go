package nio

import "net/http"

type (
	// MiddlewareFunc defines a function to process middleware.
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(Context) bool
)

// WrapMiddleware wraps `func(http.Handler) http.Handler` into `nio.MiddlewareFunc`
func WrapMiddleware(m func(http.Handler) http.Handler) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) (err error) {
			m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				err = next(c)
			})).ServeHTTP(c.Response(), c.Request())
			return
		}
	}
}

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(Context) bool {
	return false
}
