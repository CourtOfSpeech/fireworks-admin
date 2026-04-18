// Package api 提供了 API 相关的通用类型和工具函数。
// 包含分页查询、响应格式化等常用功能，用于构建统一的 API 接口。
package api

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// ApiResponse 统一的 API 响应结构体。
// 用于标准化所有 API 接口的返回格式，包含状态码、消息和数据。
type ApiResponse struct {
	Code    int    `json:"code"`    // Code HTTP 状态码或业务状态码
	Message string `json:"message"` // Message 响应消息，描述请求处理结果
	Data    any    `json:"data"`    // Data 响应数据，可以是任意类型
}

// Success 发送成功的 API 响应。
// c 是 Echo 框架的上下文指针，data 是要返回的数据可以是任意类型。
// 返回 JSON 格式的成功响应，状态码为 200。
func Success(c *echo.Context, data any) error {
	return c.JSON(http.StatusOK, ApiResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}
