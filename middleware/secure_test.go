package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dostack/nio"
	"github.com/stretchr/testify/assert"
)

func TestSecure(t *testing.T) {
	e := nio.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := func(c nio.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Default
	Secure()(h)(c)
	assert.Equal(t, "1; mode=block", rec.Header().Get(nio.HeaderXXSSProtection))
	assert.Equal(t, "nosniff", rec.Header().Get(nio.HeaderXContentTypeOptions))
	assert.Equal(t, "SAMEORIGIN", rec.Header().Get(nio.HeaderXFrameOptions))
	assert.Equal(t, "", rec.Header().Get(nio.HeaderStrictTransportSecurity))
	assert.Equal(t, "", rec.Header().Get(nio.HeaderContentSecurityPolicy))

	// Custom
	req.Header.Set(nio.HeaderXForwardedProto, "https")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	SecureWithConfig(SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	})(h)(c)
	assert.Equal(t, "", rec.Header().Get(nio.HeaderXXSSProtection))
	assert.Equal(t, "", rec.Header().Get(nio.HeaderXContentTypeOptions))
	assert.Equal(t, "", rec.Header().Get(nio.HeaderXFrameOptions))
	assert.Equal(t, "max-age=3600; includeSubdomains", rec.Header().Get(nio.HeaderStrictTransportSecurity))
	assert.Equal(t, "default-src 'self'", rec.Header().Get(nio.HeaderContentSecurityPolicy))
}
