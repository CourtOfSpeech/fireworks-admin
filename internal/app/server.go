package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/speech/fireworks-admin/internal/middleware"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	"github.com/speech/fireworks-admin/internal/pkg/config"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
	"github.com/speech/fireworks-admin/internal/pkg/validator"
)

// Server HTTP 服务器。
type Server struct {
	echo    *echo.Echo
	httpSrv *http.Server
	app     *App
	cleanup func()
}

// NewServer 创建 HTTP 服务器实例并注册所有路由。
func NewServer(a *App, cleanup func()) *Server {
	e := newEcho(a)
	setupMiddleware(e, a.Config)
	RegisterRoutes(e, a.Registrars)
	httpSrv := newHTTPServer(e, a.Config)

	return &Server{
		echo:    e,
		httpSrv: httpSrv,
		app:     a,
		cleanup: cleanup,
	}
}

// newEcho 创建并配置 Echo 实例。
func newEcho(a *App) *echo.Echo {
	return echo.NewWithConfig(echo.Config{
		Logger:             a.Logger,
		HTTPErrorHandler:   customHTTPErrorHandler,
		Validator:          validator.NewValidator(),
		FormParseMaxMemory: 10 << 20,
	})
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

// setupMiddleware 配置全局中间件。
func setupMiddleware(e *echo.Echo, cfg *config.Config) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS(cfg.Server.AllowOrigins))
	e.Use(middleware.Timeout(cfg.Server.Timeout))
}

// newHTTPServer 创建 HTTP 服务器实例。
func newHTTPServer(e *echo.Echo, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        e,
		ReadTimeout:    durationOrDefault(cfg.Server.ReadTimeout, 15*time.Second),
		WriteTimeout:   durationOrDefault(cfg.Server.WriteTimeout, 15*time.Second),
		IdleTimeout:    durationOrDefault(cfg.Server.IdleTimeout, 60*time.Second),
		MaxHeaderBytes: bytesOrDefault(cfg.Server.MaxHeaderBytes, 1<<20),
	}
}

// durationOrDefault 返回配置的持续时间或默认值。
func durationOrDefault(seconds int, defaultValue time.Duration) time.Duration {
	if seconds <= 0 {
		return defaultValue
	}
	return time.Duration(seconds) * time.Second
}

// bytesOrDefault 返回配置的字节数或默认值。
func bytesOrDefault(bytes int, defaultValue int) int {
	if bytes <= 0 {
		return defaultValue
	}
	return bytes
}

// Run 启动服务器并阻塞运行，直到收到终止信号。
// 内部处理信号监听和优雅关闭，简化调用方代码。
func (s *Server) Run() error {
	logger.Info("server starting", slog.Int("port", s.app.Config.Server.Port))

	serverErr := make(chan error, 1)
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server start failed: %w", err)
	case sig := <-quit:
		logger.Info("received shutdown signal", slog.String("signal", sig.String()))
	}

	return s.gracefulShutdown()
}

// gracefulShutdown 优雅关闭服务器。
func (s *Server) gracefulShutdown() error {
	logger.Info("shutting down server...")

	timeout := durationOrDefault(s.app.Config.Server.ShutdownTimeout, 30*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	s.cleanup()
	logger.Info("server stopped gracefully")
	return nil
}
