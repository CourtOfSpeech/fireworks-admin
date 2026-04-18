// Package tenant 提供租户管理功能，包括租户的创建、查询、更新和删除操作。
// 本文件定义了租户模块的 HTTP 处理器，负责处理租户相关的 API 请求。
package tenant

import (
	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// TenantHandler 租户 HTTP 处理器。
// 负责处理租户相关的 HTTP 请求，包括创建、查询、更新和删除操作。
type TenantHandler struct {
	service *TenantService // 租户业务逻辑服务
}

// NewTenantHandler 创建租户处理器实例。
// 参数 service 为租户业务逻辑服务，返回初始化后的处理器实例。
func NewTenantHandler(service *TenantService) *TenantHandler {
	return &TenantHandler{
		service: service,
	}
}

// RegisterRoutes 注册租户模块的路由。
// 参数 public 为公开路由组，protected 为需认证的路由组。
// 租户相关接口均需要认证，注册在 /tenants 路径下。
func (h *TenantHandler) RegisterRoutes(public *echo.Group, protected *echo.Group) {
	tenants := protected.Group("/tenants")
	tenants.POST("", h.Create)
	tenants.GET("", h.List)
	tenants.GET("/:id", h.Get)
	tenants.PUT("/:id", h.Update)
	tenants.DELETE("/:id", h.Delete)
}

// Create 处理创建租户请求。
// 绑定并验证请求参数，调用服务层创建租户。
// 成功返回空响应，失败返回相应错误。
func (h *TenantHandler) Create(c *echo.Context) error {
	var req CreateTenantReq
	if err := c.Bind(&req); err != nil {
		return bizerr.InvalidParamWrap(err, "无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	_, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return api.Success(c, nil)

}

// List 处理查询租户列表请求。
// 绑定查询参数，调用服务层获取租户列表。
// 成功返回分页结果，失败返回相应错误。
func (h *TenantHandler) List(c *echo.Context) error {
	var query TenantQuery
	if err := c.Bind(&query); err != nil {
		return bizerr.InvalidParamWrap(err, "无效的查询参数")
	}

	result, err := h.service.List(c.Request().Context(), &query)
	if err != nil {
		return err
	}

	return api.Success(c, result)
}

// Get 处理查询单个租户请求。
// 从 URL 路径获取租户 ID，调用服务层查询租户详情。
// 成功返回租户信息，失败返回相应错误。
func (h *TenantHandler) Get(c *echo.Context) error {
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

// Update 处理更新租户请求。
// 从 URL 路径获取租户 ID，绑定并验证请求参数，调用服务层更新租户。
// 成功返回更新后的租户信息，失败返回相应错误。
func (h *TenantHandler) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return bizerr.InvalidParam("租户ID不能为空")
	}

	var req UpdateTenantReq
	if err := c.Bind(&req); err != nil {
		return bizerr.InvalidParamWrap(err, "无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	tenant, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return err
	}

	return api.Success(c, tenant)
}

// Delete 处理删除租户请求。
// 从 URL 路径获取租户 ID，调用服务层删除租户。
// 成功返回空响应，失败返回相应错误。
func (h *TenantHandler) Delete(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return bizerr.InvalidParam("租户ID不能为空")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return api.Success(c, nil)
}
