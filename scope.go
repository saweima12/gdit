package gdit

import "sync"

type Scope struct {
	parent     Container
	mu         sync.RWMutex
	typeMap    sync.Map
	namedMap   sync.Map
	initFuncs  []HookFunc
	startHooks []HookFunc
	stopHooks  []HookFunc
}

func (sc *Scope) addProvider(k string, p any, isNamed bool) {
	if isNamed {
		sc.namedMap.Store(k, p)
	} else {
		sc.typeMap.Store(k, p)
	}
}

func (sc *Scope) addInvoke(f HookFunc) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.initFuncs = append(sc.initFuncs, f)
}

func (sc *Scope) getProvider(k string, isNamed bool) (val any, ok bool) {
	if isNamed {
		return sc.namedMap.Load(k)
	} else {
		return sc.typeMap.Load(k)
	}
}
