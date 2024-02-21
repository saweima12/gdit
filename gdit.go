package gdit

func New() *App {
	return &App{}
}

// func Inject[T any](ctx *Context) (T, error) {
// }
//
// func InjectNamed[T any](ctx *Context, name string) (T, error) {
// }
//
// func MustInject[T any](ctx *Context) T {
// }
//
// func MustInjectNamed[T any](ctx *Context, name string) T {
// }

func ProvideValue[T any](c Container, item T) Builder[T] {
	return &builder[T]{
		bType:    eager,
		instance: item,
	}
}

func ProvideLazy[T InitFunc[any]](f InitFunc[T]) Builder[T] {
	return &builder[T]{
		bType:   lazy,
		factory: f,
	}
}

func ProvideFactory[T InitFunc[any]](f InitFunc[T]) Builder[T] {
	return &builder[T]{
		bType:   transient,
		factory: f,
	}
}
