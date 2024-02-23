package gdit

import (
	"errors"
	"sync"
)

type LifeState uint8

const (
	Uninitialized LifeState = iota
	Initialzing
	Ready
	ShuttingDown
	Terminated
)

type App struct {
	*Scope
	subScopes []*Scope
	stopCh    chan struct{}
	once      sync.Once
}

func createApp() *App {
	return &App{
		Scope: &Scope{
			state:  Uninitialized,
			logger: &standardLogger{},
		},
	}
}

func (app *App) Setup() error {
	app.mu.Lock()
	defer app.mu.Unlock()
	// Create context
	ctx := GetContext(app)

	app.changeState(Initialzing)
	if err := app.init(ctx); err != nil {
		return err
	}

	if err := app.start(ctx); err != nil {
		return err
	}

	app.changeState(Ready)
	return nil
}

func (app *App) Teardown() error {

	ctx := GetContext(app)
	app.stop(ctx)

	return nil
}

func (app *App) SetLogger(l Logger) *App {
	app.logger = l
	return app
}

func (app *App) GetScope(scopeName string) Container {
	return &Scope{
		parent: app,
		state:  app.state,
		logger: app.logger,
	}
}

func (app *App) init(ctx Context) error {
	for i := range app.invokeFuncs {
		if err := app.invokeFuncs[i](ctx); err != nil {
			return err
		}
	}

	for i := range app.subScopes {
		if err := app.subScopes[i].init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (ap *App) start(ctx Context) error {
	for i := range ap.startHooks {
		if err := ap.startHooks[i](ctx); err != nil {
			return err
		}
	}
	for i := range ap.subScopes {
		if err := ap.subScopes[i].start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ap *App) stop(ctx Context) error {
	errOccurred := false
	for i := len(ap.stopHooks) - 1; i >= 0; i-- {
		if err := ap.stopHooks[i](ctx); err != nil {
			ap.logger.Error("Execute stop hook failed, err:", err)
			errOccurred = true
		}
	}

	for i := range ap.subScopes {
		if err := ap.subScopes[i].stop(ctx); err != nil {
			ap.logger.Error("Execute scope stop hook failed, err:", err)
			errOccurred = true
		}
	}

	if errOccurred {
		return errors.New("errors occurred during stop process, see logs for details")
	}

	return nil
}

func (ap *App) addInvoke(f HookFunc) {
	ap.mu.Lock()
	ap.invokeFuncs = append(ap.invokeFuncs, f)
	ap.mu.Unlock()
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

func (ap *App) changeState(state LifeState) {
	ap.state = state
}
