package log

import (
	"context"
	"os"

	"dzhgo/internal/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gzdzh-cn/dzhcore/utility/util"
)

type SRunLogger struct {
	gLog *glog.Logger
	ctx  context.Context
}

func NewRunLogger() *SRunLogger {

	isProd := config.Cfg.IsProd
	gLog := glog.New()
	logPath := util.GetLoggerPath(isProd, config.AppName)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.MkdirAll(logPath, 0755)
	}
	gLog.SetConfigWithMap(g.Map{
		"path":     logPath,
		"level":    "all",
		"stdout":   true,
		"flags":    44,
		"StStatus": 1,
		"StSkip":   1,
	})
	logger := &SRunLogger{
		gLog: gLog,
		ctx:  context.TODO(),
	}
	return logger
}

// Print works like Sprintf.
func (l *SRunLogger) Print(message string) {
	l.gLog.Print(l.ctx, message)
}

// Trace level logging. Works like Sprintf.
func (l *SRunLogger) Trace(message string) {
	l.gLog.Error(l.ctx, message)
}

// Debug level logging. Works like Sprintf.
func (l *SRunLogger) Debug(message string) {
	l.gLog.Debug(l.ctx, message)
}

// Info level logging. Works like Sprintf.
func (l *SRunLogger) Info(message string) {
	l.gLog.Info(l.ctx, message)
}

// Warning level logging. Works like Sprintf.
func (l *SRunLogger) Warning(message string) {
	l.gLog.Warning(l.ctx, message)
}

// Error level logging. Works like Sprintf.
func (l *SRunLogger) Error(message string) {
	l.gLog.Error(l.ctx, message)
}

// Fatal level logging. Works like Sprintf.
func (l *SRunLogger) Fatal(message string) {
	l.gLog.Error(l.ctx, message)
	os.Exit(1)
}
