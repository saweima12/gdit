package gdit

import "github.com/saweima12/gdit/utils"

const (
	lazy = iota
	factory
	value
)

type ProviderBuilder[T any] interface {
	// WithName assigns a unique name to the provider for named dependency resolution.
	WithName(name string) ProviderBuilder[T]
	// Attach adds the configured provider to the specified container.
	Attach(c Container)
}

type providerBuilder[T any] struct {
	buildType uint8
	name      string
	instance  T
	factory   CtorFunc[T]
}

func (b *providerBuilder[T]) WithName(name string) ProviderBuilder[T] {
	b.name = name
	return b
}

func (b *providerBuilder[T]) Attach(c Container) {
	p := b.getProvider()
	c.AddProvider(p.Key(), p, p.IsNamed())
}

func (b *providerBuilder[T]) getProvider() provider[T] {
	key, named := utils.GetProviderKey[T](b.name)
	switch b.buildType {
	case value:
		return &valueProvider[T]{
			instance:     b.instance,
			baseProvider: baseProvider{named: named, key: key},
		}
	case lazy:
		return &lazyProvider[T]{
			factory:      b.factory,
			baseProvider: baseProvider{named: named, key: key},
		}
	case factory:
		return &factoryProvider[T]{
			factory:      b.factory,
			baseProvider: baseProvider{named: named, key: key},
		}
	}
	return nil
}
