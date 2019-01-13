package mw

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-nio/nio"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	e := nio.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	f := func(u, p string, c nio.Context) (bool, error) {
		if u == "joe" && p == "secret" {
			return true, nil
		}
		return false, nil
	}
	h := BasicAuth(f)(func(c nio.Context) error {
		return c.String(http.StatusOK, "test")
	})

	assert := assert.New(t)

	// Valid credentials
	auth := basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:secret"))
	req.Header.Set(nio.HeaderAuthorization, auth)
	assert.NoError(h(c))

	h = BasicAuthWithConfig(BasicAuthConfig{
		Skipper:   nil,
		Validator: f,
		Realm:     "someRealm",
	})(func(c nio.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Valid credentials
	auth = basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:secret"))
	req.Header.Set(nio.HeaderAuthorization, auth)
	assert.NoError(h(c))

	// Case-insensitive header scheme
	auth = strings.ToUpper(basic) + " " + base64.StdEncoding.EncodeToString([]byte("joe:secret"))
	req.Header.Set(nio.HeaderAuthorization, auth)
	assert.NoError(h(c))

	// Invalid credentials
	auth = basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:invalid-password"))
	req.Header.Set(nio.HeaderAuthorization, auth)
	he := h(c).(*nio.HTTPError)
	assert.Equal(http.StatusUnauthorized, he.Code)
	assert.Equal(basic+` realm="someRealm"`, res.Header().Get(nio.HeaderWWWAuthenticate))

	// Missing Authorization header
	req.Header.Del(nio.HeaderAuthorization)
	he = h(c).(*nio.HTTPError)
	assert.Equal(http.StatusUnauthorized, he.Code)

	// Invalid Authorization header
	auth = base64.StdEncoding.EncodeToString([]byte("invalid"))
	req.Header.Set(nio.HeaderAuthorization, auth)
	he = h(c).(*nio.HTTPError)
	assert.Equal(http.StatusUnauthorized, he.Code)
}
