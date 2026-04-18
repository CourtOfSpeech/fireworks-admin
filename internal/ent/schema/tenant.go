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

// Tenant 表示租户实体的 Schema 定义。
// 租户可以是企业或个人，包含证件信息、联系方式等基本信息。
type Tenant struct {
	ent.Schema
}

// Mixin 返回 Tenant 实体使用的 Mixin 列表。
// 该方法组合了公共字段，包括 ID、状态、创建时间、更新时间和软删除字段。
func (Tenant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Id{},
		mixin.Status{},
		mixin.CreateTime{},
		mixin.UpdateTime{},
		mixin.SoftDelete{},
	}
}

// Fields 定义 Tenant 实体的字段列表。
// 包含证件号码、租户名称、类型、联系人信息、邮箱、电话和过期时间等字段。
func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("certificate_no").
			NotEmpty().
			MaxLen(50).
			Comment("证件号码：企业-统一社会信用代码 个人-身份证号"),
		field.String("name").
			NotEmpty().
			MaxLen(100).
			Comment("租户名称"),
		field.Int8("type").
			Default(1).
			Min(1).
			Max(2).
			Comment("租户类型"),
		field.String("contact_name").
			NotEmpty().
			MaxLen(50).
			Comment("联系人姓名"),
		field.String("email").
			NotEmpty().
			MaxLen(255).
			Comment("联系邮箱"),
		field.String("phone").
			NotEmpty().
			MaxLen(20).
			Comment("联系电话"),
		field.Time("expired_at").
			Optional().
			Comment("过期时间"),
	}
}

// Indexes 定义 Tenant 实体的索引列表。
// 使用部分索引（Partial Index）确保唯一约束只对未删除的记录生效，
// 包括邮箱、电话和证件号码的唯一索引。
func (Tenant) Indexes() []ent.Index {
	return []ent.Index{
		// 邮箱唯一索引（仅对未删除记录）
		index.Fields("email").
			Unique().
			StorageKey("uk_email").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
		// 电话唯一索引（仅对未删除记录）
		index.Fields("phone").
			Unique().
			StorageKey("uk_phone").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
		// 证件号唯一索引（仅对未删除记录）
		index.Fields("certificate_no").
			Unique().
			StorageKey("uk_certificate_no").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
	}
}
