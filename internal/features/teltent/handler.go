package teltent

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
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
	var req CreateTeltentReq
	if err := c.Bind(&req); err != nil {
		return api.BadRequest(c, "无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return api.BadRequest(c, err.Error())
	}

	teltent, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return h.handleError(c, "Create.Teltent", err)
	}

	return c.JSON(http.StatusCreated, api.ApiResponse{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    teltent,
	})
}

// FindByPage 处理 GET /api/v1/tenants 分页查询租户列表请求。
func (h *Handler) FindByPage(c *echo.Context) error {
	var query TeltentQuery
	if err := c.Bind(&query); err != nil {
		return api.BadRequest(c, "无效的查询参数")
	}

	result, err := h.service.FindByPage(c.Request().Context(), &query)
	if err != nil {
		return h.handleError(c, "FindByPage.Teltent", err)
	}

	return api.Success(c, result)
}

// GetByID 处理 GET /api/v1/tenants/:id 查询单个租户请求。
func (h *Handler) GetByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return api.BadRequest(c, "租户ID不能为空")
	}

	teltent, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return h.handleError(c, "GetByID.Teltent", err)
	}

	return api.Success(c, teltent)
}

// Update 处理 PUT /api/v1/tenants/:id 更新租户请求。
func (h *Handler) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return api.BadRequest(c, "租户ID不能为空")
	}

	var req UpdateTeltentReq
	if err := c.Bind(&req); err != nil {
		return api.BadRequest(c, "无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return api.BadRequest(c, err.Error())
	}

	teltent, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return h.handleError(c, "Update.Teltent", err)
	}

	return api.Success(c, teltent)
}

// Delete 处理 DELETE /api/v1/tenants/:id 删除租户请求。
func (h *Handler) Delete(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return api.BadRequest(c, "租户ID不能为空")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return h.handleError(c, "Delete.Teltent", err)
	}

	return c.NoContent(http.StatusNoContent)
}

// handleError 统一处理业务层返回的错误，根据错误类型映射到对应的 HTTP 响应。
// 错误处理策略：
//   - NotFoundError    → 404 Not Found
//   - ConflictError    → 409 Conflict
//   - InvalidArgumentError → 400 Bad Request
//   - 其他 BizError     → 使用 BizError 中定义的 HTTPStatus
//   - 未知错误          → 500 Internal Server Error
//
// 同时记录包含请求上下文的详细错误日志，便于问题排查。
func (h *Handler) handleError(c *echo.Context, operation string, err error) error {
	requestID, _ := c.Get("request_id").(string)

	switch {
	case bizerr.IsNotFoundError(err):
		logger.WarnCtx(c.Request().Context(), "BUSINESS_WARNING",
			slog.String("operation", operation),
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return api.NotFound(c, err.Error())

	case bizerr.IsConflictError(err):
		logger.WarnCtx(c.Request().Context(), "BUSINESS_WARNING",
			slog.String("operation", operation),
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return api.Conflict(c, err.Error())

	case bizerr.IsBizError(err):
		logger.WarnCtx(c.Request().Context(), "BUSINESS_ERROR",
			slog.String("operation", operation),
			slog.String("request_id", requestID),
			slog.Int("http_status", bizerr.HTTPStatus(err)),
			slog.String("error", err.Error()),
		)
		status := bizerr.HTTPStatus(err)
		return api.Error(c, status, err.Error())

	default:
		logger.ErrorCtx(c.Request().Context(), "INTERNAL_ERROR",
			slog.String("operation", operation),
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return api.InternalError(c, "服务器内部错误")
	}
}
