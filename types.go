package gdit

type Container interface {
	AddProvider(k string, p any, isNamed bool)
	GetProvider(k string, isNamed bool) (any, bool)
	addStartHook(f HookFunc)
	addStopHook(f HookFunc)
}

type CtorFunc[T any] func(ctx Context) (T, error)
type HookFunc func(ctx Context) error
