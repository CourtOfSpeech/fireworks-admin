// Package middleware 提供了 Echo 框架的 HTTP 中间件集合。
// 包含 CORS、Gzip 压缩、JWT 认证、日志记录、异常恢复、请求 ID 和超时控制等中间件。
// 所有中间件都提供了项目级别的默认配置，同时支持自定义配置。
package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

// Timeout 返回使用指定超时时间的超时中间件。
// 该中间件为请求设置超时限制，超时后返回 503 错误。
// 跳过 OPTIONS 预检请求。
func Timeout(timeout int) echo.MiddlewareFunc {
	return echoMiddleware.ContextTimeoutWithConfig(echoMiddleware.ContextTimeoutConfig{
		Skipper: func(c *echo.Context) bool {
			return c.Request().Method == http.MethodOptions
		},
		Timeout: time.Duration(timeout) * time.Second,
	})
}
