package teltent

import (
	"context"

	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
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
// 执行步骤：1) 证件号唯一性校验 2) 设置默认状态 3) 持久化存储。
func (s *Service) Create(ctx context.Context, req *CreateTeltentReq) (*Teltent, error) {
	exists, err := s.repo.ExistsByCertificateNo(ctx, req.CertificateNo)
	if err != nil {
		return nil, bizerr.Wrap(err, "检查证件号唯一性失败")
	}
	if exists {
		return nil, ErrDuplicateCertNo
	}

	if req.Status == 0 {
		req.Status = TeltentStatusEnabled
	}
	return s.repo.Create(ctx, req)
}

// Update 根据ID和请求参数更新租户信息。
// 执行步骤：1) 存在性校验 2) 证件号唯一性校验（排除自身）3) 状态有效性校验 4) 持久化更新。
func (s *Service) Update(ctx context.Context, id string, req *UpdateTeltentReq) (*Teltent, error) {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, NewTeltentNotFound(id)
	}

	if req.CertificateNo != nil {
		exists, err := s.repo.ExistsByCertificateNoExcludingID(ctx, *req.CertificateNo, id)
		if err != nil {
			return nil, bizerr.Wrap(err, "检查证件号唯一性失败")
		}
		if exists {
			return nil, ErrDuplicateCertNo
		}
	}

	if req.Status != nil && !isValidStatus(*req.Status) {
		return nil, ErrInvalidStatus
	}

	return s.repo.Update(ctx, id, req)
}

// Delete 根据ID删除租户。
// 先执行存在性校验，确保目标记录存在后再执行删除操作，
// 从而为调用方提供更精确的错误信息（ NotFound vs InternalError）。
func (s *Service) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return NewTeltentNotFound(id)
	}
	return s.repo.Delete(ctx, id)
}

// GetByID 根据ID查询租户。
// 当记录不存在时返回具体的 NotFoundError，而非泛型数据库错误。
func (s *Service) GetByID(ctx context.Context, id string) (*Teltent, error) {
	teltent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, NewTeltentNotFound(id)
	}
	return teltent, nil
}

// FindByPage 根据查询条件分页查询租户列表。
func (s *Service) FindByPage(ctx context.Context, query *TeltentQuery) (*api.PageResult[*Teltent], error) {
	list, total, err := s.repo.FindByPage(ctx, query)
	if err != nil {
		return nil, err
	}

	return api.NewPageResult(list, total, query.Page, query.PageSize), nil
}

// isValidStatus 校验状态值是否在允许的范围内。
// 当前仅支持 TeltentStatusDisabled(1) 和 TeltentStatusEnabled(2)。
func isValidStatus(status int8) bool {
	return status == TeltentStatusDisabled || status == TeltentStatusEnabled
}
