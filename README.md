# GDIT - Go Dependency Injection Toolkit

A dependency injection toolkit based on Golang Generics.

leveraging the new Generics feature to provide a type-safe and flexible way to manage dependencies in  Go applications. 


## Getting Started

To get started with GDIT, ensure you have a Go version that supports Generics (Go 1.18 or later). Then, follow these simple steps:

### installation
```sh
go get github.com/saweima12/gdit
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

