package app

import "github.com/labstack/echo/v5"

// RouterRegistrar 路由注册器接口。
// 各功能模块通过实现此接口将路由注册到公开组和受保护组。
// 公开组（public）用于无需认证的端点，如健康检查。
// 受保护组（protected）用于需要认证的端点，如业务 API。
type RouterRegistrar interface {
	// RegisterRoutes 注册路由到公开组和受保护组。
	// public: 公开路由组，无需认证即可访问
	// protected: 受保护路由组，需要 JWT 认证才能访问
	RegisterRoutes(public *echo.Group, protected *echo.Group)
}

// RegisterRoutes 注册所有 HTTP 路由。
// 公开组（public）用于无需认证的端点，如健康检查。
// 受保护组（protected）用于需要认证的端点，如业务 API。
func RegisterRoutes(e *echo.Echo, registrars []RouterRegistrar) {
	public := e.Group("")
	protected := e.Group("/api/v1")

	for _, r := range registrars {
		r.RegisterRoutes(public, protected)
	}
}
