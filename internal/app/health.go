package app

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// HealthRouter 健康检查路由注册器。
// 实现 RouterRegistrar 接口，提供 /health 和 /ready 端点。
type HealthRouter struct {
	client *ent.Client
}

// NewHealthRouter 创建健康检查路由注册器实例。
func NewHealthRouter(client *ent.Client) *HealthRouter {
	return &HealthRouter{client: client}
}

// RegisterRoutes 实现 RouterRegistrar 接口。
// 健康检查端点注册到公开组，无需认证即可访问。
func (h *HealthRouter) RegisterRoutes(public *echo.Group, protected *echo.Group) {
	public.GET("/health", h.livenessHandler)
	public.GET("/ready", h.readinessHandler)
}

// livenessHandler 处理 GET /health 请求，用于 Kubernetes Liveness 探针。
// 仅返回应用进程存活状态，不依赖外部服务，始终返回 {status: "ok"}。
func (h *HealthRouter) livenessHandler(c *echo.Context) error {
	return api.Success(c, map[string]string{
		"status": "ok",
	})
}

// readinessHandler 处理 GET /ready 请求，用于 Kubernetes Readiness 探针。
// 通过执行轻量数据库查询来验证应用是否已准备好接收流量。
// 数据库连接成功返回 200，失败返回 503 Service Unavailable。
func (h *HealthRouter) readinessHandler(c *echo.Context) error {
	ctx := c.Request().Context()

	if _, err := h.client.Tenant.Query().Count(ctx); err != nil {
		return bizerr.New(http.StatusServiceUnavailable, "数据库连接不可用: "+err.Error(), http.StatusServiceUnavailable)
	}

	return api.Success(c, map[string]string{
		"status": "ok",
	})
}

// 确保实现 RouterRegistrar 接口。
var _ RouterRegistrar = (*HealthRouter)(nil)
