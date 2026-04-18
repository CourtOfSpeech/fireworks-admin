//go:build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/speech/fireworks-admin/internal/features/tenant"
)

// InitializeApp 初始化应用依赖。
// 该函数由 Wire 框架自动生成实现代码（见 wire_gen.go）。
// 通过组合多个 ProviderSet 和提供者函数，构建完整的应用依赖图。
// 返回初始化完成的 App 实例或错误。
func InitializeApp() (*App, error) {
	wire.Build(
		ProviderSet,
		tenant.ProviderSet,
		NewHealthRouter,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
