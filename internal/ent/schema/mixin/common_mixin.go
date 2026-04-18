// Package mixin 提供了 Ent Schema 的公共字段 Mixin。
// 该包定义了可复用的字段组合，包括 ID、租户ID、状态、时间戳和软删除等常用字段。
package mixin

import (
	"slices"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/speech/fireworks-admin/internal/pkg/idgen"
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

// TenantId 是一个提供租户ID字段的 Mixin。
// 用于实现多租户架构中的租户隔离功能。
type TenantId struct{ ent.Schema }

// Fields 返回 TenantId Mixin 的字段列表。
// 包含一个 UUID v7 格式的租户ID字段，该字段不可变。
func (TenantId) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("tenant_id", uuid.UUID{}).
			Default(idgen.NewV7Safe).
			Immutable().
			Comment("租户ID"),
	}
}

// tenantId 实现了 ent.Mixin 接口。
var _ ent.Mixin = (*TenantId)(nil)

// Id 是一个提供主键ID字段的 Mixin。
// 使用 UUID v7 作为主键，确保全局唯一性和时间有序性。
type Id struct{ ent.Schema }

// Fields 返回 Id Mixin 的字段列表。
// 包含一个 UUID v7 格式的主键字段，该字段唯一且不可变。
func (Id) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(idgen.NewV7Safe).
			Unique().
			Immutable().
			Comment("主键"),
	}
}

// id 实现了 ent.Mixin 接口。
var _ ent.Mixin = (*Id)(nil)

// SoftDelete 是一个提供软删除功能的 Mixin。
// 通过 deleted_at 字段实现软删除，当该字段不为空时表示记录已被删除。
type SoftDelete struct{ ent.Schema }

// Fields 返回 SoftDelete Mixin 的字段列表。
// 包含一个可选的删除时间字段，用于标记记录是否被软删除。
func (SoftDelete) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Optional().
			Comment("删除时间，非空表示已软删除"),
	}
}

// softDelete 实现了 ent.Mixin 接口。
var _ ent.Mixin = (*SoftDelete)(nil)

// Status 是一个提供状态字段的 Mixin。
// 用于记录实体的当前状态，如启用、禁用等。
type Status struct{ ent.Schema }

// Fields 返回 Status Mixin 的字段列表。
// 包含一个 int8 类型的状态字段。
func (Status) Fields() []ent.Field {
	return []ent.Field{
		field.Int8("status").
			Comment("状态"),
	}
}

// status 实现了 ent.Mixin 接口。
var _ ent.Mixin = (*Status)(nil)

// CreateTime 是一个提供创建时间字段的 Mixin。
// 自动记录实体的创建时间，该字段在创建后不可变。
type CreateTime struct{ ent.Schema }

// Fields 返回 CreateTime Mixin 的字段列表。
// 包含一个自动设置为当前时间的创建时间字段，该字段不可变。
func (CreateTime) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
	}
}

// createTime 实现了 ent.Mixin 接口。
var _ ent.Mixin = (*CreateTime)(nil)

// UpdateTime 是一个提供更新时间字段的 Mixin。
// 自动记录实体的最后更新时间，每次更新时自动刷新。
type UpdateTime struct{ ent.Schema }

// Fields 返回 UpdateTime Mixin 的字段列表。
// 包含一个自动设置为当前时间的更新时间字段，每次更新时自动刷新。
func (UpdateTime) Fields() []ent.Field {
	return []ent.Field{
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
	}
}

// updateTime 实现了 ent.Mixin 接口。
var _ ent.Mixin = (*UpdateTime)(nil)
