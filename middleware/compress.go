package middleware

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/dostack/dapi"
)

type (
	// GzipConfig defines the config for Gzip middleware.
	GzipConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Gzip compression level.
		// Optional. Default value -1.
		Level int `yaml:"level"`
	}

	gzipResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

const (
	gzipScheme = "gzip"
)

var (
	// DefaultGzipConfig is the default Gzip middleware config.
	DefaultGzipConfig = GzipConfig{
		Skipper: DefaultSkipper,
		Level:   -1,
	}
)

// Gzip returns a middleware which compresses HTTP response using gzip compression
// scheme.
func Gzip() dapi.MiddlewareFunc {
	return GzipWithConfig(DefaultGzipConfig)
}

// GzipWithConfig return Gzip middleware with config.
// See: `Gzip()`.
func GzipWithConfig(config GzipConfig) dapi.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultGzipConfig.Skipper
	}
	if config.Level == 0 {
		config.Level = DefaultGzipConfig.Level
	}

	return func(next dapi.HandlerFunc) dapi.HandlerFunc {
		return func(c dapi.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			res := c.Response()
			res.Header().Add(dapi.HeaderVary, dapi.HeaderAcceptEncoding)
			if strings.Contains(c.Request().Header.Get(dapi.HeaderAcceptEncoding), gzipScheme) {
				res.Header().Set(dapi.HeaderContentEncoding, gzipScheme) // Issue #806
				rw := res.Writer
				w, err := gzip.NewWriterLevel(rw, config.Level)
				if err != nil {
					return err
				}
				defer func() {
					if res.Size == 0 {
						if res.Header().Get(dapi.HeaderContentEncoding) == gzipScheme {
							res.Header().Del(dapi.HeaderContentEncoding)
						}
						// We have to reset response to it's pristine state when
						// nothing is written to body or error is returned.
						// See issue #424, #407.
						res.Writer = rw
						w.Reset(ioutil.Discard)
					}
					w.Close()
				}()
				grw := &gzipResponseWriter{Writer: w, ResponseWriter: rw}
				res.Writer = grw
			}
			return next(c)
		}
	}
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	if code == http.StatusNoContent { // Issue #489
		w.ResponseWriter.Header().Del(dapi.HeaderContentEncoding)
	}
	w.Header().Del(dapi.HeaderContentLength) // Issue #444
	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get(dapi.HeaderContentType) == "" {
		w.Header().Set(dapi.HeaderContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	w.Writer.(*gzip.Writer).Flush()
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *gzipResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
