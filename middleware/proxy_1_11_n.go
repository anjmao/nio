// +build !go1.11

package middleware

import (
	"net/http"
	"net/http/httputil"

	"github.com/dostack/nio"
)

func proxyHTTP(t *ProxyTarget, c nio.Context, config ProxyConfig) http.Handler {
	return httputil.NewSingleHostReverseProxy(t.URL)
}
