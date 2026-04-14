package tenant

import (
	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

type TenantHandler struct {
	service *TenantService
}

func NewTenantHandler(service *TenantService) *TenantHandler {
	return &TenantHandler{
		service: service,
	}
}

func (h *TenantHandler) RegisterRoutes(public *echo.Group, protected *echo.Group) {
	tenants := protected.Group("/tenants")
	tenants.POST("", h.Create)
	tenants.GET("", h.List)
	tenants.GET("/:id", h.Get)
	tenants.PUT("/:id", h.Update)
	tenants.DELETE("/:id", h.Delete)
}

func (h *TenantHandler) Create(c *echo.Context) error {
	var req CreateTenantReq
	if err := c.Bind(&req); err != nil {
		return bizerr.InvalidParam(err.Error())
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

func (h *TenantHandler) List(c *echo.Context) error {
	var query TenantQuery
	if err := c.Bind(&query); err != nil {
		return bizerr.InvalidParam("无效的查询参数")
	}

	result, err := h.service.List(c.Request().Context(), &query)
	if err != nil {
		return err
	}

	return api.Success(c, result)
}

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

func (h *TenantHandler) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return bizerr.InvalidParam("租户ID不能为空")
	}

	var req UpdateTenantReq
	if err := c.Bind(&req); err != nil {
		return bizerr.InvalidParam("无效的请求参数")
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
