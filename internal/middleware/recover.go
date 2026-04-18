// Package middleware 提供了 Echo 框架的 HTTP 中间件集合。
// 包含 CORS、Gzip 压缩、JWT 认证、日志记录、异常恢复、请求 ID 和超时控制等中间件。
// 所有中间件都提供了项目级别的默认配置，同时支持自定义配置。
package middleware

import (
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

// Recover 返回一个配置了项目默认设置的 Recover 中间件。
// 该中间件捕获处理器中的 panic，防止应用崩溃并返回 500 错误。
func Recover() echo.MiddlewareFunc {
	return echoMiddleware.RecoverWithConfig(defaultRecoverConfig)
}

// defaultRecoverConfig 是项目级别的 Recover 中间件默认配置。
// 堆栈大小 4KB，打印所有 goroutine 堆栈。
var defaultRecoverConfig = echoMiddleware.RecoverConfig{
	Skipper:           nil,
	StackSize:         4 << 10,
	DisablePrintStack: false,
	DisableStackAll:   false,
}
