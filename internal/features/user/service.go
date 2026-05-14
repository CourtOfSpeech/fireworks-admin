// Package user 提供User功能，包括User的创建、查询、更新和删除操作。
// 本文件定义了User模块的业务逻辑层，封装User相关的业务规则和操作。
package user

import (
	"context"
	"regexp"
	"time"

	"github.com/speech/fireworks-admin/internal/features/tenant"
	"github.com/speech/fireworks-admin/internal/middleware"
	"github.com/speech/fireworks-admin/internal/pkg/api"
	"github.com/speech/fireworks-admin/internal/pkg/crypto"
	"github.com/speech/fireworks-admin/internal/pkg/ctxutil"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
	"github.com/speech/fireworks-admin/internal/pkg/idgen"
)

// UserService 封装User业务逻辑操作。
// 负责协调 Repository 层完成User的增删改查，并处理业务规则验证。
type UserService struct {
	repo       *UserRepo          // User数据持久化操作
	tenantRepo *tenant.TenantRepo // 租户数据持久化操作
}

// NewUserService 创建User Service 实例。
// 参数 repo 为User Repository，tenantRepo 为租户 Repository，返回初始化后的 Service 实例。
func NewUserService(repo *UserRepo, tenantRepo *tenant.TenantRepo) *UserService {
	return &UserService{
		repo:       repo,
		tenantRepo: tenantRepo,
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

// emailRegex 邮箱格式正则表达式。
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// phoneRegex 手机号格式正则表达式（中国大陆手机号）。
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// Login 用户登录。
// 参数 ctx 为上下文，req 为登录请求参数。
// 根据 Identity 字段判断是用户名、邮箱还是手机号进行登录。
// 验证租户标识、用户状态和密码，成功后返回 JWT token 和用户信息。
// 返回登录响应和可能的错误。
func (s *UserService) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	var tenantID string
	var err error

	if req.TenantID != nil && *req.TenantID != "" {
		tenantID = *req.TenantID
	} else if req.TenantName != nil && *req.TenantName != "" {
		tenant, err := s.tenantRepo.FindByName(ctx, *req.TenantName)
		if err != nil {
			return nil, ErrTenantMismatch()
		}
		tenantID = tenant.ID
	} else {
		return nil, bizerr.InvalidParam("租户标识不能为空")
	}

	tenantUUID, err := idgen.Parse(tenantID)
	if err != nil {
		return nil, bizerr.InvalidParamWrap(err, "无效的租户ID")
	}

	loginCtx := ctxutil.WithTenant(ctx, tenantUUID)

	var user *User
	switch {
	case isEmail(req.Identity):
		user, err = s.repo.FindByEmail(loginCtx, req.Identity)
	case isPhone(req.Identity):
		user, err = s.repo.FindByPhone(loginCtx, req.Identity)
	default:
		user, err = s.repo.FindByUsername(loginCtx, req.Identity)
	}

	if err != nil {
		return nil, ErrLoginFailed()
	}

	if user.TenantID != tenantID {
		return nil, ErrTenantMismatch()
	}

	if user.Status == 2 {
		return nil, ErrUserDisabled()
	}

	if !crypto.CheckPassword(req.Password, user.Password) {
		return nil, ErrLoginFailed()
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.TenantID)
	if err != nil {
		return nil, bizerr.Internal(err)
	}

	return &LoginResp{
		Token:     token,
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		TenantID:  user.TenantID,
		ExpiresAt: time.Now().Add(middleware.GetExpireDuration()).Unix(),
	}, nil
}

// RefreshToken 刷新 JWT token。
// 参数 ctx 为上下文，userID 为用户 ID。
// 验证用户状态后生成新的 JWT token。
// 返回刷新 token 响应和可能的错误。
func (s *UserService) RefreshToken(ctx context.Context, userID string) (*RefreshTokenResp, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, wrapError(err)
	}

	if user.Status == 2 {
		return nil, ErrUserDisabled()
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.TenantID)
	if err != nil {
		return nil, bizerr.Internal(err)
	}

	return &RefreshTokenResp{
		Token:     token,
		ExpiresAt: time.Now().Add(middleware.GetExpireDuration()).Unix(),
	}, nil
}

// GetCurrentUserInfo 获取当前用户信息。
// 参数 ctx 为上下文，userID 为用户 ID。
// 返回当前用户信息响应和可能的错误。
func (s *UserService) GetCurrentUserInfo(ctx context.Context, userID string) (*CurrentUserResp, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, wrapError(err)
	}

	return &CurrentUserResp{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		TenantID: user.TenantID,
		Status:   user.Status,
	}, nil
}

// Logout 用户登出。
// 参数 ctx 为上下文。
// 由于 JWT 是无状态的，不实现 token 黑名单，直接返回 nil。
// 返回可能的错误。
func (s *UserService) Logout(ctx context.Context) error {
	return nil
}

// isEmail 判断字符串是否为邮箱格式。
// 参数 s 为待判断的字符串。
// 返回 true 表示是邮箱格式，false 表示不是。
func isEmail(s string) bool {
	return emailRegex.MatchString(s)
}

// isPhone 判断字符串是否为手机号格式。
// 参数 s 为待判断的字符串。
// 返回 true 表示是手机号格式，false 表示不是。
func isPhone(s string) bool {
	return phoneRegex.MatchString(s)
}
