package lifecycle

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/speech/fireworks-admin/internal/pkg/logger"
)

// Hook 定义组件的启动和停止逻辑。
type Hook struct {
	Name    string
	OnStart func(ctx context.Context) error
	OnStop  func(ctx context.Context) error
}

// Lifecycle 管理所有组件的生命周期。
type Lifecycle struct {
	hooks []Hook
	mu    sync.Mutex
}

// NewLifecycle 创建生命周期管理器。
func NewLifecycle() *Lifecycle {
	return &Lifecycle{}
}

// Append 注册新钩子。
func (l *Lifecycle) Append(h Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks = append(l.hooks, h)
}

// Start 顺序启动所有组件。
func (l *Lifecycle) Start(ctx context.Context) error {
	for _, h := range l.hooks {
		if h.OnStart != nil {
			logger.Info("正在启动组件", slog.String("name", h.Name))
			if err := h.OnStart(ctx); err != nil {
				return fmt.Errorf("组件 [%s] 启动失败: %w", h.Name, err)
			}
		}
	}
	return nil
}

// Stop 逆序停止所有组件。
func (l *Lifecycle) Stop(ctx context.Context) {
	for i := len(l.hooks) - 1; i >= 0; i-- {
		h := l.hooks[i]
		if h.OnStop != nil {
			logger.Info("正在停止组件", slog.String("name", h.Name))
			if err := h.OnStop(ctx); err != nil {
				logger.Error("组件停止失败", slog.String("name", h.Name), slog.Any("error", err))
			}
		}
	}
}
