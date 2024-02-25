# GDIT - Go Dependency Injection Toolkit

A dependency injection toolkit leveraging Golang Generics, GDIT provides a type-safe and flexible way to manage dependencies in Go applications. Designed to simplify dependency injection with minimal overhead, GDIT empowers developers to build scalable and maintainable Go applications.

## Features

- **Developer-Centric Design**: Prioritizes seamless developer experience with intuitive APIs and tools for swift integration and business logic focus.
- **Generic Support**: Utilizes Go's generics for type-safe dependency injection and resolution, enhancing code robustness without runtime type assertions.
- **Emphasis on Flexibility**: Offers control over service lifecycle management and logging details, allowing customization to fit project-specific needs.
- **Avoiding "Magic"**: Favors simplicity and transparency by steering clear of excessive use of reflection, making it easier to understand and maintain.
- **Focused on Core Functionality**: Delivers essential DI container services, utilities, and lifecycle hooks, streamlining efficiency without unnecessary abstractions.

## Getting Started

GDIT requires Go 1.18 or later. Start by installing GDIT in your project:

```sh
go get github.com/saweima12/gdit
```

### Minial Example

Use gdit.New() to create a basic container, and finally execute the registered start hooks with app.Startup(). This forms the most fundamental block.

```go
package main

import (
	"github.com/saweima12/gdit"
)

func main() {
	app := gdit.New()

	// Startup will trigger all registered start hooks.
	app.Startup()
}
```


### Basic Example
```go
package main

import (
	"fmt"

	"github.com/saweima12/gdit"
)

// Define a configuration struct.
type TestCfg struct {
	DomainURL string
}

// Define a repository interface that depends on the TestCfg.
type TestRepo interface {
	GetDomainURL() string
}

type testRepo struct {
	cfg *TestCfg
}

func (te *testRepo) GetDomainURL() string {
	return te.cfg.DomainURL
}

// Create a factory method to register a lazy provider.
func NewTestRepo(ctx gdit.InvokeCtx) (TestRepo, error) {
	cfg := gdit.MustInject[*TestCfg](ctx)

	ctx.OnStart(func(startCtx gdit.StartCtx) error {
		fmt.Println("Hi", cfg)
		return nil
	})

	return &testRepo{
		cfg: cfg,
	}, nil
}

func main() {
	app := gdit.New().SetLogLevel(gdit.LOG_DEBUG)

	// Attach a valueProvider to the container.
	gdit.ProvideValue[*TestCfg](&TestCfg{DomainURL: "http://example.com"}).
		Attach(app)
	// InvokeProvide will trigger the factory function and register the provider in the container.
	repo, err := gdit.InvokeProvide[TestRepo](app, NewTestRepo)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Print the domain URL.
	fmt.Println(repo.GetDomainURL())

	// Startup will trigger all registered start hooks.
	app.Startup()
}

```



## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

