# GDIT - Go Dependency Injection Toolkit

A toolkit aimed at becoming a graceful solution for dependency injection in Golang, offering type safety and flexible interfaces based on generics. It hopes to assist people in building applications that are easy to extend and manage.

## Features

- **Type Safety**: Utilizes Go's generics to offer type-safe dependency injection, strengthening code safety and eliminating unnecessary type assertions.
- **No "Magic"**: Employs reflection solely for type acquisition, avoiding extra function parsing and complex abstractions to prioritize ease of understanding.
- **High Flexibility**: Provides dependency injection, lifecycle hooks, and triggers, all of which are optional and not mandatory, aiming to allow users to customize their own experience.
- **Simplicity**: Maintains straightforward implementations with the goal of making underlying concepts easy to understand, reducing developer burden and preventing errors through interface design.
- **Toolkit, Not a Framework**: Focuses on offering tools and methodologies without the need for convention over configuration, facilitating integration into various services.

## Getting Started

GDIT requires Go 1.18 or later. Start by installing GDIT in your project:

```sh
go get github.com/saweima12/gdit
```

### Minial Example

- Use `gdit.New()` to create a container.
- Provide dependencies with `gdit.Provide[T]()` 
- Initialize function with dependency injection using `gdit.Invoke[T]()`.
- Call `app.Startup()` to execute the Start hooks registered during the dependency injection process.
- Upon completion, call `app.Teardown()` to execute all registered Stop hooks.

These are the basic building blocks.

```go
package main

import (
	"github.com/saweima12/gdit"
)

func main() {
	app := gdit.New()
	// Provide dependencies
	// Invoke function

	// Startup will trigger all registered start hooks.
	app.Startup()

	// Do whatever you want... like running an HTTP server.

	// Teardown will trigger all registered stop hooks.
	app.Teardown()
}
```


### Basic Example

#### Define struct and factory function.

```go
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

	// Use the MustInject method to obtain dependencies provided through Provide.
	// Or
	// cfg, err := gdit.Inject[*TestCfg](ctx)
	cfg := gdit.MustInject[*TestCfg](ctx)

	ctx.OnStart(func(startCtx gdit.StartCtx) error {
		fmt.Println("TestRepo OnStart", cfg)
		return nil
	})

	ctx.OnStop(func(stopCtx gdit.StopCtx) error {
		fmt.Println("TestRepo OnStop")
		return nil
	})

	return &testRepo{
		cfg: cfg,
	}, nil
}
```

#### Create a Container and use `Provide()` & `Invoke()`
```go
func main() {
	app := gdit.New()

	// Create a provider, assign it an instance, and then attach it to the container.
	gdit.ProvideValue[*TestCfg](&TestCfg{DomainURL: "http://example.com"}).
		Attach(app)

	// Or use InvokeProvide to register the value as a dependency under the TestRepo interface after initialization.
	repo, err := gdit.InvokeProvide[TestRepo](app, NewTestRepo)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Invoke a function to initalize.
	gdit.InvokeFunc(app, TestInvoke)

	// 
	fmt.Println(repo.GetDomainURL())

	// Startup will trigger all registered start hooks.
	app.Startup()

	//
	app.Teardown()
}
```

#### Injection 
- Within the `InvokeFunction` (which has the `invokeCtx` parameter), use `Inject[T](ctx)`, `InjectNamed[T](ctx, name)`, `MustInject[T](ctx)`, and `MustInjectNamed[T](ctx, name)` to inject providers.

```go
func TestInvoke(ctx gdit.InvokeCtx) error {
	repo, err := gdit.Inject[TestRepo](ctx)
	if err != nil {
		return nil
	}
	fmt.Println("TestInvoke", repo.GetDomainURL())

	return nil
}
```


#### Lifecycle 
- Within the `InvokeFunction`, `InvokeCtx` provides two methods, `OnStart()` and `OnStop()`, to register lifecycle hooks.
```go
ctx.OnStart(func(startCtx gdit.StartCtx) error {
    fmt.Println("TestRepo OnStart", cfg)
    return nil
})

ctx.OnStop(func(stopCtx gdit.StopCtx) error {
    fmt.Println("TestRepo OnStop")
    return nil
})
```

The complete code, combining all the elements mentioned above, is as follows:

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

	// Use the MustInject method to obtain dependencies provided through Provide.
	// Or
	// cfg, err := gdit.Inject[*TestCfg](ctx)
	cfg := gdit.MustInject[*TestCfg](ctx)

	ctx.OnStart(func(startCtx gdit.StartCtx) error {
		fmt.Println("TestRepo OnStart", cfg)
		return nil
	})

	ctx.OnStop(func(stopCtx gdit.StopCtx) error {
		fmt.Println("TestRepo OnStop")
		return nil
	})

	return &testRepo{
		cfg: cfg,
	}, nil
}

func main() {
	app := gdit.New()

	// Create a provider, assign it an instance, and then attach it to the container.
	gdit.ProvideValue[*TestCfg](&TestCfg{DomainURL: "http://example.com"}).
		Attach(app)

	// Or use InvokeProvide to register the value as a dependency under the TestRepo interface after initialization.
	repo, err := gdit.InvokeProvide[TestRepo](app, NewTestRepo)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Invoke a function to initalize.
	gdit.InvokeFunc(app, TestInvoke)

	// Print the domain URL.
	fmt.Println("Main:", repo.GetDomainURL())

	// Startup will trigger all registered start hooks.
	app.Startup()

	//
	app.Teardown()
}

func TestInvoke(ctx gdit.InvokeCtx) error {
	repo, err := gdit.Inject[TestRepo](ctx)
	if err != nil {
		return nil
	}
	fmt.Println("TestInvoke", repo.GetDomainURL())

	return nil
}
```

Result:
```
$ go run ./main.go
TestInvoke http://example.com
Main: http://example.com
TestRepo OnStart &{http://example.com}
TestRepo OnStop
```


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

