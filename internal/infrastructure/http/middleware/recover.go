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
