package gdit

type Container interface {
	addInvoke(f HookFunc)
	addProvider(k string, p any, isNamed bool)
	getProvider(k string, isNamed bool) (any, bool)
	init(ctx *Context) error
	start(ctx *Context) error
	stop(ctx *Context) error
}

type CtorFunc[T any] func(ctx *Context) (T, error)
type HookFunc func(ctx *Context) error
