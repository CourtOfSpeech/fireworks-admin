package logger

import (
	"log/slog"

	"github.com/google/wire"
	"github.com/speech/fireworks-admin/internal/pkg/config"
)

// ProvideLogger 提供日志记录器实例。
func ProvideLogger(cfg *config.Config) *slog.Logger {
	return NewLogger(cfg.Log.Level, cfg.Log.Format, cfg.Log.AddSource)
}

// ProviderSet 日志依赖提供者集合。
var ProviderSet = wire.NewSet(ProvideLogger)
