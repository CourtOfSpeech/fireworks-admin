package mixin

import (
	"context"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/speech/fireworks-admin/internal/ent/hook"
	"github.com/speech/fireworks-admin/internal/ent/intercept"
	"github.com/speech/fireworks-admin/internal/pkg/ctxutil"
	"github.com/speech/fireworks-admin/internal/pkg/idgen"
)

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

// Interceptors 返回 TenantId Mixin 的拦截器列表。
// 包含一个查询拦截器，用于在查询时自动添加 WHERE tenant_id = ? 条件。
// 这确保了每个查询只返回当前租户的数据，防止越权操作。
func (TenantId) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
			tenantID, ok := ctx.Value(ctxutil.TenantKey{}).(uuid.UUID)
			if !ok {
				return fmt.Errorf("安全拦截: 上下文中缺失 tenant_id")
			}

			if w, ok := q.(interface{ WhereP(...func(*sql.Selector)) }); ok {
				w.WhereP(sql.FieldEQ("tenant_id", tenantID))
			}
			return nil
		}),
	}
}

// Hooks 返回 TenantId Mixin 的钩子列表。
// 包含一个钩子，用于在写操作时自动注入租户ID。
// 这确保了每个写操作都只影响当前租户的数据，防止越权操作。
func (TenantId) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
					tenantID, ok := ctx.Value(ctxutil.TenantKey{}).(uuid.UUID)
					if !ok {
						return nil, fmt.Errorf("安全拦截: 变更操作必须携带 tenant_id")
					}

					// 如果是创建，自动注入租户 ID
					if m.Op().Is(ent.OpCreate) {
						type tenantMutator interface {
							SetTenantID(uuid.UUID)
						}
						if mx, ok := m.(tenantMutator); ok {
							mx.SetTenantID(tenantID)
						}
					}

					// 如果是更新或删除，自动追加 WHERE 条件，防止越权操作别人的数据
					if m.Op().Is(ent.OpUpdateOne | ent.OpUpdate | ent.OpDeleteOne | ent.OpDelete) {
						if w, ok := m.(interface{ WhereP(...func(*sql.Selector)) }); ok {
							w.WhereP(sql.FieldEQ("tenant_id", tenantID))
						}
					}

					return next.Mutate(ctx, m)
				})
			},
			// 应用于所有写操作
			ent.OpCreate|ent.OpUpdateOne|ent.OpUpdate|ent.OpDeleteOne|ent.OpDelete,
		),
	}
}
