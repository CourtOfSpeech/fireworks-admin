package db

import "github.com/google/wire"

// ProviderSet 数据库依赖提供者集合。
var ProviderSet = wire.NewSet(
	NewEntClient,
)
