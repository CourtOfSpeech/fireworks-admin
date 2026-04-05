package app

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/middleware"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	"github.com/speech/fireworks-admin/internal/pkg/config"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
	"github.com/speech/fireworks-admin/internal/pkg/validator"
)

// NewEcho 创建 Echo 实例并注册中间件和路由。
func NewEcho(
	l *slog.Logger,
	cfg *config.Config,
	registrars []RouterRegistrar,
) *echo.Echo {
	e := newEcho(l)
	setupMiddleware(e, cfg)
	RegisterRoutes(e, registrars)
	return e
}

// newEcho 创建并配置 Echo 实例。
func newEcho(l *slog.Logger) *echo.Echo {
	return echo.NewWithConfig(echo.Config{
		Logger:             l,
		HTTPErrorHandler:   customHTTPErrorHandler,
		Validator:          validator.NewValidator(),
		FormParseMaxMemory: 10 << 20,
	})
}

// setupMiddleware 配置全局中间件。
func setupMiddleware(e *echo.Echo, cfg *config.Config) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS(cfg.Server.AllowOrigins))
	e.Use(middleware.Timeout(cfg.Server.Timeout))
}

// customHTTPErrorHandler 自定义 HTTP 错误处理器。
func customHTTPErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return
		}
	}
	code := http.StatusInternalServerError
	message := ""
	if he, ok := errors.AsType[*echo.HTTPError](err); ok {
		code = he.Code
		message = he.Message
	}

	// 根据状态码分级记录
	logMsg := "HTTP response error"
	attrs := []slog.Attr{
		slog.Int("status", code),
		slog.String("method", c.Request().Method),
		slog.String("path", c.Request().URL.Path),
		slog.Any("error", err),
	}
	if code >= 500 {
		logger.Error(logMsg, attrs...)
	} else {
		logger.Warn(logMsg, attrs...)
	}

	if c.Request().Method == http.MethodHead {
		_ = c.NoContent(code)
		return
	}
	_ = api.Error(c, code, message)
}
