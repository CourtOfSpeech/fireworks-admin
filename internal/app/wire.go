//go:build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/speech/fireworks-admin/internal/features/tenant"
)

// InitializeApp 初始化应用依赖。
func InitializeApp() (*App, error) {
	wire.Build(
		ProviderSet,
		tenant.ProviderSet,
		NewHealthRouter,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
