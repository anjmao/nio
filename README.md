# nio



### Example

```go
package main

import (
	"net/http"
	"log"
	"github.com/anjmao/nio"
	"github.com/anjmao/nio/mw"
)

func main() {
	// Nio instance
	n := nio.New()

	// Middleware
	n.Use(mw.Logger())
	n.Use(mw.Recover())

	// Routes
	n.GET("/", hello)

	// Start server
	log.Fatal(n.Start(":1323"))
}

// Handler
func hello(c nio.Context) error {
	retueturn c.String(http.StatusOK, "Hello, World!")
}
```
