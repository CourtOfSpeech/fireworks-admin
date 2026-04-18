// Package tenant 提供租户管理功能，包括租户的创建、查询、更新和删除操作。
// 本文件定义了租户模块的 Wire 依赖注入提供者集合。
package tenant

import "github.com/google/wire"

// ProviderSet 租户模块依赖提供者集合。
// 包含租户模块所有需要注入的组件：Repository、Service 和 Handler。
var ProviderSet = wire.NewSet(
	NewTenantRepo,
	NewTenantService,
	NewTenantHandler,
)
