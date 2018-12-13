# Nio

Minimalist Go web framework

## Getting Started

### Prerequisites

You need to have at least go 1.11 installed on you local machine.

### Installing

Install nio package with go get

```
go get -u github.com/go-nio/nio
```

## Built With

* [Go](https://www.golang.org/) - The best programming language in the world

## Contributing

Please read [CONTRIBUTING.md](https://github.com/go-nio/nio/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Billie Thompson** - *Initial work* - [PurpleBooth](https://github.com/PurpleBooth)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc





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
