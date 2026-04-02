package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

// RequestID 返回一个配置了项目默认设置的 RequestID 中间件。
// 该中间件为每个请求生成唯一的 UUID v4，便于日志追踪和问题排查。
func RequestID() echo.MiddlewareFunc {
	return echoMiddleware.RequestIDWithConfig(defaultRequestIDConfig)
}

// defaultRequestIDConfig 是项目级别的 RequestID 中间件默认配置。
// 跳过 OPTIONS 预检请求，使用 UUID v4 生成请求ID，并将ID存入 context。
var defaultRequestIDConfig = echoMiddleware.RequestIDConfig{
	Skipper: func(c *echo.Context) bool {
		return c.Request().Method == http.MethodOptions
	},
	Generator: func() string {
		return uuid.NewString()
	},
	RequestIDHandler: func(c *echo.Context, requestID string) {
		c.Set("request_id", requestID)
	},
	TargetHeader: echo.HeaderXRequestID,
}
