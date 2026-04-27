package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/speech/fireworks-admin/internal/pkg/idgen"
)

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
