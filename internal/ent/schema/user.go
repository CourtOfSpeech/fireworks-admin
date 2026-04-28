// Package schema 定义了数据库实体的 Schema 结构。
// 该包使用 Ent ORM 框架来定义数据库表结构、字段、索引等元数据。
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/speech/fireworks-admin/internal/ent/schema/mixin"
)

// User 表示用户实体的 Schema 定义。
// 用户包含用户名、邮箱、手机号、密码等基本信息。
type User struct {
	ent.Schema
}

// Mixin 返回 User 实体使用的 Mixin 列表。
// 该方法组合了公共字段，包括 ID、租户ID、状态、创建时间、更新时间和软删除字段。
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Id{},
		mixin.TenantId{},
		mixin.Status{},
		mixin.CreateTime{},
		mixin.UpdateTime{},
		mixin.SoftDelete{},
	}
}

// Fields 定义 User 实体的字段列表。
// 包含用户名、邮箱、手机号、密码、昵称和头像URL等字段。
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			NotEmpty().
			MaxLen(50).
			Comment("用户名"),
		field.String("email").
			NotEmpty().
			MaxLen(255).
			Comment("邮箱"),
		field.String("phone").
			NotEmpty().
			MaxLen(20).
			Comment("手机号"),
		field.String("password").
			NotEmpty().
			MaxLen(255).
			Comment("密码"),
		field.String("nickname").
			Optional().
			MaxLen(50).
			Comment("昵称"),
		field.String("avatar").
			Optional().
			MaxLen(500).
			Comment("头像URL"),
	}
}

// Indexes 定义 User 实体的索引列表。
// 使用部分索引（Partial Index）确保唯一约束只对未删除的记录生效，
// 包括用户名、邮箱和手机号的唯一索引。
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// 用户名唯一索引（仅对未删除记录）
		index.Fields("username").
			Unique().
			StorageKey("uk_username").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
		// 邮箱唯一索引（仅对未删除记录）
		index.Fields("email").
			Unique().
			StorageKey("uk_email").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
		// 手机号唯一索引（仅对未删除记录）
		index.Fields("phone").
			Unique().
			StorageKey("uk_phone").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
	}
}
