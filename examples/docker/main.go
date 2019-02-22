package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/go-nio/nio"
)

var (
	addr = flag.String("addr", ":9000", "Server serve address")
)

func main() {
	flag.Parse()

	n := nio.New()
	n.GET("/", hello)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         *addr,
		Handler:      n,
	}
	srv.ListenAndServe()
}

// Handler
func hello(c nio.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
