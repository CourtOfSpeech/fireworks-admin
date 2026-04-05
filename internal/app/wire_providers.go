package app

import (
	"reflect"

	"github.com/google/wire"
	"github.com/speech/fireworks-admin/internal/features/tenant"
	"github.com/speech/fireworks-admin/internal/pkg/config"
	"github.com/speech/fireworks-admin/internal/pkg/db"
	"github.com/speech/fireworks-admin/internal/pkg/lifecycle"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
)

// ProviderSet app依赖提供者集合。
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

// Registrars 内部结构体，用于收集所有 RouterRegistrar 实现。
type RegistrarIn struct {
	Tenant *tenant.Handler
	Health *HealthRouter
}

// ProvideRegistrars 将所有 RouterRegistrar 实现收集为切片。
func ProvideRegistrars(r RegistrarIn) []RouterRegistrar {
	var registrars []RouterRegistrar
	v := reflect.ValueOf(r)

	// 防御性检查：确保传入的是结构体
	if v.Kind() != reflect.Struct {
		return nil
	}

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
