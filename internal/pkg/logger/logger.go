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

var projectRoot = detectProjectRoot()

// ContextHandler 负责在日志写入前，自动从 Context 提取信息
type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if id := ctxutil.GetRequestID(ctx); id != "" {
			r.AddAttrs(slog.String("request_id", id))
		}
	}
	return h.Handler.Handle(ctx, r)
}

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

func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	logHelper(ctx, slog.LevelInfo, msg, attrs...)
}

func Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	logHelper(ctx, slog.LevelDebug, msg, attrs...)
}

func Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	logHelper(ctx, slog.LevelWarn, msg, attrs...)
}

func Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	logHelper(ctx, slog.LevelError, msg, attrs...)
}

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

func detectProjectRoot() string {
	_, file, _, _ := runtime.Caller(0)
	if idx := strings.Index(file, "/internal"); idx != -1 {
		return file[:idx+1]
	}
	return ""
}
