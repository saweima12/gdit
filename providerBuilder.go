package gdit

import "github.com/saweima12/gdit/internal/utils"

const (
	lazy = iota
	factory
	value
)

type ProviderBuilder[T any] interface {
	// When determines whether to register the provider based on a static boolean condition.
	When(condition bool) ProviderBuilder[T]
	// WhenFunc determines whether to register the provider based on a dynamic condition evaluated at runtime.
	WhenFunc(condition func() bool) ProviderBuilder[T]
	// WithName assigns a unique name to the provider for named dependency resolution.
	WithName(name string) ProviderBuilder[T]
	// Attach adds the configured provider to the specified container.
	Attach(c Container)
}

func newProviderBuilder[T any](bType uint8) *providerBuilder[T] {
	return &providerBuilder[T]{
		buildType: bType,
		condition: true,
	}
}

type providerBuilder[T any] struct {
	buildType     uint8
	condition     bool
	conditionFunc func() bool
	name          string
	instance      T
	factory       CtorFunc[T]
}

func (b *providerBuilder[T]) WithName(name string) ProviderBuilder[T] {
	b.name = name
	return b
}

func (b *providerBuilder[T]) When(condition bool) ProviderBuilder[T] {
	b.condition = condition
	return b
}

func (b *providerBuilder[T]) WhenFunc(conditionFunc func() bool) ProviderBuilder[T] {
	b.conditionFunc = conditionFunc
	return b
}

func (b *providerBuilder[T]) Attach(c Container) {
	if !b.shouldRegister() {
		return
	}
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

func (b *providerBuilder[T]) shouldRegister() bool {
	if !b.condition {
		return false
	}
	if b.conditionFunc != nil && !b.conditionFunc() {
		return false
	}
	return true
}
