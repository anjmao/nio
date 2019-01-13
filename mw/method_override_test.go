package mw

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-nio/nio"
	"github.com/stretchr/testify/assert"
)

func TestMethodOverride(t *testing.T) {
	e := nio.New()
	m := MethodOverride()
	h := func(c nio.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Override with http header
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	req.Header.Set(nio.HeaderXHTTPMethodOverride, http.MethodDelete)
	c := e.NewContext(req, rec)
	m(h)(c)
	assert.Equal(t, http.MethodDelete, req.Method)

	// Override with form parameter
	m = MethodOverrideWithConfig(MethodOverrideConfig{Getter: MethodFromForm("_method")})
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("_method="+http.MethodDelete)))
	rec = httptest.NewRecorder()
	req.Header.Set(nio.HeaderContentType, nio.MIMEApplicationForm)
	c = e.NewContext(req, rec)
	m(h)(c)
	assert.Equal(t, http.MethodDelete, req.Method)

	// Override with query parameter
	m = MethodOverrideWithConfig(MethodOverrideConfig{Getter: MethodFromQuery("_method")})
	req = httptest.NewRequest(http.MethodPost, "/?_method="+http.MethodDelete, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	m(h)(c)
	assert.Equal(t, http.MethodDelete, req.Method)

	// Ignore `GET`
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(nio.HeaderXHTTPMethodOverride, http.MethodDelete)
	assert.Equal(t, http.MethodGet, req.Method)
}
