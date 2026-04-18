// Package middleware 提供了 Echo 框架的 HTTP 中间件集合。
// 包含 CORS、Gzip 压缩、JWT 认证、日志记录、异常恢复、请求 ID 和超时控制等中间件。
// 所有中间件都提供了项目级别的默认配置，同时支持自定义配置。
package middleware

import (
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

// defaultGzipLevel 是项目默认的 Gzip 压缩级别
// 压缩级别范围: 1-9，其中 1 压缩速度最快但压缩率最低，9 压缩率最高但速度最慢
// 6 是性能和压缩率的平衡点
const defaultGzipLevel = 6

// Gzip 返回一个配置了项目默认参数的 Gzip 中间件。
// 该中间件用于压缩 HTTP 响应，减少传输数据量，提高性能。
//
// 默认配置:
//   - Level: 6 (平衡性能和压缩率)
//
// 使用示例:
//
//	e.Use(middleware.Gzip())
func Gzip() echo.MiddlewareFunc {
	return GzipWithConfig(GzipConfig{})
}

// GzipConfig 定义 Gzip 中间件的配置选项
type GzipConfig = echoMiddleware.GzipConfig

// GzipWithConfig 返回一个使用自定义配置的 Gzip 中间件。
// 该函数允许覆盖默认配置，提供更灵活的压缩控制。
// config 是 Gzip 配置选项，如果 Level 为 0 则使用默认值 6。
//
// 使用示例:
//
//	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
//	    Level: 6,
//	    Skipper: func(c *echo.Context) bool {
//	        // 跳过已压缩的文件或小文件
//	        return c.Request().URL.Path == "/static/image.png"
//	    },
//	}))
func GzipWithConfig(config GzipConfig) echo.MiddlewareFunc {
	if config.Level == 0 {
		config.Level = defaultGzipLevel
	}

	return echoMiddleware.GzipWithConfig(config)
}
