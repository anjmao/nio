package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/go-nio/nio"
	"github.com/go-nio/nio/examples/restapi/home"
	"github.com/go-nio/nio/examples/restapi/todo"
	"github.com/go-nio/nio/examples/restapi/user"
	"github.com/go-nio/nio/mw"
)

var (
	addr = flag.String("addr", ":9000", "Server serve address")
)

func main() {
	flag.Parse()

	n := nio.New()

	// Add middleware.
	n.Use(mw.Recover())

	// Register home handler.
	home.RegisterHandlers(n)

	// Register todo handlers and inject todo store.
	todoStore := todo.NewTodoStore()
	todo.RegisterHandlers(n, todoStore)

	// Register user handlers.
	user.RegisterHandlers(n)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         *addr,
		Handler:      n,
	}
	srv.ListenAndServe()
}
