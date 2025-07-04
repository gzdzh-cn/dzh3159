package main

import (
	"context"
	"dzhgo/internal/cmd"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gzdzh-cn/dzhcore/utility/env"
	"github.com/gzdzh-cn/dzhcore/utility/util"
)

type GreetService struct {
	ctx context.Context
}

func NewGreetService(ctx context.Context) *GreetService {
	gs := &GreetService{
		ctx: ctx,
	}

	return gs
}

func (gs *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}
func (gs *GreetService) Shutdown(ctx context.Context) {
	g.Log().Info(ctx, "GreetService Shutdown")
}

func (gs *GreetService) SetLogger() {
	type Config struct {
		IsProd      bool
		AppName     string
		IsDesktop   bool
		ConfigMap   g.Map
		defaultPath string
	}
	config := &Config{}
	config.defaultPath = env.GetCfgWithDefault(ctx, "core.gfLogger.path", g.NewVar("path")).String()
	config.IsProd = env.GetCfgWithDefault(ctx, "core.isProd", g.NewVar(false)).Bool()
	config.AppName = env.GetCfgWithDefault(ctx, "core.appName", g.NewVar("dzhgo")).String()
	config.IsDesktop = env.GetCfgWithDefault(ctx, "core.isDesktop", g.NewVar(false)).Bool()
	logPath := util.NewToolUtil().GetLoggerPath(config.IsProd, config.AppName, config.IsDesktop, config.defaultPath)
	config.ConfigMap = g.Map{
		"path":     logPath,
		"file":     env.GetCfgWithDefault(ctx, "core.gfLogger.file", g.NewVar("{Y-m-d}.log")).String(),
		"level":    env.GetCfgWithDefault(ctx, "core.gfLogger.level", g.NewVar("debug")).String(),
		"stdout":   env.GetCfgWithDefault(ctx, "core.gfLogger.stdout", g.NewVar(true)).Bool(),
		"flags":    env.GetCfgWithDefault(ctx, "core.gfLogger.flags", g.NewVar(44)).Int(),
		"stStatus": env.GetCfgWithDefault(ctx, "core.gfLogger.stStatus", g.NewVar(1)).Int(),
		"stSkip":   env.GetCfgWithDefault(ctx, "core.gfLogger.stSkip", g.NewVar(0)).Int(),
	}
	g.Log().SetConfigWithMap(config.ConfigMap)
}

func (gs *GreetService) StartGfServer() {
	// 启动 goframe 服务
	go cmd.Main.Run(ctx)
}
