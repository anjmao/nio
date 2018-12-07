package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dostack/dapi"
	"github.com/stretchr/testify/assert"
)

func TestSecure(t *testing.T) {
	e := dapi.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := func(c dapi.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Default
	Secure()(h)(c)
	assert.Equal(t, "1; mode=block", rec.Header().Get(dapi.HeaderXXSSProtection))
	assert.Equal(t, "nosniff", rec.Header().Get(dapi.HeaderXContentTypeOptions))
	assert.Equal(t, "SAMEORIGIN", rec.Header().Get(dapi.HeaderXFrameOptions))
	assert.Equal(t, "", rec.Header().Get(dapi.HeaderStrictTransportSecurity))
	assert.Equal(t, "", rec.Header().Get(dapi.HeaderContentSecurityPolicy))

	// Custom
	req.Header.Set(dapi.HeaderXForwardedProto, "https")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	SecureWithConfig(SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	})(h)(c)
	assert.Equal(t, "", rec.Header().Get(dapi.HeaderXXSSProtection))
	assert.Equal(t, "", rec.Header().Get(dapi.HeaderXContentTypeOptions))
	assert.Equal(t, "", rec.Header().Get(dapi.HeaderXFrameOptions))
	assert.Equal(t, "max-age=3600; includeSubdomains", rec.Header().Get(dapi.HeaderStrictTransportSecurity))
	assert.Equal(t, "default-src 'self'", rec.Header().Get(dapi.HeaderContentSecurityPolicy))
}
