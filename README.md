# dapi



### Example

```go
package main

import (
	"net/http"

	"github.com/dostack/dapi"
	"github.com/dostack/dapi/middleware"
)

func main() {
	// Dapi instance
	d := dapi.New()

	// Middleware
	d.Use(middleware.Logger())
	d.Use(middleware.Recover())

	// Routes
	d.GET("/", hello)

	// Start server
	d.Logger.Fatal(d.Start(":1323"))
}

// Handler
func hello(c dapi.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```