// Package tenant 提供租户管理功能，包括租户的创建、查询、更新和删除操作。
// 本文件定义了租户模块的错误码和错误处理函数，用于统一管理业务错误。
package tenant

import (
	"net/http"

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

// repoParser 租户模块的 Repository 错误解析器。
// 将数据库返回的未找到、约束冲突等错误转换为对应的业务错误。
var repoParser = bizerr.NewRepoErrorParser(TenantNotFound, []bizerr.ConstraintMapping{
	{Constraint: "uk_certificate_no", ErrFactory: ErrDuplicateCertNo},
	{Constraint: "uk_email", ErrFactory: ErrDuplicateEmail},
	{Constraint: "uk_phone", ErrFactory: ErrDuplicatePhone},
})
