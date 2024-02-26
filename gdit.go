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
func Invoke[T any](c Container, f func(InvokeCtx) (T, error)) (T, error) {
	ctx := getContext(c)
	defer ctx.recycle()

	resp, err := f(ctx)
	if err != nil {
		return resp, err
	}
	if err := ctx.tryAddOrRunHook(); err != nil {
		return resp, err
	}
	return resp, nil
}

// InvokeProvide combines service initialization with automatic registration in the container.
// [c] -> Container where the function is executed, managing service lifecycles and dependencies.
// [f] -> Constructor function that accepts a Context and returns a service instance (of type T) and an error.
// Returns the service instance and any error encountered during execution.
func InvokeProvide[T any](c Container, f func(InvokeCtx) (T, error)) (T, error) {
	instance, err := Invoke[T](c, f)
	if err != nil {
		return instance, err
	}
	ProvideValue[T](instance).Attach(c)
	return instance, nil
}

// InvokeFunc executes a function within the container's context, for initialization tasks.
// [c] -> Container in which the function is executed.
// [f] -> Function that takes a Context and performs initialization, returning an error if it fails.
// Returns an error if the initialization task fails.
func InvokeFunc(c Container, f func(InvokeCtx) error) error {
	ctx := getContext(c)
	defer ctx.recycle()
	if err := f(ctx); err != nil {
		return err
	}
	if err := ctx.tryAddOrRunHook(); err != nil {
		return err
	}
	return nil
}

// Provide registers a lazy-loaded service constructor within the DI system.
// [f] -> A constructor function that takes a Context and returns an instance of type T and an error.
//
//	This constructor is called lazily, i.e., the service is instantiated when first requested.
//
// Returns a ProviderBuilder to further configure the provided service.
func Provide[T any](f func(InvokeCtx) (T, error)) ProviderBuilder[T] {
	pb := newProviderBuilder[T](provider_lazy)
	pb.factory = f
	return pb
}

// ProvideFactory registers a factory-based service constructor within the DI system.
// [f] -> A factory function that takes a Context and returns an instance of type T and an error.
//
//	Unlike lazy-loaded services, factory-based services can be instantiated multiple times.
//
// Returns a ProviderBuilder to further configure the provided service.
func ProvideFactory[T any](f func(InvokeCtx) (T, error)) ProviderBuilder[T] {
	pb := newProviderBuilder[T](provider_factory)
	pb.factory = f
	return pb
}

// ProvideValue registers a pre-instantiated service instance within the DI system.
// [item] -> The pre-instantiated service instance of type T to be registered.
// Returns a ProviderBuilder to further configure the provided service.
func ProvideValue[T any](item T) ProviderBuilder[T] {
	pb := newProviderBuilder[T](provider_value)
	pb.instance = item
	return pb
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
	err = indCtx.tryAddOrRunHook()
	if err != nil {
		return utils.Empty[T](), fmt.Errorf("Execution of the startup hook for the %s failed. err: %v", key, err)
	}

	return instance, nil
}
