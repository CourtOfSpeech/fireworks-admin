// Package user 提供User功能，包括User的创建、查询、更新和删除操作。
// 本文件定义了User模块的数据传输对象（DTO），用于 API 请求和响应的数据结构定义。
package user

import (
	"time"

	"github.com/speech/fireworks-admin/internal/pkg/api"
)

// CreateUserReq 创建User请求结构体。
// 包含创建新User所需的所有必填字段信息。
type CreateUserReq struct {
	Username string `json:"username" validate:"required"`           // 用户名，必填
	Email    string `json:"email" validate:"required"`              // 邮箱，必填
	Phone    string `json:"phone" validate:"required"`              // 手机号，必填
	Password string `json:"password" validate:"required"`           // 密码，必填
	Status   int8   `json:"status" validate:"required,min=1,max=2"` // 状态，必填
}

// UpdateUserReq 更新User请求结构体。
// 所有字段均为指针类型，支持部分更新，仅更新非空字段。
type UpdateUserReq struct {
	Username  *string    `json:"username" validate:"omitempty"`           // 用户名，可选
	Email     *string    `json:"email" validate:"omitempty"`              // 邮箱，可选
	Phone     *string    `json:"phone" validate:"omitempty"`              // 手机号，可选
	Password  *string    `json:"password" validate:"omitempty"`           // 密码，可选
	Nickname  *string    `json:"nickname" validate:"omitempty"`           // 昵称，可选
	Avatar    *string    `json:"avatar" validate:"omitempty"`             // 头像URL，可选
	Status    *int8      `json:"status" validate:"omitempty,min=1,max=2"` // 状态，可选
	UpdatedAt *time.Time `json:"updated_at" validate:"omitempty"`         // 更新时间，可选
	DeletedAt *time.Time `json:"deleted_at" validate:"omitempty"`         // 删除时间，可选
}

// UserQuery User查询条件结构体。
// 用于构建User列表查询的过滤条件和分页参数。
type UserQuery struct {
	api.PageQuery        // 分页查询基础字段
	Username      string `query:"username"`  // 用户名：精确匹配
	Email         string `query:"email"`     // 邮箱：精确匹配
	Phone         string `query:"phone"`     // 手机号：精确匹配
	Password      string `query:"password"`  // 密码：精确匹配
	Nickname      string `query:"nickname"`  // 昵称：精确匹配
	Avatar        string `query:"avatar"`    // 头像URL：精确匹配
	TenantID      string `query:"tenant_id"` // 租户ID：精确匹配
	Status        *int8  `query:"status"`    // 状态：精确匹配
}

// HasUsername 检查是否设置了用户名查询条件。
// 返回 true 表示需要按用户名进行精确查询。
func (q *UserQuery) HasUsername() bool {
	return q.Username != ""
}

// HasEmail 检查是否设置了邮箱查询条件。
// 返回 true 表示需要按邮箱进行精确查询。
func (q *UserQuery) HasEmail() bool {
	return q.Email != ""
}

// HasPhone 检查是否设置了手机号查询条件。
// 返回 true 表示需要按手机号进行精确查询。
func (q *UserQuery) HasPhone() bool {
	return q.Phone != ""
}

// HasPassword 检查是否设置了密码查询条件。
// 返回 true 表示需要按密码进行精确查询。
func (q *UserQuery) HasPassword() bool {
	return q.Password != ""
}

// HasNickname 检查是否设置了昵称查询条件。
// 返回 true 表示需要按昵称进行精确查询。
func (q *UserQuery) HasNickname() bool {
	return q.Nickname != ""
}

// HasAvatar 检查是否设置了头像URL查询条件。
// 返回 true 表示需要按头像URL进行精确查询。
func (q *UserQuery) HasAvatar() bool {
	return q.Avatar != ""
}

// HasTenantID 检查是否设置了租户ID查询条件。
// 返回 true 表示需要按租户ID进行精确查询。
func (q *UserQuery) HasTenantID() bool {
	return q.TenantID != ""
}

// HasStatus 检查是否设置了状态查询条件。
// 返回 true 表示需要按状态进行精确查询。
func (q *UserQuery) HasStatus() bool {
	return q.Status != nil
}
