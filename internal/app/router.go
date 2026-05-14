package app

import (
	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/middleware"
	"github.com/speech/fireworks-admin/internal/pkg/config"
)

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
// 该函数创建公开组和受保护组，并遍历所有路由注册器进行路由注册。
// 公开组（public）用于无需认证的端点，如健康检查，路径前缀为空。
// 受保护组（protected）用于需要认证的端点，如业务 API，路径前缀为 /api/v1。
// e 是 Echo 实例，registrars 是路由注册器列表。
func RegisterRoutes(e *echo.Echo, registrars []RouterRegistrar, cfg *config.Config) {
	public := e.Group("")
	protected := e.Group("/api/v1")
	protected.Use(middleware.NewJWTMiddleware(&middleware.JWTConfig{
		Secret:     cfg.JWT.Secret,
		ExpireTime: cfg.JWT.ExpireTime,
	}))

	for _, r := range registrars {
		r.RegisterRoutes(public, protected)
	}
}
