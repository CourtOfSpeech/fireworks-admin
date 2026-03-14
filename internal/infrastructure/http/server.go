package http

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

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
	log    *slog.Logger
}

func NewServer() (*Server, error) {
	cfg, err := config.LoadByEnv()
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	log := logger.NewLogger(cfg.Log.Level, cfg.Log.Format, cfg.Log.AddSource)

	e := echo.NewWithConfig(echo.Config{
		Logger:           log,
		HTTPErrorHandler: customHTTPErrorHandler,
		Validator:        validate.NewValidator(),
	})

	return &Server{
		echo:   e,
		config: cfg,
		log:    log,
	}, nil
}

func (s *Server) Start() error {
	logger.Info("server starting", slog.Int("port", s.config.Server.Port))
	if err := s.echo.Start(fmt.Sprintf(":%d", s.config.Server.Port)); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
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
