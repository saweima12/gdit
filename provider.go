package gdit

import "sync"

type provider[T any] interface {
	Get(ctx *Context) (T, error)
	IsNamed() bool
	Key() string
}

type baseProvider struct {
	key   string
	named bool
}

func (p *baseProvider) IsNamed() bool {
	return p.named
}

func (p *baseProvider) Key() string {
	return p.key
}

type valueProvider[T any] struct {
	baseProvider
	instance T
}

func (p *valueProvider[T]) Get(ctx *Context) (T, error) {
	return p.instance, nil
}

type lazyProvider[T any] struct {
	baseProvider
	instance T
	factory  CtorFunc[T]
	once     sync.Once
}

func (p *lazyProvider[T]) Get(ctx *Context) (T, error) {
	var err error
	p.once.Do(func() {
		instance, ferr := p.factory(ctx)
		if ferr != nil {
			err = ferr
		}
		p.instance = instance
	})
	return p.instance, err
}

type factoryProvider[T any] struct {
	baseProvider
	factory CtorFunc[T]
	scoped  uint
}

func (p *factoryProvider[T]) Get(ctx *Context) (T, error) {
	instance, err := p.factory(ctx)
	if err != nil {
		var zero T
		return zero, err
	}
	return instance, nil
}
