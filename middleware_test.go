package nio

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNioWrapMiddleware(t *testing.T) {
	e := New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	buf := new(bytes.Buffer)
	mw := WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf.Write([]byte("mw"))
			h.ServeHTTP(w, r)
		})
	})
	h := mw(func(c Context) error {
		return c.String(http.StatusOK, "OK")
	})
	if assert.NoError(t, h(c)) {
		assert.Equal(t, "mw", buf.String())
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())
	}
}
