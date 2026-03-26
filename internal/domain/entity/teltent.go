package entity

import (
	"time"

	"github.com/speech/fireworks-admin/pkg/response"
)

// 租户类型
const (
	TeltentTypeCompany  int8 = 1 // 企业
	TeltentTypePersonal int8 = 2 // 个人
)

// 租户状态
const (
	TeltentStatusDisabled int8 = 1 // 禁用
	TeltentStatusEnabled  int8 = 2 // 正常
)

// Teltent 租户实体
type Teltent struct {
	ID            string    `json:"id"`
	CertificateNo string    `json:"certificate_no"`
	Name          string    `json:"name"`
	Type          int8      `json:"type"`
	ContactName   string    `json:"contact_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Status        int8      `json:"status"`
	ExpiredAt     time.Time `json:"expired_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}

// CreateTeltentReq 创建租户请求
type CreateTeltentReq struct {
	CertificateNo string    `json:"certificate_no" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	Type          int8      `json:"type" validate:"required,min=1,max=2"`
	ContactName   string    `json:"contact_name" validate:"required"`
	Email         string    `json:"email" validate:"required,email"`
	Phone         string    `json:"phone" validate:"required"`
	ExpiredAt     time.Time `json:"expired_at" validate:"required"`
	Status        int8      `json:"status" validate:"required,min=1,max=2"`
}

// UpdateTeltentReq 更新租户请求
type UpdateTeltentReq struct {
	CertificateNo *string    `json:"certificate_no" validate:"omitempty"`
	Name          *string    `json:"name" validate:"omitempty"`
	Type          *int8      `json:"type" validate:"omitempty,min=1,max=2"`
	ContactName   *string    `json:"contact_name" validate:"omitempty"`
	Email         *string    `json:"email" validate:"omitempty,email"`
	Phone         *string    `json:"phone" validate:"omitempty"`
	ExpiredAt     *time.Time `json:"expired_at" validate:"omitempty"`
	Status        *int8      `json:"status" validate:"omitempty,min=0,max=1"`
}

// TeltentQuery 租户查询条件
type TeltentQuery struct {
	response.PageQuery
	Keyword string `query:"keyword"` // 关键字：模糊匹配证件号码（前缀）、名称
	Status  *int8  `query:"status"`  // 状态：精确匹配
	Email   string `query:"email"`   // 邮箱：精确匹配
	Phone   string `query:"phone"`   // 电话：精确匹配
}

// HasKeyword 是否有关键字查询
func (q *TeltentQuery) HasKeyword() bool {
	return q.Keyword != ""
}

// HasStatus 是否有状态查询
func (q *TeltentQuery) HasStatus() bool {
	return q.Status != nil
}

// HasEmail 是否有邮箱查询
func (q *TeltentQuery) HasEmail() bool {
	return q.Email != ""
}

// HasPhone 是否有电话查询
func (q *TeltentQuery) HasPhone() bool {
	return q.Phone != ""
}
