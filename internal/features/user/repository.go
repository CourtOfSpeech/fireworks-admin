// Package user 提供User功能，包括User的创建、查询、更新和删除操作。
// 本文件定义了User模块的数据持久化层，负责与数据库交互。
package user

import (
	"context"
	"fmt"

	entgo "github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/ent/user"
	"github.com/speech/fireworks-admin/internal/pkg/db"
	"github.com/speech/fireworks-admin/internal/pkg/idgen"
)

// UserRepo User数据持久化操作的具体实现。
// 封装了User相关的数据库操作，包括增删改查。
type UserRepo struct {
	tx *db.TxManager // 数据库事务管理器
}

// NewUserRepo 创建User Repository 实例。
// 参数 txManager 为数据库事务管理器，返回初始化后的 Repository 实例。
func NewUserRepo(txManager *db.TxManager) *UserRepo {
	return &UserRepo{
		tx: txManager,
	}
}

// toEntity 将 Ent 框架的 User 模型转换为领域模型 User。
// 参数 t 为 Ent 框架的User模型，返回领域模型的User实体。
func toEntity(t *entgo.User) *User {
	return &User{
		Username:  t.Username,
		Email:     t.Email,
		Phone:     t.Phone,
		Password:  t.Password,
		Nickname:  t.Nickname,
		Avatar:    t.Avatar,
		ID:        idgen.ToString(t.ID),
		TenantID:  idgen.ToString(t.TenantID),
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: t.DeletedAt,
	}
}

// Create 创建新User记录。
// 参数 ctx 为上下文，req 为创建请求参数。
// 返回创建成功的User实体和可能的错误。
func (r *UserRepo) Create(ctx context.Context, req *CreateUserReq) (*User, error) {
	builder := r.tx.DB(ctx).User.Create()
	builder.SetUsername(req.Username)
	builder.SetEmail(req.Email)
	builder.SetPhone(req.Phone)
	builder.SetPassword(req.Password)
	builder.SetStatus(req.Status)

	t, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo:Create: %w", err)
	}

	return toEntity(t), nil
}

// Delete 根据User ID 删除User记录。
// 参数 ctx 为上下文，id 为User ID 字符串。
// 返回删除操作可能发生的错误。
func (r *UserRepo) Delete(ctx context.Context, id string) error {
	userId, err := idgen.Parse(id)
	if err != nil {
		return fmt.Errorf("repo:Delete id parse id=%s: %w", id, err)
	}
	err = r.tx.DB(ctx).User.DeleteOneID(userId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("repo:Delete id=%s: %w", id, err)
	}
	return nil
}

// List 根据查询条件获取User列表。
// 参数 ctx 为上下文，query 为查询条件。
// 返回User列表、总数和可能的错误。支持分页和条件过滤。
func (r *UserRepo) List(ctx context.Context, query *UserQuery) ([]*User, int64, error) {
	builder := r.tx.DB(ctx).User.Query()
	if query.HasUsername() {
		builder.Where(user.UsernameEQ(query.Username))
	}
	if query.HasEmail() {
		builder.Where(user.EmailEQ(query.Email))
	}
	if query.HasPhone() {
		builder.Where(user.PhoneEQ(query.Phone))
	}
	if query.HasNickname() {
		builder.Where(user.NicknameEQ(query.Nickname))
	}
	if query.HasAvatar() {
		builder.Where(user.AvatarEQ(query.Avatar))
	}
	if query.HasTenantID() {
		tenantid, err := idgen.Parse(query.TenantID)
		if err != nil {
			return nil, 0, fmt.Errorf("repo:List parse tenantid=%s: %w", query.TenantID, err)
		}
		builder.Where(user.TenantIDEQ(tenantid))
	}
	if query.HasStatus() {
		builder.Where(user.StatusEQ(*query.Status))
	}

	total, err := builder.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo:List count: %w", err)
	}

	items, err := builder.
		Offset(query.GetOffset()).
		Limit(query.GetLimit()).
		Order(entgo.Desc(user.FieldID)).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo:List query: %w", err)
	}

	result := make([]*User, 0, len(items))
	for _, t := range items {
		result = append(result, toEntity(t))
	}

	return result, int64(total), nil
}

// GetByID 根据User ID 获取User详情。
// 参数 ctx 为上下文，id 为User ID 字符串。
// 返回User实体和可能的错误。
func (r *UserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	userId, err := idgen.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("repo:GetByID id parse id=%s: %w", id, err)
	}

	t, err := r.tx.DB(ctx).User.Get(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("repo:GetByID id=%s: %w", id, err)
	}
	return toEntity(t), nil
}

// Update 根据User ID 更新User信息。
// 参数 ctx 为上下文，id 为User ID 字符串，req 为更新请求参数。
// 仅更新请求中非空字段。返回更新后的User实体和可能的错误。
func (r *UserRepo) Update(ctx context.Context, id string, req *UpdateUserReq) (*User, error) {
	userId, err := idgen.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("repo:Update id parse id=%s: %w", id, err)
	}

	builder := r.tx.DB(ctx).User.UpdateOneID(userId)
	if req.Username != nil {
		builder.SetUsername(*req.Username)
	}
	if req.Email != nil {
		builder.SetEmail(*req.Email)
	}
	if req.Phone != nil {
		builder.SetPhone(*req.Phone)
	}
	if req.Password != nil {
		builder.SetPassword(*req.Password)
	}
	if req.Nickname != nil {
		builder.SetNickname(*req.Nickname)
	}
	if req.Avatar != nil {
		builder.SetAvatar(*req.Avatar)
	}
	if req.Status != nil {
		builder.SetStatus(*req.Status)
	}

	t, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo:Update id=%s: %w", id, err)
	}
	return toEntity(t), nil
}
