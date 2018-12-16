package main

import (
	"net/http"

	"github.com/go-nio/nio"
)

func main() {
	// 1. create  nio instance
	n := nio.New()

	// 2. add some routes
	n.GET("/", func(c nio.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	// 3. pass nio to http
	http.ListenAndServe(":9000", n)
}
