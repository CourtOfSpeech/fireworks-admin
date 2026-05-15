// Package user 提供用户管理功能，包括用户的创建、查询、更新和删除操作。
// 本文件定义了用户模块的 Wire 依赖注入提供者集合。
package user

import "github.com/google/wire"

// ProviderSet 用户模块依赖提供者集合。
// 包含用户模块所有需要注入的组件：Repository、Service 和 Handler。
var ProviderSet = wire.NewSet(
	NewUserRepo,
	NewUserService,
	NewUserHandler,
)
