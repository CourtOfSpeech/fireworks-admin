package tenant

import "time"

// 租户类型
const (
	TenantTypeCompany  int8 = 1 // 企业
	TenantTypePersonal int8 = 2 // 个人
)

// 租户状态
const (
	TenantStatusDisabled int8 = 1 // 禁用
	TenantStatusEnabled  int8 = 2 // 正常
)

// Tenant 租户实体
type Tenant struct {
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

func IsValidStatus(status int8) bool {
	return status == TenantStatusDisabled || status == TenantStatusEnabled
}
