package app

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/middleware"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
	"github.com/speech/fireworks-admin/internal/pkg/validator"
)

// Server HTTP 服务器。
type Server struct {
	echo    *echo.Echo
	app     *App
	cleanup func()
}

// NewServer 创建 HTTP 服务器实例。
func NewServer(a *App, cleanup func()) *Server {
	e := echo.NewWithConfig(echo.Config{
		Logger:             a.Logger,
		HTTPErrorHandler:   customHTTPErrorHandler,
		Validator:          validator.NewValidator(),
		FormParseMaxMemory: 10 << 20,
	})

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS(a.Config.Server.AllowOrigins))
	e.Use(middleware.Timeout(a.Config.Server.Timeout))

	return &Server{
		echo:    e,
		app:     a,
		cleanup: cleanup,
	}
}

// Start 启动 HTTP 服务器。
func (s *Server) Start() error {
	logger.Info("server starting", slog.Int("port", s.app.Config.Server.Port))
	if err := s.echo.Start(fmt.Sprintf(":%d", s.app.Config.Server.Port)); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

// Echo 返回 Echo 实例。
func (s *Server) Echo() *echo.Echo {
	return s.echo
}

// customHTTPErrorHandler 自定义 HTTP 错误处理器。
func customHTTPErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return
		}
	}
	code := http.StatusInternalServerError
	message := ""
	if he, ok := errors.AsType[*echo.HTTPError](err); ok {
		code = he.Code
		message = he.Message
	}

	logger.Error("HTTP error",
		slog.Int("status", code),
		slog.String("method", c.Request().Method),
		slog.String("path", c.Request().URL.Path),
		slog.String("error", message),
	)

	if c.Request().Method == http.MethodHead {
		_ = c.NoContent(code)
		return
	}
	_ = api.Error(c, code, message)
}
