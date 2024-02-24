package gdit

import "sync"

type contextPool struct {
	pool sync.Pool
}

func (p *contextPool) Get() *context {
	return p.pool.Get().(*context)
}

func (p *contextPool) Put(ctx *context) {
	ctx.container = nil
	ctx.startHook = nil
	ctx.stopHook = nil
	p.pool.Put(ctx)
}

var ctxPool = contextPool{
	pool: sync.Pool{
		New: func() any {
			return new(context)
		},
	},
}

func getContext(c Container) Context {
	ctx := ctxPool.Get()
	ctx.container = c
	return ctx
}

type LifecycleAdder interface {
	// OnStart registers a hook function to be executed when the application starts.
	// [f] -> The hook function to execute during the application's startup process.
	// This hook allows for custom initialization logic to be executed as part of the startup sequence.
	OnStart(f HookFunc)

	// OnStop registers a hook function to be executed when the application stops.
	// [f] -> The hook function to execute during the application's shutdown process.
	// This hook allows for custom cleanup logic to be executed as part of the shutdown sequence.
	OnStop(f HookFunc)
}

type Context interface {
	LifecycleAdder
	getProvider(key string, isNamed bool) (any, bool)
	clone() Context
	recycle()
	tryRegisterHook()
}

type context struct {
	container Container
	startHook HookFunc
	stopHook  HookFunc
}

func (ctx *context) OnStart(f HookFunc) {
	ctx.startHook = f
}

func (ctx *context) OnStop(f HookFunc) {
	ctx.stopHook = f
}

func (ctx *context) getProvider(key string, isNamed bool) (any, bool) {
	return ctx.container.GetProvider(key, isNamed)
}

func (ctx *context) clone() Context {
	nCtx := ctxPool.Get()
	nCtx.container = ctx.container
	return nCtx
}

func (ctx *context) recycle() {
	ctxPool.Put(ctx)
}

func (ctx *context) tryRegisterHook() {
	if ctx.startHook != nil {
		ctx.container.addStartHook(ctx.startHook)
	}

	if ctx.stopHook != nil {
		ctx.container.addStopHook(ctx.stopHook)
	}
}
