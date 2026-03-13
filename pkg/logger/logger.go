package logger

import (
	"log/slog"
	"os"
	"time"
)

// NewLogger 创建并配置一个新的 slog.Logger 实例。
// 该函数接收日志级别、输出格式和是否添加源码位置三个参数，
// 内部通过解析日志级别、构建处理器配置、创建处理器等步骤完成日志记录器的初始化，
// 并将创建的 logger 设置为默认日志记录器。
func NewLogger(level string, format string, addSource bool) *slog.Logger {
	logLevel := parseLogLevel(level)
	opts := buildHandlerOptions(logLevel, addSource)

	handler := buildHandle(format, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// buildHandle 根据指定的格式创建日志处理器。
// format 参数支持 "json" 和 "text" 两种格式，分别创建 JSONHandler 和 TextHandler，
// 其他值默认使用 JSONHandler。opts 参数用于配置处理器的行为。
func buildHandle(format string, opts *slog.HandlerOptions) slog.Handler {
	switch format {
	case "json":
		return slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		return slog.NewTextHandler(os.Stdout, opts)
	default:
		return slog.NewJSONHandler(os.Stdout, opts)
	}
}

// buildHandlerOptions 构建日志处理器的配置选项。
// 该函数设置日志级别、是否添加源码位置信息，以及属性转换函数。
// 属性转换函数会将 password 字段脱敏为 "***REDACTED***"，并将时间字段重命名为 timestamp，
// 同时格式化为 "2006-01-02 15:04:05" 格式。
func buildHandlerOptions(level slog.Level, addSource bool) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "password" {
				return slog.String("password", "***REDACTED***")
			}
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.String("timestamp", a.Value.Time().Format(time.DateTime))
			}
			return a
		},
	}
}

// parseLogLevel 将字符串日志级别转换为 slog.Level 类型。
// 支持的级别包括："debug"、"info"、"warn"、"error"，
// 分别对应 LevelDebug、LevelInfo、LevelWarn、LevelError。
// 如果传入无法识别的字符串，默认返回 LevelInfo。
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
