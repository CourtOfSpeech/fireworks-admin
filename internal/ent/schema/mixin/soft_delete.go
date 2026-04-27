package mixin

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"github.com/speech/fireworks-admin/internal/ent/hook"
	"github.com/speech/fireworks-admin/internal/ent/intercept"
	"github.com/speech/fireworks-admin/internal/pkg/ctxutil"
)

// SoftDelete 是一个提供软删除功能的 Mixin。
// 通过 deleted_at 字段实现软删除，当该字段不为空时表示记录已被删除。
type SoftDelete struct{ ent.Schema }

// Fields 返回 SoftDelete Mixin 的字段列表。
// 包含一个可选的删除时间字段，用于标记记录是否被软删除。
func (SoftDelete) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Optional().
			Comment("删除时间，非空表示已软删除"),
	}
}

// softDelete 实现了 ent.Mixin 接口。
var _ ent.Mixin = (*SoftDelete)(nil)

// Interceptors 返回 SoftDelete Mixin 的拦截器列表。
// 包含一个查询拦截器，用于在查询时自动添加 WHERE deleted_at IS NULL 条件。
// 这确保了每个查询只返回未被删除的记录，防止查询到已删除的数据。
func (d SoftDelete) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
			// 如果 context 里标记了跳过，就不加 deleted_at is null
			if skip, _ := ctx.Value(ctxutil.SoftDeleteKey{}).(bool); skip {
				return nil
			}
			// 自动追加 WHERE deleted_at IS NULL
			if w, ok := q.(interface{ WhereP(...func(*sql.Selector)) }); ok {
				w.WhereP(sql.FieldIsNull("deleted_at"))
			}
			return nil
		}),
	}
}

// Hooks 返回 SoftDelete Mixin 的钩子列表。
// 包含一个钩子，用于在写操作时自动处理软删除。
// 这确保了每个写操作都只影响未被删除的记录，防止删除已删除的数据。
func (d SoftDelete) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
					if skip, _ := ctx.Value(ctxutil.SoftDeleteKey{}).(bool); skip {
						return next.Mutate(ctx, m)
					}

					// 将 DELETE 转换为 UPDATE deleted_at = now()
					if m.Op().Is(ent.OpDeleteOne | ent.OpDelete) {
						type softDelete interface {
							SetOp(ent.Op)
							SetDeletedAt(time.Time)
						}
						if mx, ok := m.(softDelete); ok {
							mx.SetOp(ent.OpUpdate)      // 修改操作类型为 Update
							mx.SetDeletedAt(time.Now()) // 赋值软删时间
						}
					}

					// 防止更新已经软删除的数据
					if w, ok := m.(interface{ WhereP(...func(*sql.Selector)) }); ok {
						w.WhereP(sql.FieldIsNull("deleted_at"))
					}

					return next.Mutate(ctx, m)
				})
			},
			ent.OpDeleteOne|ent.OpDelete|ent.OpUpdateOne|ent.OpUpdate,
		),
	}
}
