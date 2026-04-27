package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

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
