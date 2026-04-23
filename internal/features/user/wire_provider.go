// Package user 提供User功能，包括User的创建、查询、更新和删除操作。
// 本文件定义了User模块的 Wire 依赖注入提供者集合。
package user

import "github.com/google/wire"

// ProviderSet User模块依赖提供者集合。
// 包含User模块所有需要注入的组件：Repository、Service 和 Handler。
var ProviderSet = wire.NewSet(
	NewUserRepo,
	NewUserService,
	NewUserHandler,
)
