package main

import (
	"net/http"

	"github.com/go-nio/nio"
	"github.com/go-nio/nio/examples/restapi/home"
	"github.com/go-nio/nio/examples/restapi/todo"
	"github.com/go-nio/nio/examples/restapi/user"
)

func main() {
	n := nio.New()

	// register home handler
	home.RegisterHandlers(n)

	// register todo handlers and inject todo store
	todoStore := todo.NewTodoStore()
	todo.RegisterHandlers(n, todoStore)

	// register user handlers
	user.RegisterHandlers(n)

	http.ListenAndServe(":9000", n)
}
