package gdit

import "sync"

type Scope struct {
	parent      Container
	state       LifeState
	logger      Logger
	mu          sync.RWMutex
	typeMap     sync.Map
	namedMap    sync.Map
	invokeFuncs []HookFunc
	startHooks  []HookFunc
	stopHooks   []HookFunc
}

func (sc *Scope) init(ctx Context) error {
	for i := range sc.invokeFuncs {
		if err := sc.invokeFuncs[i](ctx); err != nil {
			return err
		}
	}
	return nil
}

func (sc *Scope) start(ctx Context) error {
	for i := range sc.startHooks {
		if err := sc.startHooks[i](ctx); err != nil {
			return err
		}
	}
	return nil
}

func (sc *Scope) stop(ctx Context) error {
	for i := len(sc.stopHooks) - 1; i >= 0; i-- {
		if err := sc.stopHooks[i](ctx); err != nil {
			return err
		}
	}
	return nil
}

func (sc *Scope) addInvoke(f HookFunc) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.invokeFuncs = append(sc.invokeFuncs, f)
}

func (sc *Scope) addStartHook(f HookFunc) {
	sc.startHooks = append(sc.startHooks, f)
}

func (sc *Scope) addStopHook(f HookFunc) {
	sc.stopHooks = append(sc.stopHooks, f)
}

func (sc *Scope) changeState(state LifeState) {
	sc.state = state
}

func (sc *Scope) addProvider(k string, p any, isNamed bool) {
	if isNamed {
		sc.namedMap.Store(k, p)
	} else {
		sc.typeMap.Store(k, p)
	}
}

func (sc *Scope) getProvider(k string, isNamed bool) (val any, ok bool) {
	if isNamed {
		if val, ok := sc.namedMap.Load(k); ok {
			return val, ok
		}
	} else {
		if val, ok := sc.typeMap.Load(k); ok {
			return val, ok
		}
	}
	return sc.parent.getProvider(k, isNamed)
}
