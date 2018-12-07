package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dostack/dapi"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	e := dapi.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := func(c dapi.Context) error {
		return c.String(http.StatusOK, "test")
	}

	rid := RequestIDWithConfig(RequestIDConfig{})
	h := rid(handler)
	h(c)
	assert.Len(t, rec.Header().Get(dapi.HeaderXRequestID), 32)

	// Custom generator
	rid = RequestIDWithConfig(RequestIDConfig{
		Generator: func() string { return "customGenerator" },
	})
	h = rid(handler)
	h(c)
	assert.Equal(t, rec.Header().Get(dapi.HeaderXRequestID), "customGenerator")
}
