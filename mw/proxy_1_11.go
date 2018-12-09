// +build go1.11

package mw

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/dostack/nio"
	"github.com/dostack/nio/log"
)

func proxyHTTP(tgt *ProxyTarget, c nio.Context, config ProxyConfig) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(tgt.URL)
	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		desc := tgt.URL.String()
		if tgt.Name != "" {
			desc = fmt.Sprintf("%s(%s)", tgt.Name, tgt.URL.String())
		}
		log.Errorf("remote %s unreachable, could not forward: %v", desc, err)
		c.Error(nio.NewHTTPError(http.StatusServiceUnavailable))
	}
	proxy.Transport = config.Transport
	return proxy
}
