// Package tenant 提供租户管理功能，包括租户的创建、查询、更新和删除操作。
// 本文件定义了租户模块的数据持久化层，负责与数据库交互。
package tenant

import (
	"context"
	"fmt"
	"time"

	entgo "github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/ent/tenant"
	"github.com/speech/fireworks-admin/internal/pkg/db"
	"github.com/speech/fireworks-admin/internal/pkg/idgen"
)

// TenantRepo 租户数据持久化操作的具体实现。
// 封装了租户相关的数据库操作，包括增删改查。
type TenantRepo struct {
	tx *db.TxManager // 数据库事务管理器
}

// NewTenantRepo 创建租户 Repository 实例。
// 参数 txManager 为数据库事务管理器，返回初始化后的 Repository 实例。
func NewTenantRepo(txManager *db.TxManager) *TenantRepo {
	return &TenantRepo{
		tx: txManager,
	}
}

// toEntity 将 Ent 框架的 Tenant 模型转换为领域模型 Tenant。
// 参数 t 为 Ent 框架的租户模型，返回领域模型的租户实体。
func toEntity(t *entgo.Tenant) *Tenant {
	return &Tenant{
		ID:            idgen.ToString(t.ID),
		CertificateNo: t.CertificateNo,
		Name:          t.Name,
		Type:          t.Type,
		ContactName:   t.ContactName,
		Email:         t.Email,
		Phone:         t.Phone,
		ExpiredAt:     t.ExpiredAt,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}

// Create 创建新租户记录。
// 参数 ctx 为上下文，req 为创建请求参数。
// 返回创建成功的租户实体和可能的错误。
func (r *TenantRepo) Create(ctx context.Context, req *CreateTenantReq) (*Tenant, error) {
	builder := r.tx.DB(ctx).Tenant.Create().
		SetCertificateNo(req.CertificateNo).
		SetName(req.Name).
		SetType(req.Type).
		SetContactName(req.ContactName).
		SetEmail(req.Email).
		SetPhone(req.Phone).
		SetStatus(req.Status)

	if !req.ExpiredAt.IsZero() {
		builder.SetExpiredAt(req.ExpiredAt)
	}

	t, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo:Create: %w", err)
	}

	return toEntity(t), nil
}

// Delete 根据租户 ID 软删除租户记录。
// 参数 ctx 为上下文，id 为租户 ID 字符串。
// 软删除通过设置 deleted_at 字段实现，返回删除操作可能发生的错误。
func (r *TenantRepo) Delete(ctx context.Context, id string) error {
	tenantId, err := idgen.Parse(id)
	if err != nil {
		return fmt.Errorf("repo:Delete id parse id=%s: %w", id, err)
	}

	err = r.tx.DB(ctx).Tenant.UpdateOneID(tenantId).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("repo:Delete id=%s: %w", id, err)
	}
	return nil
}

// List 根据查询条件获取租户列表。
// 参数 ctx 为上下文，query 为查询条件。
// 返回租户列表、总数和可能的错误。支持分页和条件过滤。
func (r *TenantRepo) List(ctx context.Context, query *TenantQuery) ([]*Tenant, int64, error) {
	builder := r.tx.DB(ctx).Tenant.Query()

	if query.HasKeyword() {
		builder.Where(
			tenant.Or(
				tenant.CertificateNoHasPrefix(query.Keyword),
				tenant.NameContains(query.Keyword),
			),
		)
	}

	if query.HasStatus() {
		builder.Where(tenant.StatusEQ(*query.Status))
	}
	if query.HasEmail() {
		builder.Where(tenant.EmailEQ(query.Email))
	}
	if query.HasPhone() {
		builder.Where(tenant.PhoneEQ(query.Phone))
	}

	total, err := builder.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo:List count: %w", err)
	}

	tenants, err := builder.
		Offset(query.GetOffset()).
		Limit(query.GetLimit()).
		Order(entgo.Desc(tenant.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo:List query: %w", err)
	}

	result := make([]*Tenant, 0, len(tenants))
	for _, t := range tenants {
		result = append(result, toEntity(t))
	}

	return result, int64(total), nil
}

// GetByID 根据租户 ID 获取租户详情。
// 参数 ctx 为上下文，id 为租户 ID 字符串。
// 返回租户实体和可能的错误。
func (r *TenantRepo) GetByID(ctx context.Context, id string) (*Tenant, error) {
	tenantId, err := idgen.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("repo:GetByID id parse id=%s: %w", id, err)
	}

	t, err := r.tx.DB(ctx).Tenant.Get(ctx, tenantId)
	if err != nil {
		return nil, fmt.Errorf("repo:GetByID id=%s: %w", id, err)
	}
	return toEntity(t), nil
}

// Update 根据租户 ID 更新租户信息。
// 参数 ctx 为上下文，id 为租户 ID 字符串，req 为更新请求参数。
// 仅更新请求中非空字段。返回更新后的租户实体和可能的错误。
func (r *TenantRepo) Update(ctx context.Context, id string, req *UpdateTenantReq) (*Tenant, error) {
	tenantId, err := idgen.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("repo:Update id parse id=%s: %w", id, err)
	}

	builder := r.tx.DB(ctx).Tenant.UpdateOneID(tenantId)
	if req.CertificateNo != nil {
		builder.SetCertificateNo(*req.CertificateNo)
	}
	if req.Name != nil {
		builder.SetName(*req.Name)
	}
	if req.Type != nil {
		builder.SetType(*req.Type)
	}
	if req.ContactName != nil {
		builder.SetContactName(*req.ContactName)
	}
	if req.Email != nil {
		builder.SetEmail(*req.Email)
	}
	if req.Phone != nil {
		builder.SetPhone(*req.Phone)
	}
	if req.ExpiredAt != nil {
		builder.SetExpiredAt(*req.ExpiredAt)
	}
	if req.Status != nil {
		builder.SetStatus(*req.Status)
	}

	t, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo:Update id=%s: %w", id, err)
	}
	return toEntity(t), nil
}
