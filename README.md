# nio



### Example

```go
package main

import (
	"net/http"

	"github.com/dostack/nio"
	"github.com/dostack/nio/middleware"
)

func main() {
	// Nio instance
	d := nio.New()

	// Middleware
	d.Use(middleware.Logger())
	d.Use(middleware.Recover())

	// Routes
	d.GET("/", hello)

	// Start server
	d.Logger.Fatal(d.Start(":1323"))
}

// Handler
func hello(c nio.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```
