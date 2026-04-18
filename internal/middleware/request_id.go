// Package middleware 提供了 Echo 框架的 HTTP 中间件集合。
// 包含 CORS、Gzip 压缩、JWT 认证、日志记录、异常恢复、请求 ID 和超时控制等中间件。
// 所有中间件都提供了项目级别的默认配置，同时支持自定义配置。
package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/speech/fireworks-admin/internal/pkg/ctxutil"
)

// RequestID 返回一个配置了项目默认设置的请求 ID 中间件。
// 该中间件为每个请求生成唯一的 UUID 作为请求 ID，便于日志追踪和问题排查。
// 跳过 OPTIONS 预检请求，将请求 ID 存储到上下文中。
func RequestID() echo.MiddlewareFunc {
	return echoMiddleware.RequestIDWithConfig(echoMiddleware.RequestIDConfig{
		Skipper: func(c *echo.Context) bool {
			return c.Request().Method == http.MethodOptions
		},
		Generator: func() string {
			return uuid.NewString()
		},
		RequestIDHandler: func(c *echo.Context, requestID string) {
			c.Set("request_id", requestID)
			ctx := ctxutil.SetRequestID(c.Request().Context(), requestID)
			c.SetRequest(c.Request().WithContext(ctx))
		},
		TargetHeader: echo.HeaderXRequestID,
	})
}
