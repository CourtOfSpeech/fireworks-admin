// Package tenant 提供租户管理功能，包括租户的创建、查询、更新和删除操作。
// 本文件定义了租户模块的错误码和错误处理函数，用于统一管理业务错误。
package tenant

import (
	"net/http"
	"strings"

	entgo "github.com/speech/fireworks-admin/internal/ent"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// 租户模块错误码常量定义。
const (
	ErrCodeTenantNotFound  = 40401 // 租户不存在错误码
	ErrCodeDuplicateCertNo = 40901 // 证件号重复错误码
	ErrCodeDuplicateEmail  = 40902 // 邮箱重复错误码
	ErrCodeDuplicatePhone  = 40903 // 电话号码重复错误码
	ErrCodeInvalidStatus   = 40001 // 无效状态错误码
	ErrCodeTenantExpired   = 40002 // 租户已过期错误码
)

// ErrDuplicateCertNo 创建证件号重复错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 409 Conflict。
func ErrDuplicateCertNo(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicateCertNo, "证件号已被使用", http.StatusConflict)
}

// ErrDuplicateEmail 创建邮箱重复错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 409 Conflict。
func ErrDuplicateEmail(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicateEmail, "邮箱已被注册", http.StatusConflict)
}

// ErrDuplicatePhone 创建电话号码重复错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 409 Conflict。
func ErrDuplicatePhone(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicatePhone, "电话号码已被使用", http.StatusConflict)
}

// ErrInvalidStatus 创建无效状态错误。
// 返回业务错误，HTTP 状态码为 400 Bad Request。
func ErrInvalidStatus() error {
	return bizerr.New(ErrCodeInvalidStatus, "无效的租户状态值", http.StatusBadRequest)
}

// TenantNotFound 创建租户不存在错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 404 Not Found。
func TenantNotFound(err error) error {
	return bizerr.Wrap(err, ErrCodeTenantNotFound, "租户不存在", http.StatusNotFound)
}

// ParseRepoError 解析 Repository 层返回的错误并转换为业务错误。
// 支持处理：未找到错误、约束冲突错误（唯一键冲突）。
// 如果错误无法识别则返回原始错误。
func ParseRepoError(err error) error {
	if err == nil {
		return nil
	}

	if entgo.IsNotFound(err) {
		return TenantNotFound(err)
	}

	if entgo.IsConstraintError(err) {
		return parseConstraintError(err)
	}

	return err
}

// parseConstraintError 解析数据库约束冲突错误。
// 根据错误信息中的约束名称判断具体冲突类型：
// - uk_certificate_no: 证件号重复
// - uk_email: 邮箱重复
// - uk_phone: 电话号码重复
func parseConstraintError(err error) error {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "uk_certificate_no"):
		return ErrDuplicateCertNo(err)
	case strings.Contains(errMsg, "uk_email"):
		return ErrDuplicateEmail(err)
	case strings.Contains(errMsg, "uk_phone"):
		return ErrDuplicatePhone(err)
	default:
		return err
	}
}
