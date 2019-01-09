package home

import (
	"net/http"

	"github.com/go-nio/nio"
)

// RegisterHandlers registers routes and handlers
func RegisterHandlers(n *nio.Nio) {
	h := &handlers{}
	n.GET("/", h.GetAllModules)
}

type handlers struct{}

type module struct {
	Name string `json:"title"`
	Path string `json:"path"`
}

func (h *handlers) GetAllModules(c nio.Context) error {
	articles := []*module{
		&module{Name: "Todos list", Path: "http://localhost:9000/todos"},
		&module{Name: "Current user", Path: "http://localhost:9000/user"},
	}
	return c.JSON(http.StatusOK, articles)
}
