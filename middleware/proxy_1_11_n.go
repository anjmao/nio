// +build !go1.11

package middleware

import (
	"net/http"
	"net/http/httputil"

	"github.com/dostack/dapi"
)

func proxyHTTP(t *ProxyTarget, c dapi.Context, config ProxyConfig) http.Handler {
	return httputil.NewSingleHostReverseProxy(t.URL)
}
