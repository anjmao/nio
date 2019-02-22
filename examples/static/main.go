package main

import (
	"flag"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/go-nio/nio"
)

var (
	addr = flag.String("addr", ":9000", "Server serve address")
)

type templateRenderer struct {
	templates *template.Template
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c nio.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	flag.Parse()

	// Register views.
	renderer := &templateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	n := nio.New(nio.WithRenderer(renderer))

	// Serve static files under /static path from assets folder.
	n.Static("/static", "assets")

	// Render index.html.
	n.GET("/", func(c nio.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{"name": "Nio!"})
	})

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         *addr,
		Handler:      n,
	}
	srv.ListenAndServe()
}
