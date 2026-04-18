// Package logger 提供结构化的日志记录功能。
// 该包基于 Go 标准库的 log/slog 实现，支持 JSON 和文本两种输出格式，
// 并提供自动从 Context 中提取请求 ID 等上下文信息的能力。
// 同时支持日志脱敏、源码位置追踪等特性。
package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/speech/fireworks-admin/internal/pkg/ctxutil"
)

// projectRoot 存储项目根目录路径，用于在日志输出中简化源码路径显示。
var projectRoot = detectProjectRoot()

// ContextHandler 负责在日志写入前，自动从 Context 提取信息。
// 它是一个 slog.Handler 的包装器，在处理日志记录时会自动附加
// Context 中的上下文信息（如请求 ID）到日志属性中。
type ContextHandler struct {
	slog.Handler
}

// Handle 处理日志记录，自动从 Context 中提取请求 ID 等信息并添加到日志属性中。
// 该方法实现了 slog.Handler 接口，在调用底层 Handler 之前，
// 会检查 Context 中是否存在请求 ID，如果存在则将其添加到日志记录中。
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if id := ctxutil.GetRequestID(ctx); id != "" {
			r.AddAttrs(slog.String("request_id", id))
		}
	}
	return h.Handler.Handle(ctx, r)
}

// NewLogger 创建并返回一个新的结构化日志记录器。
// level 是日志级别支持 "debug"、"info"、"warn"、"error"，
// format 是输出格式支持 "json" 和 "text"（不区分大小写），
// addSource 表示是否在日志中添加源码位置信息。
// 该函数会自动对 "password" 字段进行脱敏处理，并将日志记录器设置为默认记录器。
func NewLogger(level string, format string, addSource bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     parseLogLevel(level),
		AddSource: addSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "password" {
				return slog.String("password", "***REDACTED***")
			}
			if addSource && a.Key == slog.SourceKey {
				if src, ok := a.Value.Any().(*slog.Source); ok {
					file := src.File
					if projectRoot != "" && strings.HasPrefix(file, projectRoot) {
						file = file[len(projectRoot):]
					}
					return slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", file, src.Line))
				}
			}
			return a
		},
	}

	var baseHandler slog.Handler
	if strings.ToLower(format) == "json" {
		baseHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		baseHandler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(&ContextHandler{Handler: baseHandler})
	slog.SetDefault(logger)
	return logger
}

// logHelper 是一个内部辅助函数，用于统一处理日志记录的创建和写入。
// 它通过 runtime.Callers 获取调用者的程序计数器，以确保日志记录中
// 显示正确的源码位置。该函数会检查日志级别是否启用，只有启用时才会创建日志记录。
func logHelper(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	l := slog.Default()
	if !l.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = l.Handler().Handle(ctx, r)
}

// Info 记录一条 INFO 级别的日志消息。
// 该函数接受 Context、消息字符串和可选的日志属性，并将日志输出到标准输出。
func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	logHelper(ctx, slog.LevelInfo, msg, attrs...)
}

// Debug 记录一条 DEBUG 级别的日志消息。
// 该函数接受 Context、消息字符串和可选的日志属性，并将日志输出到标准输出。
func Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	logHelper(ctx, slog.LevelDebug, msg, attrs...)
}

// Warn 记录一条 WARN 级别的日志消息。
// 该函数接受 Context、消息字符串和可选的日志属性，并将日志输出到标准输出。
func Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	logHelper(ctx, slog.LevelWarn, msg, attrs...)
}

// Error 记录一条 ERROR 级别的日志消息。
// 该函数接受 Context、消息字符串和可选的日志属性，并将日志输出到标准输出。
func Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	logHelper(ctx, slog.LevelError, msg, attrs...)
}

// parseLogLevel 将字符串形式的日志级别转换为 slog.Level 类型。
// 支持的级别（不区分大小写）：
//   - "debug" -> slog.LevelDebug
//   - "warn" -> slog.LevelWarn
//   - "error" -> slog.LevelError
//   - 其他值 -> slog.LevelInfo（默认）
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// detectProjectRoot 检测并返回项目的根目录路径。
// 该函数通过获取当前源文件的路径，查找 "/internal" 目录的位置，
// 从而推断出项目根目录。返回的路径包含末尾的路径分隔符。
// 如果无法检测到项目根目录，则返回空字符串。
func detectProjectRoot() string {
	_, file, _, _ := runtime.Caller(0)
	if idx := strings.Index(file, "/internal"); idx != -1 {
		return file[:idx+1]
	}
	return ""
}
