package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/go-nio/nio"
)

type templateRenderer struct {
	templates *template.Template
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c nio.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// register views
	renderer := &templateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	n := nio.New(nio.WithRenderer(renderer))

	// serve static files under /static path from assets folder
	n.Static("/static", "assets")

	// render index.html
	n.GET("/", func(c nio.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{"name": "Nio!"})
	})

	http.ListenAndServe(":9000", n)
}
