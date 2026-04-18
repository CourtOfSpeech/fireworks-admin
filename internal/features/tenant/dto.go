// Package tenant 提供租户管理功能，包括租户的创建、查询、更新和删除操作。
// 本文件定义了租户模块的数据传输对象（DTO），用于 API 请求和响应的数据结构定义。
package tenant

import (
	"time"

	"github.com/speech/fireworks-admin/internal/pkg/api"
)

// CreateTenantReq 创建租户请求结构体。
// 包含创建新租户所需的所有必填字段信息。
type CreateTenantReq struct {
	CertificateNo string    `json:"certificate_no" validate:"required"` // 证件号码，必填
	Name          string    `json:"name" validate:"required"`           // 租户名称，必填
	Type          int8      `json:"type" validate:"required,min=1,max=2"` // 租户类型：1-企业，2-个人，必填
	ContactName   string    `json:"contact_name" validate:"required"`   // 联系人姓名，必填
	Email         string    `json:"email" validate:"required,email"`    // 联系邮箱，必填，需符合邮箱格式
	Phone         string    `json:"phone" validate:"required"`          // 联系电话，必填
	ExpiredAt     time.Time `json:"expired_at" validate:"required"`     // 过期时间，必填
	Status        int8      `json:"status" validate:"required,min=1,max=2"` // 状态：1-禁用，2-正常，必填
}

// UpdateTenantReq 更新租户请求结构体。
// 所有字段均为指针类型，支持部分更新，仅更新非空字段。
type UpdateTenantReq struct {
	CertificateNo *string    `json:"certificate_no" validate:"omitempty"` // 证件号码，可选
	Name          *string    `json:"name" validate:"omitempty"`           // 租户名称，可选
	Type          *int8      `json:"type" validate:"omitempty,min=1,max=2"` // 租户类型：1-企业，2-个人，可选
	ContactName   *string    `json:"contact_name" validate:"omitempty"`   // 联系人姓名，可选
	Email         *string    `json:"email" validate:"omitempty,email"`    // 联系邮箱，可选，需符合邮箱格式
	Phone         *string    `json:"phone" validate:"omitempty"`          // 联系电话，可选
	ExpiredAt     *time.Time `json:"expired_at" validate:"omitempty"`     // 过期时间，可选
	Status        *int8      `json:"status" validate:"omitempty,min=1,max=2"` // 状态：1-禁用，2-正常，可选
}

// TenantQuery 租户查询条件结构体。
// 用于构建租户列表查询的过滤条件和分页参数。
type TenantQuery struct {
	api.PageQuery                                                 // 分页查询基础字段
	Keyword string `query:"keyword"` // 关键字：模糊匹配证件号码（前缀）、名称
	Status  *int8  `query:"status"`  // 状态：精确匹配
	Email   string `query:"email"`   // 邮箱：精确匹配
	Phone   string `query:"phone"`   // 电话：精确匹配
}

// HasKeyword 检查是否设置了关键字查询条件。
// 返回 true 表示需要按关键字进行模糊查询。
func (q *TenantQuery) HasKeyword() bool {
	return q.Keyword != ""
}

// HasStatus 检查是否设置了状态查询条件。
// 返回 true 表示需要按状态进行精确查询。
func (q *TenantQuery) HasStatus() bool {
	return q.Status != nil
}

// HasEmail 检查是否设置了邮箱查询条件。
// 返回 true 表示需要按邮箱进行精确查询。
func (q *TenantQuery) HasEmail() bool {
	return q.Email != ""
}

// HasPhone 检查是否设置了电话查询条件。
// 返回 true 表示需要按电话进行精确查询。
func (q *TenantQuery) HasPhone() bool {
	return q.Phone != ""
}
