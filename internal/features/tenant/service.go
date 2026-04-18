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
	t, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, wrapError(err)
	}
	return t, nil
}

func (s *TenantService) Update(ctx context.Context, id string, req *UpdateTenantReq) (*Tenant, error) {
	if req.Status != nil && !IsValidStatus(*req.Status) {
		return nil, ErrInvalidStatus()
	}

	t, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, wrapError(err)
	}
	return t, nil
}

func (s *TenantService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return wrapError(err)
	}
	return nil
}

func (s *TenantService) GetByID(ctx context.Context, id string) (*Tenant, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, wrapError(err)
	}
	return t, nil
}

func (s *TenantService) List(ctx context.Context, query *TenantQuery) (*api.PageResult[*Tenant], error) {
	list, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, wrapError(err)
	}

	return api.NewPageResult(list, total, query.Page, query.PageSize), nil
}

// wrapError 解析 Repository 层错误并包装为业务错误。
// 如果错误已经是 BizError 则直接返回，否则包装为内部错误。
func wrapError(err error) error {
	if err == nil {
		return nil
	}
	parsed := ParseRepoError(err)
	if _, ok := parsed.(*bizerr.BizError); ok {
		return parsed
	}
	return bizerr.Internal(parsed)
}
