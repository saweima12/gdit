package gdit

import (
	"context"
	"fmt"
)

type App struct {
	Container
}

func (app *App) Run() {
	ctx := GetContext(context.Background(), app)

	fmt.Println(ctx)
}
