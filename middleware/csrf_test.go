package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dostack/nio"
	"github.com/dostack/nio/internal/random"
	"github.com/stretchr/testify/assert"
)

func TestCSRF(t *testing.T) {
	e := nio.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	csrf := CSRFWithConfig(CSRFConfig{
		TokenLength: 16,
	})
	h := csrf(func(c nio.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Generate CSRF token
	h(c)
	assert.Contains(t, rec.Header().Get(nio.HeaderSetCookie), "_csrf")

	// Without CSRF cookie
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	assert.Error(t, h(c))

	// Empty/invalid CSRF token
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	req.Header.Set(nio.HeaderXCSRFToken, "")
	assert.Error(t, h(c))

	// Valid CSRF token
	token := random.String(16)
	req.Header.Set(nio.HeaderCookie, "_csrf="+token)
	req.Header.Set(nio.HeaderXCSRFToken, token)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCSRFTokenFromForm(t *testing.T) {
	f := make(url.Values)
	f.Set("csrf", "token")
	e := nio.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Add(nio.HeaderContentType, nio.MIMEApplicationForm)
	c := e.NewContext(req, nil)
	token, err := csrfTokenFromForm("csrf")(c)
	if assert.NoError(t, err) {
		assert.Equal(t, "token", token)
	}
	_, err = csrfTokenFromForm("invalid")(c)
	assert.Error(t, err)
}

func TestCSRFTokenFromQuery(t *testing.T) {
	q := make(url.Values)
	q.Set("csrf", "token")
	e := nio.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(nio.HeaderContentType, nio.MIMEApplicationForm)
	req.URL.RawQuery = q.Encode()
	c := e.NewContext(req, nil)
	token, err := csrfTokenFromQuery("csrf")(c)
	if assert.NoError(t, err) {
		assert.Equal(t, "token", token)
	}
	_, err = csrfTokenFromQuery("invalid")(c)
	assert.Error(t, err)
	csrfTokenFromQuery("csrf")
}
