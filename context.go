package gdit

import "context"

type Context struct {
	context.Context
	container Container
}

func GetContext(parent context.Context, c Container) *Context {
	return &Context{
		container: c,
		Context:   parent,
	}
}

func (ctx *Context) getProvider(key string, isNamed bool) (any, bool) {
	return ctx.container.getProvider(key, isNamed)
}

func (ctx *Context) OnStart(f HookFunc) {
}

func (ctx *Context) OnStop(f HookFunc) {
}
