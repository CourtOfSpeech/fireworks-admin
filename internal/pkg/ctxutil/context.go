// Package ctxutil 提供了 context 相关的工具函数。
// 该包封装了常用的 context 操作，如请求 ID、租户 ID 的存取和软删除控制等。
package ctxutil

import (
	"context"

	"github.com/google/uuid"
)

// RequestIDKey 请求 ID 的 context 键类型。
type RequestIDKey struct{}

// SetRequestID 将请求 ID 存储到 context 中。
// 参数 ctx 是原始的 context.Context。
// 参数 requestID 是要存储的请求 ID 字符串。
// 返回一个新的 context.Context，其中包含了请求 ID。
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey{}, requestID)
}

// GetRequestID 从 context 中获取请求 ID。
// 参数 ctx 是包含请求 ID 的 context.Context。
// 返回存储的请求 ID 字符串，如果不存在则返回空字符串。
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// TenantKey 租户 ID 的 context 键类型。
type TenantKey struct{}

// WithTenant 将租户 ID 存储到 context 中。
// 参数 ctx 是原始的 context.Context，tenantID 是租户的 UUID。
// 返回一个新的 context.Context，其中包含了租户 ID。
func WithTenant(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, TenantKey{}, tenantID)
}

// SoftDeleteKey 软删除控制的 context 键类型。
type SoftDeleteKey struct{}

// SkipSoftDelete 设置 context 标记，告诉 ORM 忽略软删除过滤。
// 用于查询已被软删除的数据或执行硬删除操作。
func SkipSoftDelete(ctx context.Context) context.Context {
	return context.WithValue(ctx, SoftDeleteKey{}, true)
}
