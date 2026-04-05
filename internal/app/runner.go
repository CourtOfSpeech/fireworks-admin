package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/speech/fireworks-admin/internal/pkg/logger"
)

// Run 封装了所有的启动、监听信号、优雅退出的繁琐代码
func Run() {
	a, err := InitializeApp()
	if err != nil {
		logger.Error("应用初始化失败", slog.Any("error", err))
		os.Exit(1)
	}

	startTimeout := time.Duration(a.Config.Server.StartTimeout) * time.Second
	startCtx, cancelStart := context.WithTimeout(context.Background(), startTimeout)
	defer cancelStart()

	if err := a.Lifecycle.Start(startCtx); err != nil {
		logger.Error("应用启动失败", slog.Any("error", err))
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	logger.Info("接收到退出信号", slog.String("signal", s.String()))

	shutdownTimeout := time.Duration(a.Config.Server.ShutdownTimeout) * time.Second
	stopCtx, cancelStop := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelStop()

	a.Lifecycle.Stop(stopCtx)

	logger.Info("应用已安全退出")
}
