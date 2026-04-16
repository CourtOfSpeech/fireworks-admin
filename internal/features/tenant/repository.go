package tenant

import (
	"context"
	"fmt"

	entgo "github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/ent/tenant"
	"github.com/speech/fireworks-admin/internal/pkg/db"
	"github.com/speech/fireworks-admin/internal/pkg/idgen"
)

// TenantRepo 租户数据持久化操作的具体实现
type TenantRepo struct {
	tx *db.TxManager
}

// NewTenantRepo 创建 Repository 实例
func NewTenantRepo(txManager *db.TxManager) *TenantRepo {
	return &TenantRepo{
		tx: txManager,
	}
}

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

func (r *TenantRepo) Delete(ctx context.Context, id string) error {
	tenantId, err := idgen.Parse(id)
	if err != nil {
		return fmt.Errorf("repo:Delete id parse id=%s: %w", id, err)
	}

	err = r.tx.DB(ctx).Tenant.DeleteOneID(tenantId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("repo:Delete id=%s: %w", id, err)
	}
	return nil
}

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
