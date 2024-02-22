package gdit

import (
	"context"
)

type App struct {
	Container
	Scope
	scopes []*Scope
}

func (app *App) Run() error {
	return app.RunWithContext(context.Background())
}

func (app *App) RunWithContext(inputCtx context.Context) error {
	ctx := GetContext(inputCtx, app)

	// Execute all invoke hook.

	// initialize all scope.

	// Execute all start hook.

	return nil
}

func (ap *App) init(ctx *Context) error {
	for i := range ap.initFuncs {
		return ap.initFuncs[i](ctx)
	}
	return nil
}

func (ap *App) stop(ctx *Context) error {
	panic("not implemented") // TODO: Implement
}

func (ap *App) start(ctx *Context) error {
	panic("not implemented") // TODO: Implement
}

func (ap *App) addInvoke(f HookFunc) {
	ap.initFuncs = append(ap.initFuncs, f)
}

func (ap *App) addProvider(k string, p any, isNamed bool) {
	if isNamed {
		ap.namedMap.Store(k, p)
	} else {
		ap.typeMap.Store(k, p)
	}
}

func (ap *App) getProvider(k string, isNamed bool) (val any, ok bool) {
	if isNamed {
		return ap.namedMap.Load(k)
	} else {
		return ap.typeMap.Load(k)
	}
}
