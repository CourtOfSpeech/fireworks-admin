package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

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
