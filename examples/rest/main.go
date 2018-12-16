package main

import (
	"net/http"

	"github.com/go-nio/nio"
)

func main() {
	// 1. create nio instance
	n := nio.New()

	// 2. create some mock in memory todos store
	store := newTodoStore()

	// 3. defince todo crud rest apis
	n.GET("/todos", func(c nio.Context) error {
		allTodos, err := store.GetAll()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, allTodos)
	})

	n.GET("/todos/:id", func(c nio.Context) error {
		todoID := c.Param("id")
		todo, err := store.GetByID(todoID)
		if err != nil {
			return err
		}
		if todo == nil {
			return nio.ErrNotFound
		}
		return c.JSON(http.StatusOK, todo)
	})

	n.POST("/todos", func(c nio.Context) error {
		newTodo := &todo{}
		if err := c.Bind(newTodo); err != nil {
			return err
		}
		if err := store.AddTodo(newTodo); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})

	n.DELETE("/todos/:id", func(c nio.Context) error {
		todoID := c.Param("id")
		if err := store.DeleteTodo(todoID); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})

	// 4. pass nio to http
	http.ListenAndServe(":9000", n)
}
