// Package tenant 提供租户管理功能，包括租户的创建、查询、更新和删除操作。
// 本文件定义了租户模块的领域模型和常量，包括租户类型、状态和实体结构。
package tenant

import "time"

// 租户类型常量定义。
const (
	TenantTypeCompany  int8 = 1 // 企业租户
	TenantTypePersonal int8 = 2 // 个人租户
)

// 租户状态常量定义。
const (
	TenantStatusDisabled int8 = 1 // 禁用状态
	TenantStatusEnabled  int8 = 2 // 正常状态
)

// Tenant 租户实体结构体。
// 表示系统中的一个租户，包含租户的基本信息和状态。
type Tenant struct {
	ID            string    `json:"id"`             // 租户唯一标识
	CertificateNo string    `json:"certificate_no"` // 证件号码
	Name          string    `json:"name"`           // 租户名称
	Type          int8      `json:"type"`           // 租户类型：1-企业，2-个人
	ContactName   string    `json:"contact_name"`   // 联系人姓名
	Email         string    `json:"email"`          // 联系邮箱
	Phone         string    `json:"phone"`          // 联系电话
	Status        int8      `json:"status"`         // 状态：1-禁用，2-正常
	ExpiredAt     time.Time `json:"expired_at"`     // 过期时间
	CreatedAt     time.Time `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`     // 更新时间
	DeletedAt     time.Time `json:"deleted_at"`     // 删除时间
}

// IsValidStatus 检查给定的状态值是否有效。
// 有效状态值为 TenantStatusDisabled (1) 或 TenantStatusEnabled (2)。
// 返回 true 表示状态值有效，false 表示无效。
func IsValidStatus(status int8) bool {
	return status == TenantStatusDisabled || status == TenantStatusEnabled
}
