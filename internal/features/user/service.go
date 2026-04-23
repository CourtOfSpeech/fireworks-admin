// Package user 提供User功能，包括User的创建、查询、更新和删除操作。
// 本文件定义了User模块的业务逻辑层，封装User相关的业务规则和操作。
package user

import (
	"context"

	"github.com/speech/fireworks-admin/internal/pkg/api"
	"github.com/speech/fireworks-admin/internal/pkg/crypto"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// UserService 封装User业务逻辑操作。
// 负责协调 Repository 层完成User的增删改查，并处理业务规则验证。
type UserService struct {
	repo *UserRepo // User数据持久化操作
}

// NewUserService 创建User Service 实例。
// 参数 repo 为User Repository，返回初始化后的 Service 实例。
func NewUserService(repo *UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Create 创建新User。
// 参数 ctx 为上下文，req 为创建请求参数。
// 密码字段会被加密存储。
// 返回创建成功的User实体和可能的错误。
func (s *UserService) Create(ctx context.Context, req *CreateUserReq) (*User, error) {
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, bizerr.Internal(err)
	}
	req.Password = hashedPassword
	t, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, wrapError(err)
	}
	return t, nil
}

// Update 更新User信息。
// 参数 ctx 为上下文，id 为User ID，req 为更新请求参数。
// 如果更新密码字段，密码会被加密存储。
// 返回更新后的User实体和可能的错误。
func (s *UserService) Update(ctx context.Context, id string, req *UpdateUserReq) (*User, error) {
	if req.Password != nil {
		hashedPassword, err := crypto.HashPassword(*req.Password)
		if err != nil {
			return nil, bizerr.Internal(err)
		}
		*req.Password = hashedPassword
	}

	t, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, wrapError(err)
	}
	return t, nil
}

// Delete 删除User。
// 参数 ctx 为上下文，id 为User ID。
// 返回删除操作可能发生的错误。
func (s *UserService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return wrapError(err)
	}
	return nil
}

// GetByID 根据User ID 获取User详情。
// 参数 ctx 为上下文，id 为User ID。
// 返回User实体和可能的错误。
func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, wrapError(err)
	}
	return t, nil
}

// List 根据查询条件获取User列表。
// 参数 ctx 为上下文，query 为查询条件。
// 返回分页结果和可能的错误。
func (s *UserService) List(ctx context.Context, query *UserQuery) (*api.PageResult[*User], error) {
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
