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

// maxHeaderBytes 最大请求头字节数（1MB）。
const maxHeaderBytes = 1 << 20

// Server HTTP 服务器。
// 封装了 Echo 框架实例和标准库 http.Server，
// 并通过生命周期管理器实现优雅启动和关闭。
type Server struct {
	echo    *echo.Echo   // Echo 框架实例
	httpSrv *http.Server // 标准 HTTP 服务器
}

// NewServer 创建 HTTP 服务器实例。
// 该函数是 Wire 依赖注入的提供者，负责创建 HTTP 服务器并注册生命周期钩子。
// 服务器在生命周期启动时开始监听，在停止时执行优雅关闭。
// e 是 Echo 实例作为 HTTP 处理器，cfg 是应用配置用于获取服务器端口和超时配置，
// lc 是生命周期管理器用于注册启动和停止钩子。
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
					logger.Error(context.Background(), "HTTP服务异常", slog.Any("error", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info(ctx, "正在停止 HTTP 服务 (Graceful Shutdown)...")
			return httpSrv.Shutdown(ctx)
		},
	})
	return &Server{
		echo:    e,
		httpSrv: httpSrv,
	}
}

// newHTTPServer 创建 HTTP 服务器实例。
// 该函数根据配置创建标准库的 http.Server，设置监听地址、处理器和超时参数。
// e 是 Echo 实例作为 HTTP 处理器，cfg 是应用配置用于获取服务器端口和超时配置。
func newHTTPServer(e *echo.Echo, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),    // 监听地址
		Handler:        e,                                      // HTTP 处理器
		ReadTimeout:    durationOrDefault(cfg.Server.ReadTimeout, 15*time.Second),    // 读取超时
		WriteTimeout:   durationOrDefault(cfg.Server.WriteTimeout, 15*time.Second),   // 写入超时
		IdleTimeout:    durationOrDefault(cfg.Server.IdleTimeout, 60*time.Second),    // 空闲超时
		MaxHeaderBytes: bytesOrDefault(cfg.Server.MaxHeaderBytes, maxHeaderBytes),
	}
}

// durationOrDefault 返回配置的持续时间或默认值。
// 如果配置的秒数小于等于 0，则返回默认值。
// seconds 是配置的秒数，defaultValue 是默认持续时间。
func durationOrDefault(seconds int, defaultValue time.Duration) time.Duration {
	if seconds <= 0 {
		return defaultValue
	}
	return time.Duration(seconds) * time.Second
}

// bytesOrDefault 返回配置的字节数或默认值。
// 如果配置的字节数小于等于 0，则返回默认值。
// bytes 是配置的字节数，defaultValue 是默认字节数。
func bytesOrDefault(bytes int, defaultValue int) int {
	if bytes <= 0 {
		return defaultValue
	}
	return bytes
}
