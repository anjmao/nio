package user

import (
	"net/http"

	"github.com/go-nio/nio"
)

// RegisterHandlers registers routes and handlers
func RegisterHandlers(n *nio.Nio) {
	h := &handlers{}
	n.GET("/user", h.GetCurrentUser)
}

type handlers struct{}

type user struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (h *handlers) GetCurrentUser(c nio.Context) error {
	u := &user{
		ID:        "u1",
		FirstName: "Apolo",
		LastName:  "Omega",
	}
	return c.JSON(http.StatusOK, u)
}
