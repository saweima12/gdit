package gdit

const (
	lazy = iota
	factory
	value
)

type Builder[T any] interface {
	WithName(name string) Builder[T]
	Attach(c Container)
}

type builder[T any] struct {
	buildType uint8
	name      string
	instance  T
	factory   CtorFunc[T]
}

func (b *builder[T]) WithName(name string) Builder[T] {
	b.name = name
	return b
}

func (b *builder[T]) Attach(c Container) {
	p := b.getProvider()
	c.addProvider(p.Key(), p, p.IsNamed())
}

func (b *builder[T]) getProvider() provider[T] {
	key, named := getProviderKey[T](b.name)
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
