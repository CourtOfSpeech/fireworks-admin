package handle

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/domain/entity"
	"github.com/speech/fireworks-admin/internal/usecase"
	"github.com/speech/fireworks-admin/pkg/response"
)

// TeltentHandle 处理租户相关的 HTTP 请求。
type TeltentHandle struct {
	usecase *usecase.TeltentUsecase
}

// NewTeltentHandle 创建 TeltentHandle 实例。
func NewTeltentHandle(usecase *usecase.TeltentUsecase) *TeltentHandle {
	return &TeltentHandle{
		usecase: usecase,
	}
}

// Create 处理 POST /api/v1/tenants 创建租户请求。
func (h *TeltentHandle) Create(c *echo.Context) error {
	var req entity.CreateTeltentReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	teltent, err := h.usecase.Create(c.Request().Context(), &req)
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
func (h *TeltentHandle) FindByPage(c *echo.Context) error {
	var query entity.TeltentQuery
	if err := c.Bind(&query); err != nil {
		return response.BadRequest(c, "无效的查询参数")
	}

	result, err := h.usecase.FindByPage(c.Request().Context(), &query)
	if err != nil {
		return response.InternalError(c, "获取租户列表失败")
	}

	return response.Success(c, result)
}

// GetByID 处理 GET /api/v1/tenants/:id 查询单个租户请求。
func (h *TeltentHandle) GetByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "租户ID不能为空")
	}

	teltent, err := h.usecase.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.NotFound(c, "租户不存在")
	}

	return response.Success(c, teltent)
}

// Update 处理 PUT /api/v1/tenants/:id 更新租户请求。
func (h *TeltentHandle) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "租户ID不能为空")
	}

	var req entity.UpdateTeltentReq
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	teltent, err := h.usecase.Update(c.Request().Context(), id, &req)
	if err != nil {
		return response.InternalError(c, "更新租户失败")
	}

	return response.Success(c, teltent)
}

// Delete 处理 DELETE /api/v1/tenants/:id 删除租户请求。
func (h *TeltentHandle) Delete(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "租户ID不能为空")
	}

	if err := h.usecase.Delete(c.Request().Context(), id); err != nil {
		return response.InternalError(c, "删除租户失败")
	}

	return c.NoContent(http.StatusNoContent)
}
