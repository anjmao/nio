// +build !go1.11

package mw

import (
	"net/http"
	"net/http/httputil"

	"github.com/anjmao/nio"
)

func proxyHTTP(t *ProxyTarget, c nio.Context, config ProxyConfig) http.Handler {
	return httputil.NewSingleHostReverseProxy(t.URL)
}
