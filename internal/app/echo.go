package app

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/middleware"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	"github.com/speech/fireworks-admin/internal/pkg/config"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
	"github.com/speech/fireworks-admin/internal/pkg/validator"
)

// NewEcho 创建 Echo 实例并注册中间件和路由。
// 该函数是 Wire 依赖注入的提供者，负责初始化 Echo 框架并配置完整的请求处理链。
// l 是日志记录器用于记录请求和错误信息，cfg 是应用配置用于获取服务器相关配置，
// registrars 是路由注册器列表用于注册各业务模块的路由。
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
// 该函数初始化 Echo 框架的基础配置，包括日志、错误处理器、验证器和表单解析限制。
// l 是日志记录器。
func newEcho(l *slog.Logger) *echo.Echo {
	return echo.NewWithConfig(echo.Config{
		Logger:             l,                      // 日志记录器
		HTTPErrorHandler:   customHTTPErrorHandler, // 自定义错误处理器
		Validator:          validator.NewValidator(), // 请求验证器
		FormParseMaxMemory: 10 << 20,               // 表单解析最大内存 (10MB)
	})
}

// setupMiddleware 配置全局中间件。
// 该函数按顺序注册以下中间件：
//   - RequestID: 为每个请求生成唯一 ID
//   - Logger: 记录请求日志
//   - Recover: 恢复 panic 并返回 500 错误
//   - CORS: 跨域资源共享配置
//   - Timeout: 请求超时控制
// e 是 Echo 实例，cfg 是应用配置用于获取 CORS 和超时配置。
func setupMiddleware(e *echo.Echo, cfg *config.Config) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS(cfg.Server.AllowOrigins))
	e.Use(middleware.Timeout(cfg.Server.Timeout))
}

// customHTTPErrorHandler 自定义 HTTP 错误处理器。
// 该函数统一处理所有 HTTP 错误，包括业务错误和框架错误。
// 对于 5xx 错误记录 Error 级别日志，对于 4xx 错误记录 Warn 级别日志。
// 响应格式统一为 JSON 格式的 ApiResponse 结构。
// c 是 Echo 上下文，err 是错误对象。
func customHTTPErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return
		}
	}
	code := bizerr.ErrInternal
	message := "internal server error"
	httpStatus := http.StatusInternalServerError
	var stackAttrs slog.Value
	if biz, ok := errors.AsType[*bizerr.BizError](err); ok {
		httpStatus = biz.HTTPStatus
		message = biz.Message
		code = biz.Code
		if len(biz.Stack) > 0 {
			stackAttrs = biz.StackValue()
		}
	} else if he, ok := errors.AsType[*echo.HTTPError](err); ok {
		code = he.Code
		message = he.Message
		httpStatus = he.Code
	}

	logMsg := "HTTP response error"
	attrs := []slog.Attr{
		slog.Int("status", code),
		slog.String("method", c.Request().Method),
		slog.String("path", c.Request().URL.Path),
		slog.String("message", message),
		slog.Any("error", err),
	}
	if stackAttrs.Kind() == slog.KindString {
		attrs = append(attrs, slog.Any("stack", stackAttrs))
	}
	if code >= 500 {
		logger.Error(c.Request().Context(), logMsg, attrs...)
	} else {
		logger.Warn(c.Request().Context(), logMsg, attrs...)
	}

	if c.Request().Method == http.MethodHead {
		_ = c.NoContent(code)
		return
	}
	_ = c.JSON(httpStatus, api.ApiResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
