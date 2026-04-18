// Package middleware 提供了 Echo 框架的 HTTP 中间件集合。
// 包含 CORS、Gzip 压缩、JWT 认证、日志记录、异常恢复、请求 ID 和超时控制等中间件。
// 所有中间件都提供了项目级别的默认配置，同时支持自定义配置。
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
)

// Logger 返回一个配置了项目默认设置的日志中间件。
// 该中间件记录 HTTP 请求的详细信息，便于问题排查和性能监控。
func Logger() echo.MiddlewareFunc {
	m, err := defaultLoggerConfig.ToMiddleware()
	if err != nil {
		logger.Error(context.Background(), "LOGGER_MIDDLEWARE_INIT_FAILED", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return m
}

// defaultLoggerConfig 是项目级别的默认日志中间件配置。
// 跳过 OPTIONS 预检请求，记录请求方法、URI、状态码、耗时等关键信息。
var defaultLoggerConfig = echoMiddleware.RequestLoggerConfig{
	Skipper: func(c *echo.Context) bool {
		return c.Request().Method == http.MethodOptions
	},
	LogLatency:   true,
	LogMethod:    true,
	LogURI:       true,
	LogStatus:    true,
	LogRequestID: false,
	LogRemoteIP:  true,
	HandleError:  true,
	LogValuesFunc: func(c *echo.Context, v echoMiddleware.RequestLoggerValues) error {
		ctx := c.Request().Context()
		if v.Error == nil {
			logger.Info(ctx, "request processed",
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("latency", v.Latency.String()),
				slog.String("remote_ip", v.RemoteIP),
			)
		} else {
			logger.Error(ctx, "request error",
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("latency", v.Latency.String()),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("error", v.Error.Error()),
			)
		}
		return nil
	},
}
