package tenant

import (
	"fmt"

	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// 租户模块业务错误码定义。
// 使用 4xxxx 范围作为租户模块的错误码段，便于错误分类和追踪。
const (
	ErrCodeTenantNotFound   = 40401 // 租户不存在
	ErrCodeDuplicateCertNo  = 40901 // 证件号重复
	ErrCodeDuplicateEmail   = 40902 // 邮箱重复
	ErrCodeDuplicatePhone   = 40903 // 电话重复
	ErrCodeInvalidStatus    = 40001 // 状态无效
	ErrCodeTenantExpired    = 40002 // 租户已过期
)

// ErrTenantNotFound 表示租户记录不存在。
// 在 GetByID、Update、Delete 操作中当目标租户未找到时返回此错误。
var ErrTenantNotFound = bizerr.NewNotFoundError("租户", "ID")

// ErrDuplicateCertNo 表示证件号已存在。
// 创建或更新租户时，若证件号与已有记录冲突则返回此错误。
var ErrDuplicateCertNo = bizerr.NewConflictError("证件号已被使用")

// ErrDuplicateEmail 表示邮箱已存在。
var ErrDuplicateEmail = bizerr.NewConflictError("邮箱已被注册")

// ErrDuplicatePhone 表示电话号码已存在。
var ErrDuplicatePhone = bizerr.NewConflictError("电话号码已被使用")

// ErrInvalidStatus 表示状态值无效。
// 当传入的状态值不在允许范围内（非禁用/正常）时返回此错误。
var ErrInvalidStatus = bizerr.NewInvalidArgumentError("无效的租户状态值")

// NewTenantNotFound 根据给定 ID 创建具体的"租户不存在"错误实例。
// 用于在 Service 层返回带有具体查询条件的 NotFoundError。
func NewTenantNotFound(id string) error {
	return bizerr.NewNotFoundError("租户", fmt.Sprintf("id=%s", id))
}

// IsDuplicateKeyError 判断数据库错误是否为唯一约束冲突（重复键）。
// Ent/PostgreSQL 的唯一约束违反会返回包含特定关键词的错误，
// 此函数用于在 Repository 层识别此类错误并转换为业务错误。
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return containsAny(errMsg,
		"unique constraint",
		"duplicate key",
		"SQLSTATE 23505",
	)
}

// containsAny 检查字符串 s 是否包含 substrs 中的任意一个子串。
// 忽略大小写进行匹配。
func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if containsIgnoreCase(s, sub) {
			return true
		}
	}
	return false
}

// containsIgnoreCase 忽略大小写检查 s 是否包含 sub。
func containsIgnoreCase(s, sub string) bool {
	sLower := make([]byte, len(s))
	subLower := make([]byte, len(sub))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		sLower[i] = c
	}
	for i := 0; i < len(sub); i++ {
		c := sub[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		subLower[i] = c
	}
	return contains(string(sLower), string(subLower))
}

// contains 检查 s 是否包含 sub（标准库简单实现，避免引入 strings 包的额外依赖）。
func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
