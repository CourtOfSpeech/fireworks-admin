package api

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
