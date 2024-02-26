package gdit

import (
	"sync"
)

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

func getContext(c Container) *context {
	ctx := ctxPool.Get()
	ctx.container = c
	return ctx
}

type LifecycleStarter interface {
	// OnStart registers a hook function to be executed when the application starts.
	// [f] -> The hook function to execute during the application's startup process.
	// This hook allows for custom initialization logic to be executed as part of the startup sequence.
	OnStart(f StartFunc)
}

type LifecycleStoper interface {
	// OnStop registers a hook function to be executed when the application stops.
	// [f] -> The hook function to execute during the application's shutdown process.
	// This hook allows for custom cleanup logic to be executed as part of the shutdown sequence.
	OnStop(f StopFunc)
}

type InvokeCtx interface {
	Context
	LifecycleStarter
	LifecycleStoper
}

type StartCtx interface {
	Context
	LifecycleStoper
}

type StopCtx interface {
	Context
}

type Context interface {
	clone() InvokeCtx
	getProvider(key string, isNamed bool) (any, bool)
	tryAddOrRunHook() error
	recycle()
}

type context struct {
	container Container
	startHook StartFunc
	stopHook  StopFunc
}

func (ctx *context) OnStart(f StartFunc) {
	ctx.startHook = f
}

func (ctx *context) OnStop(f StopFunc) {
	ctx.stopHook = f
}

func (ctx *context) getProvider(key string, isNamed bool) (any, bool) {
	return ctx.container.GetProvider(key, isNamed)
}

func (ctx *context) clone() InvokeCtx {
	nCtx := ctxPool.Get()
	nCtx.container = ctx.container
	return nCtx
}

func (ctx *context) recycle() {
	ctxPool.Put(ctx)
}

func (ctx *context) tryAddOrRunHook() error {
	if ctx.startHook != nil {
		if ctx.container.CurState() == STATE_READY {
			return ctx.startHook(ctx)
		} else {
			ctx.container.addStartHook(ctx.startHook)
		}
	}

	if ctx.stopHook != nil {
		ctx.container.addStopHook(ctx.stopHook)
	}
	return nil
}
