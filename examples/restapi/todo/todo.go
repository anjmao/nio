package todo

import (
	"net/http"

	"github.com/go-nio/nio"
)

// RegisterHandlers registers todo routes and inject todo store
func RegisterHandlers(n *nio.Nio, store TodoStore) {
	h := &handlers{store: store}
	n.GET("/todos", h.GetAllTodos)
	n.GET("/todos/:id", h.GetTodoByID)
	n.POST("/todos", h.AddTodo)
	n.DELETE("/todos/:id", h.DeleteTodo)
}

type handlers struct {
	store TodoStore
}

func (h *handlers) GetAllTodos(c nio.Context) error {
	allTodos, err := h.store.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, allTodos)
}

func (h *handlers) GetTodoByID(c nio.Context) error {
	todoID := c.Param("id")
	todo, err := h.store.GetByID(todoID)
	if err != nil {
		return err
	}
	if todo == nil {
		return nio.ErrNotFound
	}
	return c.JSON(http.StatusOK, todo)
}

func (h *handlers) AddTodo(c nio.Context) error {
	newTodo := &todo{}
	if err := c.Bind(newTodo); err != nil {
		return err
	}
	if err := h.store.Add(newTodo); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *handlers) DeleteTodo(c nio.Context) error {
	todoID := c.Param("id")
	if err := h.store.Delete(todoID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
