package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/pkg/config"
	"github.com/speech/fireworks-admin/internal/pkg/lifecycle"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
)

// Server HTTP 服务器。
type Server struct {
	echo    *echo.Echo
	httpSrv *http.Server
}

// NewServer 创建 HTTP 服务器实例。
func NewServer(
	e *echo.Echo,
	cfg *config.Config,
	lc *lifecycle.Lifecycle,
) *Server {
	httpSrv := newHTTPServer(e, cfg)

	lc.Append(lifecycle.Hook{
		Name: "HTTP-Server",
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP服务异常", slog.Any("error", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("正在停止 HTTP 服务 (Graceful Shutdown)...")
			return httpSrv.Shutdown(ctx)
		},
	})
	return &Server{
		echo:    e,
		httpSrv: httpSrv,
	}
}

// newHTTPServer 创建 HTTP 服务器实例。
func newHTTPServer(e *echo.Echo, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        e,
		ReadTimeout:    durationOrDefault(cfg.Server.ReadTimeout, 15*time.Second),
		WriteTimeout:   durationOrDefault(cfg.Server.WriteTimeout, 15*time.Second),
		IdleTimeout:    durationOrDefault(cfg.Server.IdleTimeout, 60*time.Second),
		MaxHeaderBytes: bytesOrDefault(cfg.Server.MaxHeaderBytes, 1<<20),
	}
}

// durationOrDefault 返回配置的持续时间或默认值。
func durationOrDefault(seconds int, defaultValue time.Duration) time.Duration {
	if seconds <= 0 {
		return defaultValue
	}
	return time.Duration(seconds) * time.Second
}

// bytesOrDefault 返回配置的字节数或默认值。
func bytesOrDefault(bytes int, defaultValue int) int {
	if bytes <= 0 {
		return defaultValue
	}
	return bytes
}
