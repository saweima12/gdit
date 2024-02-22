package gdit

import (
	"fmt"
)

func New() *App {
	return &App{}
}

func Inject[T any](ctx Context) (T, error) {
	typeStr := getType[T]()
	return injectInternal[T](ctx, typeStr, false)
}

func InjectNamed[T any](ctx Context, name string) (T, error) {
	return injectInternal[T](ctx, name, true)
}

func MustInject[T any](ctx Context) T {
	typeStr := getType[T]()
	item, err := injectInternal[T](ctx, typeStr, false)
	if err != nil {
		panic(fmt.Sprintf("MustInejct failed, err: %v", err))
	}
	return item
}

func MustInjectNamed[T any](ctx Context, name string) T {
	item, err := injectInternal[T](ctx, name, true)
	if err != nil {
		panic(fmt.Sprintf("MustInejctNamed failed, err: %v", err))
	}
	return item
}

func Invoke[T any](c Container, f CtorFunc[T]) {
	c.addInvoke(func(ctx Context) error {
		_, err := f(ctx)
		return err
	})
}

func InvokeFunc(c Container, f HookFunc) {
	c.addInvoke(f)
}

func Provide[T any](f CtorFunc[T]) Builder[T] {
	return &builder[T]{
		buildType: lazy,
		factory:   f,
	}
}

func ProvideValue[T any](item T) Builder[T] {
	return &builder[T]{
		buildType: value,
		instance:  item,
	}
}

func ProvideFactory[T CtorFunc[any]](f CtorFunc[T]) Builder[T] {
	return &builder[T]{
		buildType: factory,
		factory:   f,
	}
}

func injectInternal[T any](ctx Context, key string, isNamed bool) (T, error) {
	item, ok := ctx.getProvider(key, isNamed)
	if !ok {
		return empty[T](), fmt.Errorf("The key %s is not found.", key)
	}

	p, ok := item.(provider[T])
	if !ok {
		return empty[T](), fmt.Errorf("The item %s is not a valid provider.", key)
	}
	instance, err := p.Get(ctx)
	if err != nil {
		return empty[T](), fmt.Errorf("Get %s failed, err: %v", key, err)
	}
	return instance, nil
}
