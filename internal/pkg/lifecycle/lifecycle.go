// Package lifecycle 提供应用程序组件生命周期管理功能。
// 它实现了组件的顺序启动和逆序停止机制，确保资源正确初始化和清理。
package lifecycle

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/speech/fireworks-admin/internal/pkg/logger"
)

// Hook 定义组件的启动和停止逻辑。
// Hook 用于封装单个组件的生命周期钩子函数，包含组件名称以及启动和停止时的回调函数。
type Hook struct {
	Name    string                        // 组件名称，用于日志和错误信息标识
	OnStart func(ctx context.Context) error // 启动回调函数，在组件启动时调用
	OnStop  func(ctx context.Context) error // 停止回调函数，在组件停止时调用
}

// Lifecycle 管理所有组件的生命周期。
// Lifecycle 维护一个钩子列表，提供组件的注册、顺序启动和逆序停止功能。
// 它使用互斥锁保证并发安全，适用于应用程序的优雅启动和关闭场景。
type Lifecycle struct {
	hooks []Hook     // 已注册的生命周期钩子列表
	mu    sync.Mutex // 保护 hooks 并发访问的互斥锁
}

// NewLifecycle 创建生命周期管理器。
// 返回一个初始化的 Lifecycle 实例，可用于注册和管理组件生命周期。
func NewLifecycle() *Lifecycle {
	return &Lifecycle{}
}

// Append 注册新钩子。
// 将 Hook 添加到生命周期管理器的钩子列表中，支持并发安全调用。
func (l *Lifecycle) Append(h Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks = append(l.hooks, h)
}

// Start 顺序启动所有组件。
// 按照注册顺序依次调用每个钩子的 OnStart 函数。
// 如果某个组件启动失败，立即返回错误并停止后续组件的启动。
func (l *Lifecycle) Start(ctx context.Context) error {
	for _, h := range l.hooks {
		if h.OnStart != nil {
			logger.Info(ctx, "正在启动组件", slog.String("name", h.Name))
			if err := h.OnStart(ctx); err != nil {
				return fmt.Errorf("组件 [%s] 启动失败: %w", h.Name, err)
			}
		}
	}
	return nil
}

// Stop 逆序停止所有组件。
// 按照注册的逆序依次调用每个钩子的 OnStop 函数。
// 即使某个组件停止失败，也会继续停止其他组件，错误仅记录日志。
func (l *Lifecycle) Stop(ctx context.Context) {
	for i := len(l.hooks) - 1; i >= 0; i-- {
		h := l.hooks[i]
		if h.OnStop != nil {
			logger.Info(ctx, "正在停止组件", slog.String("name", h.Name))
			if err := h.OnStop(ctx); err != nil {
				logger.Error(ctx, "组件停止失败", slog.String("name", h.Name), slog.Any("error", err))
			}
		}
	}
}
