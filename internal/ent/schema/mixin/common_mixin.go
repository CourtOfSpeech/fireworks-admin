// Package mixin 提供了 Ent Schema 的公共字段 Mixin。
// 该包定义了可复用的字段组合，包括 ID、租户ID、状态、时间戳和软删除等常用字段。
package mixin

import (
	"slices"

	"entgo.io/ent"
)

// CommonMixin 是一个组合型 Mixin，包含常用的公共字段。
// 该 Mixin 整合了 ID、租户ID、状态、创建时间、更新时间和软删除字段，
// 适用于需要完整审计和租户隔离功能的实体。
type CommonMixin struct {
	ent.Schema
}

// Fields 返回 CommonMixin 包含的所有字段列表。
// 该方法通过组合其他 Mixin 的字段来实现字段的复用。
func (CommonMixin) Fields() []ent.Field {
	return slices.Concat(
		Id{}.Fields(),
		TenantId{}.Fields(),
		Status{}.Fields(),
		CreateTime{}.Fields(),
		UpdateTime{}.Fields(),
		SoftDelete{}.Fields(),
	)
}

// commonMixin 实现了 ent.Mixin 接口。
var _ ent.Mixin = (*CommonMixin)(nil)
