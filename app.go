package gdit

import (
	"errors"
	"fmt"
	"sync"

	"github.com/saweima12/gdit/internal/ext"
)

type LifeState uint32

func (ls LifeState) String() string {
	switch ls {
	case STATE_UNINITIALIZED:
		return "UNINITIALIZED"
	case STATE_INITIALIZING:
		return "INITIALIZING"
	case STATE_READY:
		return "READY"
	case STATE_SHUTTING_DOWN:
		return "SHUTTING_DOWN"
	case STATE_TERMINATED:
		return "TERMINATED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(ls))
	}
}

const (
	STATE_UNINITIALIZED LifeState = iota
	STATE_INITIALIZING
	STATE_READY
	STATE_SHUTTING_DOWN
	STATE_TERMINATED
)

type App interface {
	Container

	// Startup initializes and starts the application. It executes all registered OnStart hooks
	// in their respective order. An error is returned if any part of the initialization process fails.
	Startup() error

	// Teardown gracefully stops the application. It executes all registered OnStop hooks
	// in reverse order to ensure proper cleanup. An error is returned if the teardown process encounters issues.
	Teardown() error

	// SetLogger assigns a custom logger to the application for capturing runtime logs.
	// Returns a reference to the App for method chaining.
	SetLogger(logger Logger) App

	// SetLogLevel adjusts the logging level of the application's logger.
	// This controls the verbosity of the application logs at runtime.
	// Returns a reference to the App for method chaining.
	SetLogLevel(level LogLevel) App

	// GetScope retrieves or creates a named scope within the application.
	// Scopes are used to manage service lifecycles and dependencies in a modular fashion.
	// If the scope does not exist, it is created and linked to the application's root container.
	GetScope(scopeName string) Container
	CurState() LifeState
}

type app struct {
	*Scope
	subScopes ext.GSyncMap[*Scope]
	once      sync.Once
}

func createApp() *app {
	return &app{
		Scope: &Scope{
			Name:  "root",
			State: STATE_UNINITIALIZED,
			Logger: &loggerWrapper{
				Level:  LOG_INFO,
				Logger: newStandardLogger(),
			},
		},
	}
}

func (ap *app) Startup() error {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	if ap.State != STATE_UNINITIALIZED {
		return errors.New("The app has been launched.")
	}

	// Create a context and execute all start hooks.
	ctx := getContext(ap)
	ap.changeState(STATE_INITIALIZING)
	if err := ap.start(ctx); err != nil {
		return err
	}
	ap.changeState(STATE_READY)
	return nil
}

func (ap *app) Teardown() error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if ap.State != STATE_READY && ap.State != STATE_INITIALIZING {
		return errors.New("The app has not been launched yet.")
	}

	ap.Logger.Debug("The app is starting teardown.")
	// Create a context and execute all stop hooks.
	ctx := getContext(ap)
	if err := ap.stop(ctx); err != nil {
		return err
	}
	ap.Logger.Debug("The app has been terminated")
	return nil
}

func (ap *app) SetLogger(l Logger) App {
	ap.Logger.Logger = l
	return ap
}

func (ap *app) SetLogLevel(level LogLevel) App {
	ap.Logger.Level = level
	return ap
}

func (ap *app) GetScope(scopeName string) Container {
	s := &Scope{
		parent: ap,
		Name:   scopeName,
		State:  ap.State,
		Logger: ap.Logger,
	}
	if _, loaded := ap.subScopes.Swap(scopeName, s); loaded {
		ap.Logger.Warn("The scope [%s] is overwritten", scopeName)
	}
	return s
}

func (ap *app) start(ctx *context) error {
	ap.Logger.Debug("The app is starting initialization.")

	for i := range ap.startHooks {
		if err := ap.startHooks[i](ctx); err != nil {
			return err
		}
	}

	var err error
	ap.subScopes.Range(func(key string, value *Scope) bool {
		if ferr := value.start(ctx); ferr != nil {
			err = ferr
			return false
		}
		return true
	})

	ap.Logger.Debug("The app is ready.")
	return err
}

func (ap *app) stop(ctx Context) error {
	errOccurred := false
	for i := len(ap.stopHooks) - 1; i >= 0; i-- {
		if err := ap.stopHooks[i](ctx); err != nil {
			ap.Logger.Error("Execute stop hook failed, err:", err)
			errOccurred = true
		}
	}

	ap.subScopes.Range(func(key string, value *Scope) bool {
		if err := value.stop(ctx); err != nil {
			ap.Logger.Error("Execute scope stop hook failed, err:", err)
			errOccurred = true
		}
		return true
	})

	if errOccurred {
		return errors.New("errors occurred during stop process, see logs for details")
	}

	return nil
}
