package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

// Timeout 返回使用指定超时时间的超时中间件。
// 超时后返回 503 错误。
func Timeout(timeout int) echo.MiddlewareFunc {
	return echoMiddleware.ContextTimeoutWithConfig(echoMiddleware.ContextTimeoutConfig{
		Skipper: func(c *echo.Context) bool {
			return c.Request().Method == http.MethodOptions
		},
		Timeout: time.Duration(timeout) * time.Second,
	})
}
