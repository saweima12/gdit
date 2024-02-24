package gdit

import (
	"fmt"

	"github.com/saweima12/gdit/internal/utils"
)

// New creates and returns a new instance of the application with a root container.
// This instance is ready for configuration and startup, allowing service registration,
// lifecycle management, and dependency injection.
func New() App {
	return createApp()
}

// Inject resolves a dependency of type T using the provided context.
// [ctx] -> The context used for dependency resolution, managing service lifecycles and dependencies.
// Returns an instance of type T and any error encountered during resolution.
func Inject[T any](ctx Context) (T, error) {
	typeStr := utils.GetType[T]()
	return injectInternal[T](ctx, typeStr, false)
}

// InjectNamed resolves a named dependency of type T using the provided context.
// [ctx] -> The context used for dependency resolution.
// [name] -> The unique name identifying the dependency to be resolved.
// Returns an instance of type T associated with the given name and any error encountered.
func InjectNamed[T any](ctx Context, name string) (T, error) {
	return injectInternal[T](ctx, name, true)
}

// MustInject resolves a dependency of type T using the provided context. Panics if resolution fails.
// [ctx] -> The context used for dependency resolution.
// Returns an instance of type T. Panics with an error message if the dependency cannot be resolved.
func MustInject[T any](ctx Context) T {
	typeStr := utils.GetType[T]()
	item, err := injectInternal[T](ctx, typeStr, false)
	if err != nil {
		panic(fmt.Sprintf("MustInejct failed, err: %v", err))
	}
	return item
}

// MustInjectNamed resolves a named dependency of type T using the provided context. Panics if resolution fails.
// [ctx] -> The context used for dependency resolution.
// [name] -> The unique name identifying the dependency to be resolved.
// Returns an instance of type T associated with the given name. Panics with an error message if the dependency cannot be resolved.
func MustInjectNamed[T any](ctx Context, name string) T {
	item, err := injectInternal[T](ctx, name, true)
	if err != nil {
		panic(fmt.Sprintf("MustInejctNamed failed, err: %v", err))
	}
	return item
}

// Invoke calls a constructor function with the container's context, used for service initialization.
// [c] -> Container where the function is executed, managing service lifecycles and dependencies.
// [f] -> Constructor function that accepts a Context and returns a service instance (of type T) and an error.
// Returns the service instance and any error encountered during execution.
func Invoke[T any](c Container, f func(Context) (T, error)) (T, error) {
	ctx := getContext(c)
	return f(ctx)
}

// InvokeFunc executes a function within the container's context, for initialization tasks.
// [c] -> Container in which the function is executed.
// [f] -> Function that takes a Context and performs initialization, returning an error if it fails.
// Returns an error if the initialization task fails.
func InvokeFunc(c Container, f func(Context) error) error {
	ctx := getContext(c)
	return f(ctx)
}

// Provide registers a lazy-loaded service constructor within the DI system.
// [f] -> A constructor function that takes a Context and returns an instance of type T and an error.
//
//	This constructor is called lazily, i.e., the service is instantiated when first requested.
//
// Returns a ProviderBuilder to further configure the provided service.
func Provide[T any](f func(Context) (T, error)) ProviderBuilder[T] {
	return &providerBuilder[T]{
		buildType: lazy,
		factory:   f,
	}
}

// ProvideFactory registers a factory-based service constructor within the DI system.
// [f] -> A factory function that takes a Context and returns an instance of type T and an error.
//
//	Unlike lazy-loaded services, factory-based services can be instantiated multiple times.
//
// Returns a ProviderBuilder to further configure the provided service.
func ProvideFactory[T any](f func(Context) (T, error)) ProviderBuilder[T] {
	return &providerBuilder[T]{
		buildType: factory,
		factory:   f,
	}
}

// ProvideValue registers a pre-instantiated service instance within the DI system.
// [item] -> The pre-instantiated service instance of type T to be registered.
// Returns a ProviderBuilder to further configure the provided service.
func ProvideValue[T any](item T) ProviderBuilder[T] {
	return &providerBuilder[T]{
		buildType: value,
		instance:  item,
	}
}

func injectInternal[T any](ctx Context, key string, isNamed bool) (T, error) {
	item, ok := ctx.getProvider(key, isNamed)
	if !ok {
		return utils.Empty[T](), fmt.Errorf("The key %s is not found.", key)
	}

	p, ok := item.(provider[T])
	if !ok {
		return utils.Empty[T](), fmt.Errorf("The item %s is not a valid provider.", key)
	}

	// Clone a independet context
	indCtx := ctx.clone()
	defer indCtx.recycle()
	instance, err := p.Get(indCtx)
	if err != nil {
		return utils.Empty[T](), fmt.Errorf("Get %s failed, err: %v", key, err)
	}

	// try to register hook.
	indCtx.tryRegisterHook()
	return instance, nil
}
