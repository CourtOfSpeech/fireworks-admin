package app

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes 注册所有 HTTP 路由。
// 采用分布式定义、集中注册的策略，各模块路由在各自的 Handler 中定义。
func RegisterRoutes(e *echo.Echo, a *App) {
	v1 := e.Group("/api/v1")

	// 注册各模块路由
	a.TeltentHandler.RegisterRoutes(v1)
}
