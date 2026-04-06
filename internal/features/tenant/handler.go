package tenant

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// Handler 处理租户相关的 HTTP 请求。
type Handler struct {
	service *Service
}

// NewHandler 创建 Handler 实例。
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes 实现 RouterRegistrar 接口。
// 租户相关路由注册到受保护组，需要 JWT 认证才能访问。
func (h *Handler) RegisterRoutes(public *echo.Group, protected *echo.Group) {
	tenants := protected.Group("/tenants")
	tenants.POST("", h.Create)
	tenants.GET("", h.FindByPage)
	tenants.GET("/:id", h.GetByID)
	tenants.PUT("/:id", h.Update)
	tenants.DELETE("/:id", h.Delete)
}

// Create 处理 POST /api/v1/tenants 创建租户请求。
func (h *Handler) Create(c *echo.Context) error {
	var req CreateTenantReq
	if err := c.Bind(&req); err != nil {
		return bizerr.InvalidParam(err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return bizerr.InvalidParam(err.Error())
	}

	tenant, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, api.ApiResponse{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    tenant,
	})
}

// FindByPage 处理 GET /api/v1/tenants 分页查询租户列表请求。
func (h *Handler) FindByPage(c *echo.Context) error {
	var query TenantQuery
	if err := c.Bind(&query); err != nil {
		return bizerr.InvalidParam("无效的查询参数")
	}

	result, err := h.service.FindByPage(c.Request().Context(), &query)
	if err != nil {
		return err
	}

	return api.Success(c, result)
}

// GetByID 处理 GET /api/v1/tenants/:id 查询单个租户请求。
func (h *Handler) GetByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return bizerr.InvalidParam("租户ID不能为空")
	}

	tenant, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return api.Success(c, tenant)
}

// Update 处理 PUT /api/v1/tenants/:id 更新租户请求。
func (h *Handler) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return bizerr.InvalidParam("租户ID不能为空")
	}

	var req UpdateTenantReq
	if err := c.Bind(&req); err != nil {
		return bizerr.InvalidParam("无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return bizerr.InvalidParam(err.Error())
	}

	tenant, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return err
	}

	return api.Success(c, tenant)
}

// Delete 处理 DELETE /api/v1/tenants/:id 删除租户请求。
func (h *Handler) Delete(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return bizerr.InvalidParam("租户ID不能为空")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
