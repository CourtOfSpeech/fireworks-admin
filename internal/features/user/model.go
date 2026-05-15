// Package user 提供用户管理功能，包括用户的创建、查询、更新和删除操作。
// 本文件定义了用户模块的领域模型和常量，包括用户状态和实体结构。
package user

import (
	"regexp"
	"time"
)

// 用户状态常量定义。
const (
	UserStatusDisabled int8 = 1 // 禁用状态
	UserStatusEnabled  int8 = 2 // 正常状态
)

// IsValidUserStatus 检查给定的用户状态值是否有效。
// 有效状态值为 UserStatusDisabled (1) 或 UserStatusEnabled (2)。
func IsValidUserStatus(status int8) bool {
	return status == UserStatusDisabled || status == UserStatusEnabled
}

// emailRegex 邮箱格式正则表达式。
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// phoneRegex 手机号格式正则表达式（中国大陆手机号）。
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// IsEmail 判断字符串是否为邮箱格式。
// 参数 s 为待判断的字符串。
// 返回 true 表示是邮箱格式，false 表示不是。
func IsEmail(s string) bool {
	return emailRegex.MatchString(s)
}

// IsPhone 判断字符串是否为手机号格式。
// 参数 s 为待判断的字符串。
// 返回 true 表示是手机号格式，false 表示不是。
func IsPhone(s string) bool {
	return phoneRegex.MatchString(s)
}

// User 用户实体结构体。
// 表示系统中的一个用户，包含基本信息和状态。
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
