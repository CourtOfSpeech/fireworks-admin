package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/speech/fireworks-admin/internal/infrastructure/persistence/ent/mixin"
)

// Teltent 租户实体
type Teltent struct {
	ent.Schema
}

// 使用公共字段 Mixin
func (Teltent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Id{},
		mixin.Status{},
		mixin.CreateTime{},
		mixin.UpdateTime{},
		mixin.SoftDelete{},
	}
}

// Fields 字段
func (Teltent) Fields() []ent.Field {
	return []ent.Field{
		field.String("certificate_no").
			NotEmpty().
			MaxLen(50).
			Unique().
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

// Indexes 索引
// 电话 和 ID 的联合唯一索引
// 邮箱和 ID 的联合唯一索引
func (Teltent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("phone", "id").StorageKey("idx_phone_id"),
		index.Fields("email", "id").StorageKey("idx_email_id"),
	}
}
