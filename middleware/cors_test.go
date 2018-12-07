package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dostack/dapi"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	e := dapi.New()

	// Wildcard origin
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := CORS()(dapi.NotFoundHandler)
	h(c)
	assert.Equal(t, "*", rec.Header().Get(dapi.HeaderAccessControlAllowOrigin))

	// Allow origins
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = CORSWithConfig(CORSConfig{
		AllowOrigins: []string{"localhost"},
	})(dapi.NotFoundHandler)
	req.Header.Set(dapi.HeaderOrigin, "localhost")
	h(c)
	assert.Equal(t, "localhost", rec.Header().Get(dapi.HeaderAccessControlAllowOrigin))

	// Preflight request
	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	req.Header.Set(dapi.HeaderOrigin, "localhost")
	req.Header.Set(dapi.HeaderContentType, dapi.MIMEApplicationJSON)
	cors := CORSWithConfig(CORSConfig{
		AllowOrigins:     []string{"localhost"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h = cors(dapi.NotFoundHandler)
	h(c)
	assert.Equal(t, "localhost", rec.Header().Get(dapi.HeaderAccessControlAllowOrigin))
	assert.NotEmpty(t, rec.Header().Get(dapi.HeaderAccessControlAllowMethods))
	assert.Equal(t, "true", rec.Header().Get(dapi.HeaderAccessControlAllowCredentials))
	assert.Equal(t, "3600", rec.Header().Get(dapi.HeaderAccessControlMaxAge))

	// Preflight request with `AllowOrigins` *
	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	req.Header.Set(dapi.HeaderOrigin, "localhost")
	req.Header.Set(dapi.HeaderContentType, dapi.MIMEApplicationJSON)
	cors = CORSWithConfig(CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h = cors(dapi.NotFoundHandler)
	h(c)
	assert.Equal(t, "localhost", rec.Header().Get(dapi.HeaderAccessControlAllowOrigin))
	assert.NotEmpty(t, rec.Header().Get(dapi.HeaderAccessControlAllowMethods))
	assert.Equal(t, "true", rec.Header().Get(dapi.HeaderAccessControlAllowCredentials))
	assert.Equal(t, "3600", rec.Header().Get(dapi.HeaderAccessControlMaxAge))
}
