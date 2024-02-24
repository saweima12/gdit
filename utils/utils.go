package utils

import (
	"reflect"
)

func GetProviderKey[T any](name string) (key string, named bool) {
	if name != "" {
		return name, true
	}
	return GetType[T](), false
}

func GetType[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}

func Empty[T any]() (t T) {
	return
}
