package teltent

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/pkg/response"
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

// RegisterRoutes 注册租户模块路由。
func (h *Handler) RegisterRoutes(g *echo.Group) {
	tenants := g.Group("/tenants")
	tenants.POST("", h.Create)
	tenants.GET("", h.FindByPage)
	tenants.GET("/:id", h.GetByID)
	tenants.PUT("/:id", h.Update)
	tenants.DELETE("/:id", h.Delete)
}

// Create 处理 POST /api/v1/tenants 创建租户请求。
func (h *Handler) Create(c *echo.Context) error {
	var req CreateTeltentReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	teltent, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return response.InternalError(c, "创建租户失败")
	}

	return c.JSON(http.StatusCreated, response.ApiResponse{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    teltent,
	})
}

// FindByPage 处理 GET /api/v1/tenants 分页查询租户列表请求。
func (h *Handler) FindByPage(c *echo.Context) error {
	var query TeltentQuery
	if err := c.Bind(&query); err != nil {
		return response.BadRequest(c, "无效的查询参数")
	}

	result, err := h.service.FindByPage(c.Request().Context(), &query)
	if err != nil {
		return response.InternalError(c, "获取租户列表失败")
	}

	return response.Success(c, result)
}

// GetByID 处理 GET /api/v1/tenants/:id 查询单个租户请求。
func (h *Handler) GetByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "租户ID不能为空")
	}

	teltent, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.NotFound(c, "租户不存在")
	}

	return response.Success(c, teltent)
}

// Update 处理 PUT /api/v1/tenants/:id 更新租户请求。
func (h *Handler) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "租户ID不能为空")
	}

	var req UpdateTeltentReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	teltent, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return response.InternalError(c, "更新租户失败")
	}

	return response.Success(c, teltent)
}

// Delete 处理 DELETE /api/v1/tenants/:id 删除租户请求。
func (h *Handler) Delete(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "租户ID不能为空")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return response.InternalError(c, "删除租户失败")
	}

	return c.NoContent(http.StatusNoContent)
}
