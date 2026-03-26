package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/di"
	"github.com/speech/fireworks-admin/internal/infrastructure/http/middleware"
	"github.com/speech/fireworks-admin/internal/infrastructure/http/router"
	"github.com/speech/fireworks-admin/pkg/logger"
	"github.com/speech/fireworks-admin/pkg/response"
	"github.com/speech/fireworks-admin/pkg/validate"
)

type Server struct {
	echo    *echo.Echo
	app     *di.App
	cleanup func()
}

func NewServer() (*Server, error) {
	app, cleanup, err := di.InitializeApp()
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		return nil, err
	}

	e := echo.NewWithConfig(echo.Config{
		Logger:             app.Logger,
		HTTPErrorHandler:   customHTTPErrorHandler,
		Validator:          validate.NewValidator(),
		FormParseMaxMemory: 10 << 20,
	})

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS(app.Config.Server.AllowOrigins))
	e.Use(middleware.Timeout(app.Config.Server.Timeout))

	router.RegisterRoutes(e, app)

	return &Server{
		echo:    e,
		app:     app,
		cleanup: cleanup,
	}, nil
}

func (s *Server) Start() error {
	logger.Info("server starting", slog.Int("port", s.app.Config.Server.Port))
	if err := s.echo.Start(fmt.Sprintf(":%d", s.app.Config.Server.Port)); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

func (s *Server) Close() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

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
	_ = response.Error(c, code, message)
}
