package middleware

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/speech/fireworks-admin/pkg/logger"
)

// Logger 返回一个配置了项目默认设置的日志中间件。
// 该中间件记录 HTTP 请求的详细信息，便于问题排查和性能监控。
func Logger() echo.MiddlewareFunc {
	m, err := defaultLoggerConfig.ToMiddleware()
	if err != nil {
		logger.Error("LOGGER_MIDDLEWARE_INIT_FAILED", slog.String("error", err.Error()))
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
	LogRequestID: true,
	LogRemoteIP:  true,
	HandleError:  true,
	LogValuesFunc: func(c *echo.Context, v echoMiddleware.RequestLoggerValues) error {
		if v.Error == nil {
			logger.Info("REQUEST",
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("latency", v.Latency.String()),
				slog.String("request_id", v.RequestID),
				slog.String("remote_ip", v.RemoteIP),
			)
		} else {
			logger.Error("REQUEST_ERROR",
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("latency", v.Latency.String()),
				slog.String("request_id", v.RequestID),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("error", v.Error.Error()),
			)
		}
		return nil
	},
}
