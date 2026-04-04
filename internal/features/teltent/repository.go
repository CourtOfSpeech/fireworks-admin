package teltent

import (
	"context"

	entgo "github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/ent/teltent"
	"github.com/speech/fireworks-admin/internal/pkg/idgen"
)

// Repository 租户数据持久化操作。
type Repository struct {
	client *entgo.Client
}

// NewRepository 使用给定的 Ent 客户端创建 Repository 实例。
func NewRepository(client *entgo.Client) *Repository {
	return &Repository{
		client: client,
	}
}

// toEntity 将 Ent Teltent 模型转换为领域 Teltent 实体。
func toEntity(t *entgo.Teltent) *Teltent {
	return &Teltent{
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
func (r *Repository) Create(ctx context.Context, req *CreateTeltentReq) (*Teltent, error) {
	id, err := idgen.NewV7()
	if err != nil {
		return nil, err
	}

	builder := r.client.Teltent.Create().
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
	telentId, err := idgen.Parse(id)
	if err != nil {
		return err
	}
	return r.client.Teltent.DeleteOneID(telentId).Exec(ctx)
}

// FindByPage 根据查询条件分页查询租户列表。
func (r *Repository) FindByPage(ctx context.Context, query *TeltentQuery) ([]*Teltent, int64, error) {
	builder := r.client.Teltent.Query()

	if query.HasKeyword() {
		builder.Where(
			teltent.Or(
				teltent.CertificateNoHasPrefix(query.Keyword),
				teltent.NameContains(query.Keyword),
			),
		)
	}

	if query.HasStatus() {
		builder.Where(teltent.StatusEQ(*query.Status))
	}

	if query.HasEmail() {
		builder.Where(teltent.EmailEQ(query.Email))
	}

	if query.HasPhone() {
		builder.Where(teltent.PhoneEQ(query.Phone))
	}

	total, err := builder.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	teltents, err := builder.
		Offset(query.GetOffset()).
		Limit(query.GetLimit()).
		Order(entgo.Desc(teltent.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*Teltent, 0, len(teltents))
	for _, t := range teltents {
		result = append(result, toEntity(t))
	}

	return result, int64(total), nil
}

// GetByID 根据ID从数据库查询租户。
func (r *Repository) GetByID(ctx context.Context, id string) (*Teltent, error) {
	telentId, err := idgen.Parse(id)
	if err != nil {
		return nil, err
	}
	t, err := r.client.Teltent.Get(ctx, telentId)
	if err != nil {
		return nil, err
	}
	return toEntity(t), nil
}

// ExistsByCertificateNo 根据证件号检查租户是否已存在。
// 返回 true 表示该证件号已被使用，false 表示未使用。
func (r *Repository) ExistsByCertificateNo(ctx context.Context, certNo string) (bool, error) {
	count, err := r.client.Teltent.Query().
		Where(teltent.CertificateNoEQ(certNo)).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByCertificateNoExcludingID 根据证件号检查除指定ID外的租户是否已存在。
// 用于更新操作时排除自身记录的唯一性校验。
func (r *Repository) ExistsByCertificateNoExcludingID(ctx context.Context, certNo string, excludeID string) (bool, error) {
	telentId, err := idgen.Parse(excludeID)
	if err != nil {
		return false, err
	}
	count, err := r.client.Teltent.Query().
		Where(
			teltent.CertificateNoEQ(certNo),
			teltent.IDNEQ(telentId),
		).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update 根据ID更新数据库中的租户信息。
func (r *Repository) Update(ctx context.Context, id string, req *UpdateTeltentReq) (*Teltent, error) {
	telentId, err := idgen.Parse(id)
	if err != nil {
		return nil, err
	}
	builder := r.client.Teltent.UpdateOneID(telentId)
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
