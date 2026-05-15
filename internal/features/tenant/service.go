// Package tenant 提供租户管理功能，包括租户的创建、查询、更新和删除操作。
// 本文件定义了租户模块的业务逻辑层，封装租户相关的业务规则和操作。
package tenant

import (
	"context"

	"github.com/speech/fireworks-admin/internal/pkg/api"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// TenantService 封装租户业务逻辑操作。
// 负责协调 Repository 层完成租户的增删改查，并处理业务规则验证。
type TenantService struct {
	repo *TenantRepo // 租户数据持久化操作
}

// NewTenantService 创建租户 Service 实例。
// 参数 repo 为租户 Repository，返回初始化后的 Service 实例。
func NewTenantService(repo *TenantRepo) *TenantService {
	return &TenantService{
		repo: repo,
	}
}

// Create 创建新租户。
// 参数 ctx 为上下文，req 为创建请求参数。
// 如果未指定状态，默认设置为正常状态。返回创建成功的租户实体和可能的错误。
func (s *TenantService) Create(ctx context.Context, req *CreateTenantReq) (*Tenant, error) {
	if req.Status == 0 {
		req.Status = TenantStatusEnabled
	}
	t, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, bizerr.WrapRepoError(err, repoParser)
	}
	return t, nil
}

// Update 更新租户信息。
// 参数 ctx 为上下文，id 为租户 ID，req 为更新请求参数。
// 验证状态值有效性后执行更新。返回更新后的租户实体和可能的错误。
func (s *TenantService) Update(ctx context.Context, id string, req *UpdateTenantReq) (*Tenant, error) {
	if req.Status != nil && !IsValidStatus(*req.Status) {
		return nil, ErrInvalidStatus()
	}

	t, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, bizerr.WrapRepoError(err, repoParser)
	}
	return t, nil
}

// Delete 删除租户。
// 参数 ctx 为上下文，id 为租户 ID。
// 返回删除操作可能发生的错误。
func (s *TenantService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return bizerr.WrapRepoError(err, repoParser)
	}
	return nil
}

// GetByID 根据租户 ID 获取租户详情。
// 参数 ctx 为上下文，id 为租户 ID。
// 返回租户实体和可能的错误。
func (s *TenantService) GetByID(ctx context.Context, id string) (*Tenant, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, bizerr.WrapRepoError(err, repoParser)
	}
	return t, nil
}

// List 根据查询条件获取租户列表。
// 参数 ctx 为上下文，query 为查询条件。
// 返回分页结果和可能的错误。
func (s *TenantService) List(ctx context.Context, query *TenantQuery) (*api.PageResult[*Tenant], error) {
	list, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, bizerr.WrapRepoError(err, repoParser)
	}

	return api.NewPageResult(list, total, query.Page, query.PageSize), nil
}

// FindByName 根据租户名称查询租户。
// 参数 ctx 为上下文，name 为租户名称。
// 返回租户实体和可能的错误。
func (s *TenantService) FindByName(ctx context.Context, name string) (*Tenant, error) {
	return s.repo.FindByName(ctx, name)
}
