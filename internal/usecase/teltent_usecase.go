package usecase

import (
	"context"

	"github.com/speech/fireworks-admin/internal/domain/entity"
	"github.com/speech/fireworks-admin/internal/domain/repo"
	"github.com/speech/fireworks-admin/pkg/response"
)

// TeltentUsecase 封装租户业务逻辑操作。
type TeltentUsecase struct {
	repo repo.TeltentRepo
}

// NewTeltentUsecase 创建 TeltentUsecase 实例。
func NewTeltentUsecase(repo repo.TeltentRepo) *TeltentUsecase {
	return &TeltentUsecase{
		repo: repo,
	}
}

// Create 根据请求参数创建新租户。
func (u *TeltentUsecase) Create(ctx context.Context, req *entity.CreateTeltentReq) (*entity.Teltent, error) {
	if req.Status == 0 {
		req.Status = entity.TeltentStatusEnabled
	}
	return u.repo.Create(ctx, req)
}

// Update 根据ID和请求参数更新租户信息。
func (u *TeltentUsecase) Update(ctx context.Context, id string, req *entity.UpdateTeltentReq) (*entity.Teltent, error) {
	return u.repo.Update(ctx, id, req)
}

// Delete 根据ID删除租户。
func (u *TeltentUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

// GetByID 根据ID查询租户。
func (u *TeltentUsecase) GetByID(ctx context.Context, id string) (*entity.Teltent, error) {
	return u.repo.GetByID(ctx, id)
}

// FindByPage 根据查询条件分页查询租户列表。
func (u *TeltentUsecase) FindByPage(ctx context.Context, query *entity.TeltentQuery) (*response.PageResult[*entity.Teltent], error) {
	list, total, err := u.repo.FindByPage(ctx, query)
	if err != nil {
		return nil, err
	}

	return response.NewPageResult(list, total, query.Page, query.PageSize), nil
}
