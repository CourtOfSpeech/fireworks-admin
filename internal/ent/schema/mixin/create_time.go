package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

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
