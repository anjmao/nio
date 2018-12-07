package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dostack/dapi"
	"github.com/stretchr/testify/assert"
)

func TestAddTrailingSlash(t *testing.T) {
	e := dapi.New()
	req := httptest.NewRequest(http.MethodGet, "/add-slash", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := AddTrailingSlash()(func(c dapi.Context) error {
		return nil
	})
	h(c)

	assert := assert.New(t)
	assert.Equal("/add-slash/", req.URL.Path)
	assert.Equal("/add-slash/", req.RequestURI)

	// With config
	req = httptest.NewRequest(http.MethodGet, "/add-slash?key=value", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = AddTrailingSlashWithConfig(TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	})(func(c dapi.Context) error {
		return nil
	})
	h(c)
	assert.Equal(http.StatusMovedPermanently, rec.Code)
	assert.Equal("/add-slash/?key=value", rec.Header().Get(dapi.HeaderLocation))
}

func TestRemoveTrailingSlash(t *testing.T) {
	e := dapi.New()
	req := httptest.NewRequest(http.MethodGet, "/remove-slash/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := RemoveTrailingSlash()(func(c dapi.Context) error {
		return nil
	})
	h(c)

	assert := assert.New(t)

	assert.Equal("/remove-slash", req.URL.Path)
	assert.Equal("/remove-slash", req.RequestURI)

	// With config
	req = httptest.NewRequest(http.MethodGet, "/remove-slash/?key=value", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = RemoveTrailingSlashWithConfig(TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	})(func(c dapi.Context) error {
		return nil
	})
	h(c)
	assert.Equal(http.StatusMovedPermanently, rec.Code)
	assert.Equal("/remove-slash?key=value", rec.Header().Get(dapi.HeaderLocation))

	// With bare URL
	req = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = RemoveTrailingSlash()(func(c dapi.Context) error {
		return nil
	})
	h(c)
	assert.Equal("", req.URL.Path)
}
