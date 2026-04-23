// Package user 提供User功能，包括User的创建、查询、更新和删除操作。
// 本文件定义了User模块的领域模型和常量，包括实体结构。
package user

import "time"

// User User实体结构体。
// 表示系统中的一个User，包含基本信息和状态。
type User struct {
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 邮箱
	Phone     string    `json:"phone"`      // 手机号
	Password  string    `json:"-"`          // 密码
	Nickname  string    `json:"nickname"`   // 昵称
	Avatar    string    `json:"avatar"`     // 头像URL
	ID        string    `json:"id"`         // 主键
	TenantID  string    `json:"tenant_id"`  // 租户ID
	Status    int8      `json:"status"`     // 状态
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
	DeletedAt time.Time `json:"deleted_at"` // 删除时间
}
