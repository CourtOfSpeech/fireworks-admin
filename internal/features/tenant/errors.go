package tenant

import (
	"fmt"
	"net/http"
	"strings"

	entgo "github.com/speech/fireworks-admin/internal/ent"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// 租户模块业务错误码定义。
// 使用 4xxxx 范围作为租户模块的错误码段，便于错误分类和追踪。
const (
	ErrCodeTenantNotFound  = 40401 // 租户不存在
	ErrCodeDuplicateCertNo = 40901 // 证件号重复
	ErrCodeDuplicateEmail  = 40902 // 邮箱重复
	ErrCodeDuplicatePhone  = 40903 // 电话重复
	ErrCodeInvalidStatus   = 40001 // 状态无效
	ErrCodeTenantExpired   = 40002 // 租户已过期
)

// ErrTenantNotFound 表示租户记录不存在。
// 在 GetByID、Update、Delete 操作中当目标租户未找到时返回此错误。
func ErrTenantNotFound() error {
	return bizerr.New(ErrCodeTenantNotFound, "租户不存在", http.StatusNotFound)
}

// ErrDuplicateCertNo 表示证件号已存在。
// 创建或更新租户时，若证件号与已有记录冲突则返回此错误。
func ErrDuplicateCertNo() error {
	return bizerr.New(ErrCodeDuplicateCertNo, "证件号已被使用", http.StatusConflict)
}

// ErrDuplicateEmail 表示邮箱已存在。
func ErrDuplicateEmail() error {
	return bizerr.New(ErrCodeDuplicateEmail, "邮箱已被注册", http.StatusConflict)
}

// ErrDuplicatePhone 表示电话号码已存在。
func ErrDuplicatePhone() error {
	return bizerr.New(ErrCodeDuplicatePhone, "电话号码已被使用", http.StatusConflict)
}

// ErrInvalidStatus 表示状态值无效。
// 当传入的状态值不在允许范围内（非禁用/正常）时返回此错误。
func ErrInvalidStatus() error {
	return bizerr.New(ErrCodeInvalidStatus, "无效的租户状态值", http.StatusBadRequest)
}

// NewTenantNotFound 根据给定 ID 创建具体的"租户不存在"错误实例。
// 用于在 Service 层返回带有具体查询条件的 NotFoundError。
func NewTenantNotFound(id string) error {
	return bizerr.New(ErrCodeTenantNotFound, fmt.Sprintf("租户不存在: id=%s", id), http.StatusNotFound)
}

// ParseRepoError 解析 Repository 层返回的错误，转换为对应的业务错误。
// 此函数封装了 ent 错误类型的判断逻辑，使 Service 层无需依赖 ent 包。
// 参数 id 用于 NotFound 错误时提供具体的查询条件信息。
func ParseRepoError(err error, id string) error {
	if err == nil {
		return nil
	}

	if entgo.IsNotFound(err) {
		return NewTenantNotFound(id)
	}

	if entgo.IsConstraintError(err) {
		return parseConstraintError(err)
	}

	return err
}

// parseConstraintError 解析约束错误，返回对应的领域错误。
func parseConstraintError(err error) error {
	if !entgo.IsConstraintError(err) {
		return err
	}
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "uk_certificate_no"):
		return ErrDuplicateCertNo()
	case strings.Contains(errMsg, "uk_email"):
		return ErrDuplicateEmail()
	case strings.Contains(errMsg, "uk_phone"):
		return ErrDuplicatePhone()
	default:
		return err
	}
}
