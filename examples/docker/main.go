package main

import (
	"log"
	"net/http"

	"github.com/go-nio/nio"
)

func main() {
	// Nio instance
	n := nio.New()

	// Routes
	n.GET("/", hello)

	// Start server
	log.Fatal(http.ListenAndServe(":1323", n))
}

// Handler
func hello(c nio.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
