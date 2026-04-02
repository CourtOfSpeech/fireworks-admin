package mixin

import (
	"slices"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/speech/fireworks-admin/internal/pkg/util"
)

// CommonMixin 通用 Mixin。
// 该 Mixin 提供 id 字段（UUID v7）和 created_at、updated_at 字段，用于记录实体的创建和更新时间。
// 该 Mixin 还提供了 status 字段，用于记录实体的状态。
// 该 Mixin 还提供了 deleted_at 字段，用于实现软删除功能。
// 该 Mixin 还提供了 tenant_id 字段，用于记录实体所属租户。
type CommonMixin struct {
	ent.Mixin
}

// Fields of the common mixin.
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

// common mixin must implement `Mixin` interface.
var _ ent.Mixin = (*CommonMixin)(nil)

// TenantId adds tenant id field.
type TenantId struct{ ent.Schema }

// Fields of the tenant id mixin.
func (TenantId) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("tenant_id", uuid.UUID{}).
			Default(utils.NewV7Safe).
			Immutable().
			Comment("租户ID"),
	}
}

// tenant id mixin must implement `Mixin` interface.
var _ ent.Mixin = (*TenantId)(nil)

// Id adds id field.
type Id struct{ ent.Schema }

// Fields of the id mixin.
func (Id) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(utils.NewV7Safe).
			Unique().
			Immutable().
			Comment("主键"),
	}
}

// id mixin must implement `Mixin` interface.
var _ ent.Mixin = (*Id)(nil)

// SoftDelete adds soft delete field.
// 该 Mixin 提供 deleted_at 字段，用于实现软删除功能。
// 当 deleted_at 不为空时，表示该记录已被软删除。
type SoftDelete struct{ ent.Schema }

func (s SoftDelete) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Optional().
			Comment("删除时间，非空表示已软删除"),
	}
}

// soft delete mixin must implement `Mixin` interface.
var _ ent.Mixin = (*SoftDelete)(nil)

// Status adds created at time field.
type Status struct{ ent.Schema }

// Fields of the status mixin.
func (Status) Fields() []ent.Field {
	return []ent.Field{
		field.Int8("status").
			Comment("状态"),
	}
}

// status mixin must implement `Mixin` interface.
var _ ent.Mixin = (*Status)(nil)

// CreateTime adds created at time field.
type CreateTime struct{ ent.Schema }

// Fields of the create time mixin.
func (CreateTime) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
	}
}

// create time mixin must implement `Mixin` interface.
var _ ent.Mixin = (*CreateTime)(nil)

// UpdateTime adds updated at time field.
type UpdateTime struct{ ent.Schema }

// Fields of the update time mixin.
func (UpdateTime) Fields() []ent.Field {
	return []ent.Field{
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
	}
}

// update time mixin must implement `Mixin` interface.
var _ ent.Mixin = (*UpdateTime)(nil)
