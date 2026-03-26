//go:build wireinject

package di

import (
	"log/slog"

	"github.com/google/wire"
	"github.com/speech/fireworks-admin/internal/di/provider"
	appconfig "github.com/speech/fireworks-admin/internal/infrastructure/config"
	apphandle "github.com/speech/fireworks-admin/internal/infrastructure/http/handle"
	"github.com/speech/fireworks-admin/internal/infrastructure/persistence/ent"
)

// App 应用依赖容器
// 包含所有需要注入的依赖组件
type App struct {
	Config        *appconfig.Config
	Logger        *slog.Logger
	EntClient     *ent.Client
	TeltentHandle *apphandle.TeltentHandle
}

// InitializeApp 初始化应用依赖
// 该函数由 wire 自动生成实现，用于创建所有依赖组件
// 返回:
//   - *App: 应用依赖容器
//   - func(): 清理函数
//   - error: 初始化错误
func InitializeApp() (*App, func(), error) {
	wire.Build(
		provider.ProvideConfig,
		provider.ProvideLogger,
		provider.ProvideDSN,
		provider.ProvideEntClient,
		provider.ProvideTeltentRepo,
		provider.ProvideTeltentUsecase,
		provider.ProvideTeltentHandle,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
