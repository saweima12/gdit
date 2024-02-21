package gdit

const (
	eager = iota
	lazy
	transient
)

type Builder[T any] interface {
	WithName(name string) Builder[T]
}

type builder[T any] struct {
	bType    uint8
	instance T
	factory  InitFunc[T]
	name     string
}

func (b *builder[T]) WithName(name string) Builder[T] {
	b.name = name
	return b
}

func (b *builder[T]) Build() provider[T] {
	return b.Build()
}

func (b *builder[T]) Attach(c Container) {

}

func (b *builder[T]) getProvider() provider[T] {
	switch b.bType {
	case eager:
		return &eagerProvider[T]{
			instance: b.instance,
		}
	case lazy:
		return &lazyProvider[T]{
			factory: b.factory,
		}
	case transient:
		return &factoryProvider[T]{
			factory: b.factory,
		}
	}
	return nil
}
