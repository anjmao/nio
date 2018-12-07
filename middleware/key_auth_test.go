package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dostack/dapi"
	"github.com/stretchr/testify/assert"
)

func TestKeyAuth(t *testing.T) {
	e := dapi.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	config := KeyAuthConfig{
		Validator: func(key string, c dapi.Context) (bool, error) {
			return key == "valid-key", nil
		},
	}
	h := KeyAuthWithConfig(config)(func(c dapi.Context) error {
		return c.String(http.StatusOK, "test")
	})

	assert := assert.New(t)

	// Valid key
	auth := DefaultKeyAuthConfig.AuthScheme + " " + "valid-key"
	req.Header.Set(dapi.HeaderAuthorization, auth)
	assert.NoError(h(c))

	// Invalid key
	auth = DefaultKeyAuthConfig.AuthScheme + " " + "invalid-key"
	req.Header.Set(dapi.HeaderAuthorization, auth)
	he := h(c).(*dapi.HTTPError)
	assert.Equal(http.StatusUnauthorized, he.Code)

	// Missing Authorization header
	req.Header.Del(dapi.HeaderAuthorization)
	he = h(c).(*dapi.HTTPError)
	assert.Equal(http.StatusBadRequest, he.Code)

	// Key from custom header
	config.KeyLookup = "header:API-Key"
	h = KeyAuthWithConfig(config)(func(c dapi.Context) error {
		return c.String(http.StatusOK, "test")
	})
	req.Header.Set("API-Key", "valid-key")
	assert.NoError(h(c))

	// Key from query string
	config.KeyLookup = "query:key"
	h = KeyAuthWithConfig(config)(func(c dapi.Context) error {
		return c.String(http.StatusOK, "test")
	})
	q := req.URL.Query()
	q.Add("key", "valid-key")
	req.URL.RawQuery = q.Encode()
	assert.NoError(h(c))

	// Key from form
	config.KeyLookup = "form:key"
	h = KeyAuthWithConfig(config)(func(c dapi.Context) error {
		return c.String(http.StatusOK, "test")
	})
	f := make(url.Values)
	f.Set("key", "valid-key")
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(dapi.HeaderContentType, dapi.MIMEApplicationForm)
	c = e.NewContext(req, rec)
	assert.NoError(h(c))
}
