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

// Run 封装了所有的启动、监听信号、优雅退出的繁琐代码。
// 该函数是应用程序的主入口点，负责以下工作：
//  1. 初始化应用依赖容器
//  2. 启动所有生命周期组件（如 HTTP 服务器）
//  3. 监听系统信号（SIGINT, SIGTERM）
//  4. 接收到退出信号后执行优雅关闭
//  5. 等待所有组件安全退出
// 如果初始化或启动失败，程序将以非零状态码退出。
func Run() {
	a, err := InitializeApp()
	if err != nil {
		logger.Error(context.Background(), "应用初始化失败", slog.Any("error", err))
		os.Exit(1)
	}

	startTimeout := time.Duration(a.Config.Server.StartTimeout) * time.Second
	startCtx, cancelStart := context.WithTimeout(context.Background(), startTimeout)
	defer cancelStart()

	if err := a.Lifecycle.Start(startCtx); err != nil {
		logger.Error(context.Background(), "应用启动失败", slog.Any("error", err))
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	logger.Info(context.Background(), "接收到退出信号", slog.String("signal", s.String()))

	shutdownTimeout := time.Duration(a.Config.Server.ShutdownTimeout) * time.Second
	stopCtx, cancelStop := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelStop()

	a.Lifecycle.Stop(stopCtx)

	logger.Info(context.Background(), "应用已安全退出")
}
