package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c *echo.Context, data any) error {
	return c.JSON(http.StatusOK, ApiResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func Error(c *echo.Context, code int, message string) error {
	return c.JSON(code, ApiResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func BadRequest(c *echo.Context, message string) error {
	return Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message)
}

func Forbidden(c *echo.Context, message string) error {
	return Error(c, http.StatusForbidden, message)
}

func NotFound(c *echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message)
}

func InternalError(c *echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, message)
}
