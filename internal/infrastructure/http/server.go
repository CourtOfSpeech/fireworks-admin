package http

import (
	"errors"

	"net/http"

	"github.com/labstack/echo/v5"
	_ "github.com/lib/pq"
	"github.com/speech/fireworks-admin/internal/infrastructure/config"
	"github.com/speech/fireworks-admin/pkg/logger"
	"github.com/speech/fireworks-admin/pkg/response"
	"github.com/speech/fireworks-admin/pkg/validate"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
}

func NewServer(cfg *config.Config) (*Server, error) {
	logger := logger.NewLogger(cfg.Log.Level, cfg.Log.Format)
	e := echo.NewWithConfig(echo.Config{
		Logger:           logger,
		HTTPErrorHandler: customHTTPErrorHandler,
		Validator:        validate.NewValidator(),
	})

	return &Server{
		echo:   e,
		config: cfg,
	}, nil
}

// customHTTPErrorHandler handles HTTP errors.
// 执行流畅，先检查响应是否已提交，若已提交则直接返回，否则提取 HTTPError 信息，记录错误日志，返回 JSON 错误响应。
func customHTTPErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return // response has been already sent to the client by handler or some middleware
		}
	}
	code := http.StatusInternalServerError
	message := ""
	if he, ok := errors.AsType[*echo.HTTPError](err); ok {
		code = he.Code
		message = he.Message
	}

	c.Logger().Error("HTTP error",
		"status", code,
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"error", message,
	)

	if c.Request().Method == http.MethodHead {
		_ = c.NoContent(code)
		return
	}
	_ = response.Error(c, code, message)
}
