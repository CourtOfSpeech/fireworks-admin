package tenant

import (
	"time"

	"github.com/speech/fireworks-admin/internal/pkg/api"
)

// CreateTenantReq 创建租户请求
type CreateTenantReq struct {
	CertificateNo string    `json:"certificate_no" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	Type          int8      `json:"type" validate:"required,min=1,max=2"`
	ContactName   string    `json:"contact_name" validate:"required"`
	Email         string    `json:"email" validate:"required,email"`
	Phone         string    `json:"phone" validate:"required"`
	ExpiredAt     time.Time `json:"expired_at" validate:"required"`
	Status        int8      `json:"status" validate:"required,min=1,max=2"`
}

// UpdateTenantReq 更新租户请求
type UpdateTenantReq struct {
	CertificateNo *string    `json:"certificate_no" validate:"omitempty"`
	Name          *string    `json:"name" validate:"omitempty"`
	Type          *int8      `json:"type" validate:"omitempty,min=1,max=2"`
	ContactName   *string    `json:"contact_name" validate:"omitempty"`
	Email         *string    `json:"email" validate:"omitempty,email"`
	Phone         *string    `json:"phone" validate:"omitempty"`
	ExpiredAt     *time.Time `json:"expired_at" validate:"omitempty"`
	Status        *int8      `json:"status" validate:"omitempty,min=0,max=1"`
}

// TenantQuery 租户查询条件
type TenantQuery struct {
	api.PageQuery
	Keyword string `query:"keyword"` // 关键字：模糊匹配证件号码（前缀）、名称
	Status  *int8  `query:"status"`  // 状态：精确匹配
	Email   string `query:"email"`   // 邮箱：精确匹配
	Phone   string `query:"phone"`   // 电话：精确匹配
}

// HasKeyword 是否有关键字查询
func (q *TenantQuery) HasKeyword() bool {
	return q.Keyword != ""
}

// HasStatus 是否有状态查询
func (q *TenantQuery) HasStatus() bool {
	return q.Status != nil
}

// HasEmail 是否有邮箱查询
func (q *TenantQuery) HasEmail() bool {
	return q.Email != ""
}

// HasPhone 是否有电话查询
func (q *TenantQuery) HasPhone() bool {
	return q.Phone != ""
}
