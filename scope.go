package gdit

import (
	"fmt"
	"sync"
)

type scope struct {
	parent     Container
	name       string
	state      LifeState
	logger     *loggerWrapper
	mu         sync.RWMutex
	typeMap    sync.Map
	namedMap   sync.Map
	startHooks []HookFunc
	stopHooks  []HookFunc
}

func (sc *scope) AddProvider(k string, p any, isNamed bool) {
	if isNamed {
		sc.storeProvider(k, p, isNamed, &sc.namedMap)
	} else {
		sc.storeProvider(k, p, isNamed, &sc.typeMap)
	}
}

func (sc *scope) storeProvider(k string, p any, isNamed bool, providerMap *sync.Map) {
	if _, loaded := providerMap.Swap(k, p); loaded {
		msg := fmt.Sprintf("[%s] -> The provider [%s] was overwritten.", sc.name, k)
		sc.logger.Warn(msg)
	}
}

func (sc *scope) GetProvider(k string, isNamed bool) (val any, ok bool) {
	if isNamed {
		if val, ok := sc.namedMap.Load(k); ok {
			return val, ok
		}
	} else {
		if val, ok := sc.typeMap.Load(k); ok {
			return val, ok
		}
	}
	return sc.parent.GetProvider(k, isNamed)
}

func (sc *scope) start(ctx Context) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.changeState(STATE_INITIALIZING)
	for i := range sc.startHooks {
		if err := sc.startHooks[i](ctx); err != nil {
			return err
		}
	}
	sc.changeState(STATE_READY)
	return nil
}

func (sc *scope) stop(ctx Context) error {
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

func (sc *scope) addStartHook(f HookFunc) {
	sc.mu.Lock()
	sc.startHooks = append(sc.startHooks, f)
	sc.mu.Unlock()
}

func (sc *scope) addStopHook(f HookFunc) {
	sc.mu.Lock()
	sc.stopHooks = append(sc.stopHooks, f)
	sc.mu.Unlock()
}

func (sc *scope) changeState(state LifeState) {
	sc.state = state
}
