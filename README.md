# nio

### Example

```go
package main

import (
    "net/http"
    "log"
	"github.com/go-nio/nio"
)

func main() {
	// Nio instance
	n := nio.New()

	// Routes
	n.GET("/", hello)

	// Start server
	log.Fatal(n.Start(":1323"))
}

// Handler
func hello(c nio.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```
