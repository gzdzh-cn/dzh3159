package main

import (
	"context"
	"dzhgo/internal/cmd"

	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	go cmd.Main.Run(ctx)

}

func (a *App) Shutdown(ctx context.Context) {
	g.Log().Info(ctx, "Shutdown")
}

// Greet returns a greeting for the given name
func (a *App) Greet(p string) string {
	// 使用 GoFrame 上下文
	// ctx := gctx.New()
	// 记录日志

	return fmt.Sprintf("Hello %s!", p)
}
