package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anjmao/nio"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	e := nio.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := Recover()(nio.HandlerFunc(func(c nio.Context) error {
		panic("test")
	}))
	h(c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

