// Package user 提供User功能，包括User的创建、查询、更新和删除操作。
// 本文件定义了User模块的 HTTP 处理器，负责处理User相关的 API 请求。
package user

import (
	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// UserHandler User HTTP 处理器。
// 负责处理User相关的 HTTP 请求，包括创建、查询、更新和删除操作。
type UserHandler struct {
	service *UserService // User业务逻辑服务
}

// NewUserHandler 创建User处理器实例。
// 参数 service 为User业务逻辑服务，返回初始化后的处理器实例。
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// RegisterRoutes 注册User模块的路由。
// 参数 public 为公开路由组，protected 为需认证的路由组。
// User相关接口均需要认证，注册在 /users 路径下。
func (h *UserHandler) RegisterRoutes(public *echo.Group, protected *echo.Group) {
	users := protected.Group("/users")
	users.POST("", h.Create)
	users.GET("", h.List)
	users.GET("/:id", h.Get)
	users.PUT("/:id", h.Update)
	users.DELETE("/:id", h.Delete)
}

// Create 处理创建User请求。
// 绑定并验证请求参数，调用服务层创建User。
// 成功返回空响应，失败返回相应错误。
func (h *UserHandler) Create(c *echo.Context) error {
	var req CreateUserReq
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

// List 处理查询User列表请求。
// 绑定查询参数，调用服务层获取User列表。
// 成功返回分页结果，失败返回相应错误。
func (h *UserHandler) List(c *echo.Context) error {
	var query UserQuery
	if err := c.Bind(&query); err != nil {
		return bizerr.InvalidParamWrap(err, "无效的查询参数")
	}

	result, err := h.service.List(c.Request().Context(), &query)
	if err != nil {
		return err
	}

	return api.Success(c, result)
}

// Get 处理查询单个User请求。
// 从 URL 路径获取User ID，调用服务层查询User详情。
// 成功返回User信息，失败返回相应错误。
func (h *UserHandler) Get(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return bizerr.InvalidParam("UserID不能为空")
	}

	user, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return api.Success(c, user)
}

// Update 处理更新User请求。
// 从 URL 路径获取User ID，绑定并验证请求参数，调用服务层更新User。
// 成功返回更新后的User信息，失败返回相应错误。
func (h *UserHandler) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return bizerr.InvalidParam("UserID不能为空")
	}

	var req UpdateUserReq
	if err := c.Bind(&req); err != nil {
		return bizerr.InvalidParamWrap(err, "无效的请求参数")
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	user, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return err
	}

	return api.Success(c, user)
}

// Delete 处理删除User请求。
// 从 URL 路径获取User ID，调用服务层删除User。
// 成功返回空响应，失败返回相应错误。
func (h *UserHandler) Delete(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return bizerr.InvalidParam("UserID不能为空")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return api.Success(c, nil)
}
