//go:build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/speech/fireworks-admin/internal/features/teltent"
	"github.com/speech/fireworks-admin/internal/pkg/config"
	"github.com/speech/fireworks-admin/internal/pkg/db"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
)

// InitializeApp 初始化应用依赖。
func InitializeApp() (*App, func(), error) {
	wire.Build(
		config.ProviderSet,
		logger.ProviderSet,
		db.ProviderSet,
		teltent.ProviderSet,
		NewHealthRouter,

		wire.Struct(new(RegistrarIn), "*"),
		ProvideRegistrars,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
