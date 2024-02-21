package gdit

type Container interface{}

type InitFunc[T any] func(ctx *Context) (T, error)
