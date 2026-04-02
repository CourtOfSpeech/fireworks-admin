package teltent

import "github.com/google/wire"

// ProviderSet 租户模块依赖提供者集合。
var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
)
