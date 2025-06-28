package main

import (
	"context"

	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

var (
// Logger = log.NewRunLogger()
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

}

func (a *App) Shutdown(ctx context.Context) {
	g.Log().Info(ctx, "Shutdown")
}

// Greet returns a greeting for the given name
func (a *App) Greet(p string) string {
	// 使用 GoFrame 上下文
	// ctx := gctx.New()
	// 记录日志
	// Logger.Info(fmt.Sprintf("Greeting %s", p))
	return fmt.Sprintf("Hello %s!", p)
}

// func (a *App) GetConfig() string {
// 	var baseSysSetting = &model.BaseSysSetting{}
// 	g.DB().Model("base_sys_setting").Where("id", 1).Scan(&baseSysSetting)
// 	// g.Log().Debugf(a.ctx, "GetConfig: %v", config)
// 	Logger.Info(fmt.Sprintf("GetConfig: %v", baseSysSetting))
// 	logPath := util.GetLoggerPath(config.Cfg.IsProd, config.AppName)

// 	return fmt.Sprintf("Hello %v,logPath:%s,isProd:%v!", baseSysSetting.ID, logPath, config.Cfg.IsProd)
// }
