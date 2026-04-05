package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/speech/fireworks-admin/internal/ent/schema/mixin"
)

// Tenant 租户实体
type Tenant struct {
	ent.Schema
}

// 使用公共字段 Mixin
func (Tenant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Id{},
		mixin.Status{},
		mixin.CreateTime{},
		mixin.UpdateTime{},
		mixin.SoftDelete{},
	}
}

// Fields 字段
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

// Indexes 定义表的索引。
// 使用部分索引（Partial Index）确保唯一约束只对未删除的记录生效。
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
