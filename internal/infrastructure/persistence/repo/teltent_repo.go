package repo

import (
	"context"

	"github.com/speech/fireworks-admin/internal/domain/entity"
	domainrepo "github.com/speech/fireworks-admin/internal/domain/repo"
	"github.com/speech/fireworks-admin/internal/infrastructure/persistence/ent"
	"github.com/speech/fireworks-admin/internal/infrastructure/persistence/ent/teltent"
	"github.com/speech/fireworks-admin/pkg/utils"
)

// teltentEnt 使用 Ent ORM 实现 TeltentRepo 接口。
type teltentEnt struct {
	entClient *ent.Client
}

// NewTeltentEnt 使用给定的 Ent 客户端创建 TeltentRepo 实例。
func NewTeltentEnt(client *ent.Client) domainrepo.TeltentRepo {
	return &teltentEnt{
		entClient: client,
	}
}

// toEntity 将 Ent Teltent 模型转换为领域 Teltent 实体。
func toEntity(t *ent.Teltent) *entity.Teltent {
	return &entity.Teltent{
		ID:            utils.ToString(t.ID),
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
func (t *teltentEnt) Create(ctx context.Context, req *entity.CreateTeltentReq) (*entity.Teltent, error) {
	id, err := utils.NewV7()
	if err != nil {
		return nil, err
	}

	builder := t.entClient.Teltent.Create().
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

	teltent, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	return toEntity(teltent), nil
}

// Delete 根据ID从数据库中删除租户。
func (t *teltentEnt) Delete(ctx context.Context, id string) error {
	telentId, err := utils.Parse(id)
	if err != nil {
		return err
	}
	return t.entClient.Teltent.DeleteOneID(telentId).Exec(ctx)
}

// GetByPage 根据查询条件分页查询租户列表。
func (t *teltentEnt) FindByPage(ctx context.Context, query *entity.TeltentQuery) ([]*entity.Teltent, int64, error) {
	builder := t.entClient.Teltent.Query()

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
		Order(ent.Desc(teltent.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.Teltent, 0, len(teltents))
	for _, t := range teltents {
		result = append(result, toEntity(t))
	}

	return result, int64(total), nil
}

// GetByID 根据ID从数据库查询租户。
func (t *teltentEnt) GetByID(ctx context.Context, id string) (*entity.Teltent, error) {
	telentId, err := utils.Parse(id)
	if err != nil {
		return nil, err
	}
	teltent, err := t.entClient.Teltent.Get(ctx, telentId)
	if err != nil {
		return nil, err
	}
	return toEntity(teltent), nil
}

// Update 根据ID更新数据库中的租户信息。
func (t *teltentEnt) Update(ctx context.Context, id string, req *entity.UpdateTeltentReq) (*entity.Teltent, error) {
	telentId, err := utils.Parse(id)
	if err != nil {
		return nil, err
	}
	builder := t.entClient.Teltent.UpdateOneID(telentId)
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

	teltent, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toEntity(teltent), nil
}
