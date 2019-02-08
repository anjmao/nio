package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/go-nio/nio"
	"github.com/go-nio/nio/mw"
)

type templateRenderer struct {
	templates *template.Template
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c nio.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type dataItem struct {
	Name string `json:"name"`
}

func main() {
	renderer := &templateRenderer{templates: template.Must(template.ParseGlob("dist/*.html"))}
	n := nio.New(nio.WithRenderer(renderer))
	n.Use(mw.Gzip())
	n.Use(mw.CORS())

	// Static files are handled via StaticWithConfig middleware.
	n.Use(mw.StaticWithConfig(mw.StaticConfig{
		Skipper: nio.DefaultSkipper,
		Root:    "./dist",
		Index:   "index.html",
		HTML5:   true,
		Browse:  false,
	}))

	// Some public endpoint for JSON API.
	n.GET("/api/data", func(c nio.Context) error {
		items := []*dataItem{
			&dataItem{Name: "Dog"},
			&dataItem{Name: "Cat"},
			&dataItem{Name: "Tiger"},
		}
		return c.JSON(http.StatusOK, items)
	})

	http.ListenAndServe(":9000", n)
}
