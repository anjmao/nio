package main

import (
	"net/http"

	"github.com/go-nio/nio"
	"github.com/go-nio/nio/examples/restapi/home"
	"github.com/go-nio/nio/examples/restapi/todo"
	"github.com/go-nio/nio/examples/restapi/user"
	"github.com/go-nio/nio/mw"
)

func main() {
	n := nio.New()

	// add middleware
	n.Use(mw.Recover())

	// register home handler
	home.RegisterHandlers(n)

	// register todo handlers and inject todo store
	todoStore := todo.NewTodoStore()
	todo.RegisterHandlers(n, todoStore)

	// register user handlers
	user.RegisterHandlers(n)

	http.ListenAndServe(":9000", n)
}
