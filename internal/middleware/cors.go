// Package middleware 提供了 Echo 框架的 HTTP 中间件集合。
// 包含 CORS、Gzip 压缩、JWT 认证、日志记录、异常恢复、请求 ID 和超时控制等中间件。
// 所有中间件都提供了项目级别的默认配置，同时支持自定义配置。
package middleware

import (
	"net/http"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

// CORS 返回使用默认配置的 CORS 中间件。
// 该函数接收允许的来源地址列表，如果列表为空则允许所有来源（"*"）。
// 适用于大多数前后端分离场景。
func CORS(allowOrigins []string) echo.MiddlewareFunc {
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}
	return echoMiddleware.CORSWithConfig(defaultCORSConfig(allowOrigins))
}

// defaultCORSConfig 返回项目级别的默认 CORS 配置。
// 该配置允许常见的 HTTP 方法和请求头，并启用凭证支持。
func defaultCORSConfig(allowOrigins []string) echoMiddleware.CORSConfig {
	return echoMiddleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXRequestID,
		},
		AllowCredentials: true,
		ExposeHeaders: []string{
			echo.HeaderXRequestID,
		},
		MaxAge: 86400,
	}
}
