package router

import (
	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/di"
)

// RegisterRoutes 注册所有 HTTP 路由。
func RegisterRoutes(e *echo.Echo, app *di.App) {
	v1 := e.Group("/api/v1")

	tenants := v1.Group("/tenants")
	tenants.POST("", app.TeltentHandle.Create)
	tenants.GET("", app.TeltentHandle.FindByPage)
	tenants.GET("/:id", app.TeltentHandle.GetByID)
	tenants.PUT("/:id", app.TeltentHandle.Update)
	tenants.DELETE("/:id", app.TeltentHandle.Delete)
}
