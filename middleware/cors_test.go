package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dostack/nio"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	e := nio.New()

	// Wildcard origin
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := CORS()(nio.NotFoundHandler)
	h(c)
	assert.Equal(t, "*", rec.Header().Get(nio.HeaderAccessControlAllowOrigin))

	// Allow origins
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = CORSWithConfig(CORSConfig{
		AllowOrigins: []string{"localhost"},
	})(nio.NotFoundHandler)
	req.Header.Set(nio.HeaderOrigin, "localhost")
	h(c)
	assert.Equal(t, "localhost", rec.Header().Get(nio.HeaderAccessControlAllowOrigin))

	// Preflight request
	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	req.Header.Set(nio.HeaderOrigin, "localhost")
	req.Header.Set(nio.HeaderContentType, nio.MIMEApplicationJSON)
	cors := CORSWithConfig(CORSConfig{
		AllowOrigins:     []string{"localhost"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h = cors(nio.NotFoundHandler)
	h(c)
	assert.Equal(t, "localhost", rec.Header().Get(nio.HeaderAccessControlAllowOrigin))
	assert.NotEmpty(t, rec.Header().Get(nio.HeaderAccessControlAllowMethods))
	assert.Equal(t, "true", rec.Header().Get(nio.HeaderAccessControlAllowCredentials))
	assert.Equal(t, "3600", rec.Header().Get(nio.HeaderAccessControlMaxAge))

	// Preflight request with `AllowOrigins` *
	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	req.Header.Set(nio.HeaderOrigin, "localhost")
	req.Header.Set(nio.HeaderContentType, nio.MIMEApplicationJSON)
	cors = CORSWithConfig(CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h = cors(nio.NotFoundHandler)
	h(c)
	assert.Equal(t, "localhost", rec.Header().Get(nio.HeaderAccessControlAllowOrigin))
	assert.NotEmpty(t, rec.Header().Get(nio.HeaderAccessControlAllowMethods))
	assert.Equal(t, "true", rec.Header().Get(nio.HeaderAccessControlAllowCredentials))
	assert.Equal(t, "3600", rec.Header().Get(nio.HeaderAccessControlMaxAge))
}
