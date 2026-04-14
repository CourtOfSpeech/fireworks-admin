package tenant

import (
	"context"

	"github.com/speech/fireworks-admin/internal/pkg/api"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// TenantService 封装租户业务逻辑操作
type TenantService struct {
	repo *TenantRepo
}

// NewTenantService 创建 Service 实例
func NewTenantService(repo *TenantRepo) *TenantService {
	return &TenantService{
		repo: repo,
	}
}

func (s *TenantService) Create(ctx context.Context, req *CreateTenantReq) (*Tenant, error) {
	if req.Status == 0 {
		req.Status = TenantStatusEnabled
	}
	return s.repo.Create(ctx, req)
}

func (s *TenantService) Update(ctx context.Context, id string, req *UpdateTenantReq) (*Tenant, error) {
	if req.Status != nil && !IsValidStatus(*req.Status) {
		return nil, ErrInvalidStatus()
	}

	return s.repo.Update(ctx, id, req)
}

func (s *TenantService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *TenantService) GetByID(ctx context.Context, id string) (*Tenant, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TenantService) List(ctx context.Context, query *TenantQuery) (*api.PageResult[*Tenant], error) {
	list, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, bizerr.Internal(err)
	}

	return api.NewPageResult(list, total, query.Page, query.PageSize), nil
}
