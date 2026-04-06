package tenant

import (
	"context"

	entgo "github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/ent/tenant"
	"github.com/speech/fireworks-admin/internal/pkg/db"
	"github.com/speech/fireworks-admin/internal/pkg/idgen"
)

// Repository 租户数据持久化操作。
type Repository struct {
	tx *db.TxManager
}

// NewRepository 使用给定的 Ent 客户端创建 Repository 实例。
func NewRepository(txManager *db.TxManager) *Repository {
	return &Repository{
		tx: txManager,
	}
}

// toEntity 将 Ent Tenant 模型转换为领域 Tenant 实体。
func toEntity(t *entgo.Tenant) *Tenant {
	return &Tenant{
		ID:            idgen.ToString(t.ID),
		CertificateNo: t.CertificateNo,
		Name:          t.Name,
		Type:          t.Type,
		ContactName:   t.ContactName,
		Email:         t.Email,
		Phone:         t.Phone,
		Status:        t.Status,
		ExpiredAt:     t.ExpiredAt,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}

// Create 在数据库中创建新租户。
func (r *Repository) Create(ctx context.Context, req *CreateTenantReq) (*Tenant, error) {
	id, err := idgen.NewV7()
	if err != nil {
		return nil, err
	}

	builder := r.tx.DB(ctx).Tenant.Create().
		SetID(id).
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
		return nil, err
	}

	return toEntity(t), nil
}

// Delete 根据ID从数据库中删除租户。
func (r *Repository) Delete(ctx context.Context, id string) error {
	tenantId, err := idgen.Parse(id)
	if err != nil {
		return err
	}
	return r.tx.DB(ctx).Tenant.DeleteOneID(tenantId).Exec(ctx)
}

// FindByPage 根据查询条件分页查询租户列表。
func (r *Repository) FindByPage(ctx context.Context, query *TenantQuery) ([]*Tenant, int64, error) {
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
		return nil, 0, err
	}

	tenants, err := builder.
		Offset(query.GetOffset()).
		Limit(query.GetLimit()).
		Order(entgo.Desc(tenant.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*Tenant, 0, len(tenants))
	for _, t := range tenants {
		result = append(result, toEntity(t))
	}

	return result, int64(total), nil
}

// GetByID 根据ID从数据库查询租户。
func (r *Repository) GetByID(ctx context.Context, id string) (*Tenant, error) {
	tenantId, err := idgen.Parse(id)
	if err != nil {
		return nil, err
	}
	t, err := r.tx.DB(ctx).Tenant.Get(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	return toEntity(t), nil
}

// ExistsByCertificateNo 根据证件号检查租户是否已存在。
// 返回 true 表示该证件号已被使用，false 表示未使用。
func (r *Repository) ExistsByCertificateNo(ctx context.Context, certNo string) (bool, error) {
	count, err := r.tx.DB(ctx).Tenant.Query().
		Where(tenant.CertificateNoEQ(certNo)).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByCertificateNoExcludingID 根据证件号检查除指定ID外的租户是否已存在。
// 用于更新操作时排除自身记录的唯一性校验。
func (r *Repository) ExistsByCertificateNoExcludingID(ctx context.Context, certNo string, excludeID string) (bool, error) {
	tenantId, err := idgen.Parse(excludeID)
	if err != nil {
		return false, err
	}
	count, err := r.tx.DB(ctx).Tenant.Query().
		Where(
			tenant.CertificateNoEQ(certNo),
			tenant.IDNEQ(tenantId),
		).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update 根据ID更新数据库中的租户信息。
func (r *Repository) Update(ctx context.Context, id string, req *UpdateTenantReq) (*Tenant, error) {
	tenantId, err := idgen.Parse(id)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return toEntity(t), nil
}
