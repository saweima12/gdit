package gdit

import (
	"reflect"
)

func getProviderKey[T any](name string) (key string, named bool) {
	if name != "" {
		return name, true
	}
	return getType[T](), false
}

func getType[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}

func empty[T any]() (t T) {
	return
}
