package app

import (
	"log/slog"

	"github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/pkg/config"
)

// App 应用依赖容器。
type App struct {
	Config     *config.Config
	Logger     *slog.Logger
	EntClient  *ent.Client
	Registrars []RouterRegistrar
}
