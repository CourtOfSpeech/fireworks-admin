package app

import (
	"reflect"

	"github.com/google/wire"
	"github.com/speech/fireworks-admin/internal/features/tenant"
	"github.com/speech/fireworks-admin/internal/features/user"
	"github.com/speech/fireworks-admin/internal/pkg/config"
	"github.com/speech/fireworks-admin/internal/pkg/db"
	"github.com/speech/fireworks-admin/internal/pkg/lifecycle"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
)

// 确保实现 RouterRegistrar 接口。
var (
	_ RouterRegistrar = (*tenant.TenantHandler)(nil)
	_ RouterRegistrar = (*user.UserHandler)(nil)
)

// ProviderSet app 依赖提供者集合。
// 该集合包含了应用程序核心组件的所有提供者，用于 Wire 依赖注入。
// 包含的提供者：
//   - config.ProviderSet: 配置相关提供者
//   - logger.ProviderSet: 日志相关提供者
//   - db.ProviderSet: 数据库相关提供者
//   - NewServer: HTTP 服务器提供者
//   - NewEcho: Echo 框架提供者
//   - lifecycle.NewLifecycle: 生命周期管理器提供者
//   - wire.Struct(new(RegistrarIn), "*"): 路由注册器输入结构体提供者
//   - ProvideRegistrars: 路由注册器列表提供者
var ProviderSet = wire.NewSet(
	config.ProviderSet,
	logger.ProviderSet,
	db.ProviderSet,
	NewServer,
	NewEcho,
	lifecycle.NewLifecycle,
	wire.Struct(new(RegistrarIn), "*"),
	ProvideRegistrars,
)

// RegistrarIn 内部结构体，用于收集所有 RouterRegistrar 实现。
// 该结构体的字段由 Wire 自动注入，每个字段都是一个实现了 RouterRegistrar 接口的处理器。
// 通过 ProvideRegistrars 函数将所有字段收集为路由注册器列表。
type RegistrarIn struct {
	Tenant *tenant.TenantHandler // 租户管理路由处理器
	User   *user.UserHandler    // 用户管理路由处理器
	Health *HealthRouter        // 健康检查路由处理器
}

// ProvideRegistrars 将所有 RouterRegistrar 实现收集为切片。
// 该函数使用反射遍历 RegistrarIn 结构体的所有字段，
// 自动收集所有实现了 RouterRegistrar 接口的非空字段。
// 这种设计使得添加新的路由注册器时无需修改此函数，只需在 RegistrarIn 中添加字段即可。
// r 是包含所有路由注册器的输入结构体。
func ProvideRegistrars(r RegistrarIn) []RouterRegistrar {
	var registrars []RouterRegistrar
	v := reflect.ValueOf(r)

	// 防御性检查：确保传入的是结构体
	if v.Kind() != reflect.Struct {
		return nil
	}

	// 遍历结构体的所有字段
	// (注意：如果 Fields() 返回的是迭代器 iter.Seq，这里可能只需 for field := range v.Fields())
	for _, field := range v.Fields() {
		// 1. 安全门禁：跳过未导出的私有字段（比如小写开头的 tenant）
		if !field.CanInterface() {
			continue
		}

		// 2. nil 检查：只有特定类型才能调用 IsNil，否则会 Panic
		k := field.Kind()
		isNillable := k == reflect.Pointer || k == reflect.Interface || k == reflect.Slice ||
			k == reflect.Map || k == reflect.Chan || k == reflect.Func

		if isNillable && field.IsNil() {
			continue
		}

		// 3. 核心逻辑：类型断言并收集接口实现
		if registrar, ok := field.Interface().(RouterRegistrar); ok {
			registrars = append(registrars, registrar)
		}
	}

	return registrars
}
