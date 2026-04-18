package logger

import (
	"log/slog"

	"github.com/google/wire"
	"github.com/speech/fireworks-admin/internal/pkg/config"
)

// ProvideLogger 根据配置创建并返回一个日志记录器实例。
// 该函数是 Wire 依赖注入的提供者函数，从配置中读取日志级别、
// 输出格式和是否添加源码位置等参数，调用 NewLogger 创建日志记录器。
func ProvideLogger(cfg *config.Config) *slog.Logger {
	return NewLogger(cfg.Log.Level, cfg.Log.Format, cfg.Log.AddSource)
}

// ProviderSet 是日志模块的 Wire 依赖提供者集合。
// 该集合包含了日志模块所需的所有依赖提供者，用于在 Wire 依赖注入中使用。
var ProviderSet = wire.NewSet(ProvideLogger)
