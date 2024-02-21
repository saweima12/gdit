package gdit

import "sync"

type provider[T any] interface {
	Get(ctx *Context) (T, error)
}

type eagerProvider[T any] struct {
	instance T
}

func (p *eagerProvider[T]) Get(ctx *Context) (T, error) {
	return p.instance, nil
}

type lazyProvider[T any] struct {
	instance T
	factory  InitFunc[T]
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
	factory InitFunc[T]
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
