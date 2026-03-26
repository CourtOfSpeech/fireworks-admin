package provider

import (
	"log/slog"

	"github.com/speech/fireworks-admin/internal/infrastructure/config"
	"github.com/speech/fireworks-admin/pkg/logger"
)

// ProvideConfig 提供配置实例
// 返回:
//   - *config.Config: 配置实例
//   - error: 加载错误
func ProvideConfig() (*config.Config, error) {
	return config.LoadByEnv()
}

// ProvideLogger 提供日志记录器
// 参数:
//   - cfg: 配置实例
// 返回:
//   - *slog.Logger: 日志记录器实例
func ProvideLogger(cfg *config.Config) *slog.Logger {
	return logger.NewLogger(cfg.Log.Level, cfg.Log.Format, cfg.Log.AddSource)
}

// ProvideDSN 提供数据库连接字符串
// 参数:
//   - cfg: 配置实例
// 返回:
//   - string: 数据库连接字符串
func ProvideDSN(cfg *config.Config) string {
	return cfg.Database.DSN()
}
