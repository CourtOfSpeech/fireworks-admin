package teltent

import (
	"context"

	"github.com/speech/fireworks-admin/internal/pkg/api"
)

// Service 封装租户业务逻辑操作。
type Service struct {
	repo *Repository
}

// NewService 创建 Service 实例。
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create 根据请求参数创建新租户。
func (s *Service) Create(ctx context.Context, req *CreateTeltentReq) (*Teltent, error) {
	if req.Status == 0 {
		req.Status = TeltentStatusEnabled
	}
	return s.repo.Create(ctx, req)
}

// Update 根据ID和请求参数更新租户信息。
func (s *Service) Update(ctx context.Context, id string, req *UpdateTeltentReq) (*Teltent, error) {
	return s.repo.Update(ctx, id, req)
}

// Delete 根据ID删除租户。
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetByID 根据ID查询租户。
func (s *Service) GetByID(ctx context.Context, id string) (*Teltent, error) {
	return s.repo.GetByID(ctx, id)
}

// FindByPage 根据查询条件分页查询租户列表。
func (s *Service) FindByPage(ctx context.Context, query *TeltentQuery) (*api.PageResult[*Teltent], error) {
	list, total, err := s.repo.FindByPage(ctx, query)
	if err != nil {
		return nil, err
	}

	return api.NewPageResult(list, total, query.Page, query.PageSize), nil
}
