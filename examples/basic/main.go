package main

import (
	"net/http"

	"github.com/go-nio/nio"
)

func main() {
	n := nio.New()

	n.GET("/", func(c nio.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	n.Start(":9000")
}
