package gdit

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Scope struct {
	parent     Container
	Name       string
	State      LifeState
	Logger     *loggerWrapper
	mu         sync.RWMutex
	TypeMap    sync.Map
	NamedMap   sync.Map
	startHooks []StartFunc
	stopHooks  []StopFunc
}

func (sc *Scope) getLogger() Logger {
	return sc.Logger
}

func (sc *Scope) AddProvider(k string, p any, isNamed bool) {
	if isNamed {
		sc.storeProvider(k, p, isNamed, &sc.NamedMap)
		sc.Logger.Debug("[%s] -> The provider [%s] is registered by name", sc.Name, k)
	} else {
		sc.storeProvider(k, p, isNamed, &sc.TypeMap)
		sc.Logger.Debug("[%s] -> The provider [%s] is registered by type", sc.Name, k)
	}
}

func (sc *Scope) storeProvider(k string, p any, isNamed bool, providerMap *sync.Map) {
	if _, loaded := providerMap.Swap(k, p); loaded {
		msg := fmt.Sprintf("[%s] -> The provider [%s] was overwritten.", sc.Name, k)
		sc.Logger.Warn(msg)
	}
}

func (sc *Scope) GetProvider(k string, isNamed bool) (val any, ok bool) {
	if isNamed {
		if val, ok := sc.NamedMap.Load(k); ok {
			return val, ok
		}
	} else {
		if val, ok := sc.TypeMap.Load(k); ok {
			return val, ok
		}
	}

	if sc.parent != nil {
		return sc.parent.GetProvider(k, isNamed)
	} else {
		return nil, false
	}
}

func (sc *Scope) CurState() LifeState {
	return sc.State
}

func (sc *Scope) start(ctx StartCtx) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.Logger.Debug("The scope [%s] is starting initialization.", sc.Name)
	sc.changeState(STATE_INITIALIZING)
	for i := range sc.startHooks {
		if err := sc.startHooks[i](ctx); err != nil {
			return err
		}
	}
	sc.Logger.Debug("The scope [%s] is ready.", sc.Name)
	sc.changeState(STATE_READY)
	return nil
}

func (sc *Scope) stop(ctx Context) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.changeState(STATE_SHUTTING_DOWN)
	for i := len(sc.stopHooks) - 1; i >= 0; i-- {
		if err := sc.stopHooks[i](ctx); err != nil {
			return err
		}
	}
	sc.changeState(STATE_TERMINATED)
	return nil
}

func (sc *Scope) addStartHook(f StartFunc) {
	sc.mu.Lock()
	sc.startHooks = append(sc.startHooks, f)
	sc.mu.Unlock()
}

func (sc *Scope) addStopHook(f StopFunc) {
	sc.mu.Lock()
	sc.stopHooks = append(sc.stopHooks, f)
	sc.mu.Unlock()
}

func (sc *Scope) changeState(newState LifeState) {
	preState := atomic.SwapUint32((*uint32)(&sc.State), uint32(newState))
	sc.Logger.Debug("ChangeState %v to %v", LifeState(preState), newState)
}
