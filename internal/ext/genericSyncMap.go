package ext

import (
	"sync"

	"github.com/saweima12/gdit/internal/utils"
)

type GSyncMap[T any] struct {
	smp sync.Map
}

func (gsm *GSyncMap[T]) Load(k string) (T, bool) {
	resp, ok := gsm.smp.Load(k)
	if !ok {
		return utils.Empty[T](), ok
	}
	return resp.(T), ok
}

func (gsm *GSyncMap[T]) Store(k string, p T) {
	gsm.smp.Store(k, p)
}

func (gsm *GSyncMap[T]) Swap(k string, p T) (prev T, loaded bool) {
	resp, loaded := gsm.smp.Swap(k, p)
	if !loaded {
		return utils.Empty[T](), loaded
	}
	return resp.(T), loaded
}

func (gsm *GSyncMap[T]) Range(f func(key string, value T) bool) {
	gsm.smp.Range(func(key, value any) bool {
		return f(key.(string), value.(T))
	})
}
