package nio

import (
	"net/http"
	"path"
)

type (
	// Group is a set of sub-routes for a specified route. It can be used for inner
	// routes that share a common middleware or functionality that should be separate
	// from the parent nio instance while still inheriting from it.
	Group struct {
		prefix     string
		middleware []MiddlewareFunc
		nio       *Nio
	}
)

// Use implements `Nio#Use()` for sub-routes within the Group.
func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
	// Allow all requests to reach the group as they might get dropped if router
	// doesn't find a match, making none of the group middleware process.
	for _, p := range []string{"", "/*"} {
		g.nio.Any(path.Clean(g.prefix+p), func(c Context) error {
			return NotFoundHandler(c)
		}, g.middleware...)
	}
}

// CONNECT implements `Nio#CONNECT()` for sub-routes within the Group.
func (g *Group) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodConnect, path, h, m...)
}

// DELETE implements `Nio#DELETE()` for sub-routes within the Group.
func (g *Group) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodDelete, path, h, m...)
}

// GET implements `Nio#GET()` for sub-routes within the Group.
func (g *Group) GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodGet, path, h, m...)
}

// HEAD implements `Nio#HEAD()` for sub-routes within the Group.
func (g *Group) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodHead, path, h, m...)
}

// OPTIONS implements `Nio#OPTIONS()` for sub-routes within the Group.
func (g *Group) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodOptions, path, h, m...)
}

// PATCH implements `Nio#PATCH()` for sub-routes within the Group.
func (g *Group) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPatch, path, h, m...)
}

// POST implements `Nio#POST()` for sub-routes within the Group.
func (g *Group) POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPost, path, h, m...)
}

// PUT implements `Nio#PUT()` for sub-routes within the Group.
func (g *Group) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPut, path, h, m...)
}

// TRACE implements `Nio#TRACE()` for sub-routes within the Group.
func (g *Group) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodTrace, path, h, m...)
}

// Any implements `Nio#Any()` for sub-routes within the Group.
func (g *Group) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route {
	routes := make([]*Route, len(methods))
	for i, m := range methods {
		routes[i] = g.Add(m, path, handler, middleware...)
	}
	return routes
}

// Match implements `Nio#Match()` for sub-routes within the Group.
func (g *Group) Match(methods []string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route {
	routes := make([]*Route, len(methods))
	for i, m := range methods {
		routes[i] = g.Add(m, path, handler, middleware...)
	}
	return routes
}

// Group creates a new sub-group with prefix and optional sub-group-level middleware.
func (g *Group) Group(prefix string, middleware ...MiddlewareFunc) *Group {
	m := make([]MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	return g.nio.Group(g.prefix+prefix, m...)
}

// Static implements `Nio#Static()` for sub-routes within the Group.
func (g *Group) Static(prefix, root string) {
	static(g, prefix, root)
}

// File implements `Nio#File()` for sub-routes within the Group.
func (g *Group) File(path, file string) {
	g.nio.File(g.prefix+path, file)
}

// Add implements `Nio#Add()` for sub-routes within the Group.
func (g *Group) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	m := make([]MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	return g.nio.Add(method, g.prefix+path, handler, m...)
}
