package gdit

type Container interface {
	AddProvider(k string, p any, isNamed bool)
	GetProvider(k string, isNamed bool) (any, bool)
	getLogger() Logger
	CurState() LifeState
	addStartHook(f StartFunc)
	addStopHook(f StopFunc)
}

type CtorFunc[T any] func(ctx InvokeCtx) (T, error)
type HookFunc func(ctx Context) error

type StartFunc func(startCtx StartCtx) error
type StopFunc func(stopCtx StopCtx) error
