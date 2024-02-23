package gdit

type Context interface {
	OnStart(f HookFunc)
	OnStop(f HookFunc)
	getProvider(key string, isNamed bool) (any, bool)
	clone() Context
	tryRegisterHook()
}

type context struct {
	container Container
	startHook HookFunc
	stopHook  HookFunc
}

func GetContext(c Container) Context {
	return &context{
		container: c,
	}
}

func (ctx *context) OnStart(f HookFunc) {
	ctx.startHook = f
}

func (ctx *context) OnStop(f HookFunc) {
	ctx.stopHook = f
}

func (ctx *context) getProvider(key string, isNamed bool) (any, bool) {
	return ctx.container.getProvider(key, isNamed)
}

func (ctx *context) clone() Context {
	return &context{
		container: ctx.container,
	}
}

func (ctx *context) tryRegisterHook() {
	if ctx.startHook != nil {
		ctx.container.addStartHook(ctx.startHook)
	}

	if ctx.stopHook != nil {
		ctx.container.addStopHook(ctx.stopHook)
	}
}
