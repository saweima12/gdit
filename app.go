package gdit

import (
	"errors"
	"sync"
)

type App struct {
	Container
	Scope
	scopes []*Scope
	stopCh chan struct{}
	once   sync.Once
}

func (app *App) Run() error {
	ctx := GetContext(app)
	// Initialize all invoke funtion.
	if err := app.init(ctx); err != nil {
		return err
	}

	if err := app.start(ctx); err != nil {
		return err
	}

	<-app.stopCh

	app.stop(ctx)

	// When all service shutdown, close the stopCh
	close(app.stopCh)
	return nil
}

func (app *App) Stop() {
	app.once.Do(func() {
		app.stopCh <- struct{}{}
	})
}

func (app *App) SetLogger(l Logger) *App {
	app.logger = l
	return app
}

func (app *App) init(ctx *Context) error {
	for i := range app.initFuncs {
		if err := app.initFuncs[i](ctx); err != nil {
			return err
		}
	}

	for i := range app.scopes {
		if err := app.scopes[i].init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (ap *App) start(ctx *Context) error {
	for i := range ap.startHooks {
		if err := ap.startHooks[i](ctx); err != nil {
			return err
		}
	}
	for i := range ap.scopes {
		if err := ap.scopes[i].start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ap *App) stop(ctx *Context) error {
	errOccurred := false
	for i := len(ap.stopHooks) - 1; i >= 0; i-- {
		if err := ap.stopHooks[i](ctx); err != nil {
			ap.logger.Error("Execute stop hook failed, err:", err)
			errOccurred = true
		}
	}

	for i := range ap.scopes {
		if err := ap.scopes[i].stop(ctx); err != nil {
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
	ap.initFuncs = append(ap.initFuncs, f)
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
