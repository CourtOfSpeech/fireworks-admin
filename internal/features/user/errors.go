// Package user 提供用户管理功能，包括用户的创建、查询、更新和删除操作。
// 本文件定义了用户模块的错误码和错误处理函数，用于统一管理业务错误。
package user

import (
	"net/http"

	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// 用户模块错误码常量定义。
const (
	ErrCodeUserNotFound      = 40401 // 用户不存在错误码
	ErrCodeDuplicateUsername = 40901 // 用户名重复错误码
	ErrCodeDuplicateEmail    = 40902 // 邮箱重复错误码
	ErrCodeDuplicatePhone    = 40903 // 手机号重复错误码
	ErrCodeLoginFailed       = 40101 // 登录失败错误码
	ErrCodeTenantMismatch    = 40301 // 租户不匹配错误码
	ErrCodeUserDisabled      = 40302 // 用户被禁用错误码
)

// ErrDuplicateUsername 创建用户名重复错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 409 Conflict。
func ErrDuplicateUsername(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicateUsername, "用户名已被使用", http.StatusConflict)
}

// ErrDuplicateEmail 创建邮箱重复错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 409 Conflict。
func ErrDuplicateEmail(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicateEmail, "邮箱已被注册", http.StatusConflict)
}

// ErrDuplicatePhone 创建手机号重复错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 409 Conflict。
func ErrDuplicatePhone(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicatePhone, "手机号已被使用", http.StatusConflict)
}

// ErrUserNotFound 创建用户不存在错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 404 Not Found。
func ErrUserNotFound(err error) error {
	return bizerr.Wrap(err, ErrCodeUserNotFound, "用户不存在", http.StatusNotFound)
}

// ErrLoginFailed 创建登录失败错误。
// 返回业务错误，HTTP 状态码为 401 Unauthorized。
func ErrLoginFailed() error {
	return bizerr.New(ErrCodeLoginFailed, "用户名或密码错误", http.StatusUnauthorized)
}

// ErrTenantMismatch 创建租户不匹配错误。
// 返回业务错误，HTTP 状态码为 403 Forbidden。
func ErrTenantMismatch() error {
	return bizerr.New(ErrCodeTenantMismatch, "租户信息不匹配", http.StatusForbidden)
}

// ErrUserDisabled 创建用户被禁用错误。
// 返回业务错误，HTTP 状态码为 403 Forbidden。
func ErrUserDisabled() error {
	return bizerr.New(ErrCodeUserDisabled, "用户已被禁用", http.StatusForbidden)
}

// repoParser 用户模块的 Repository 错误解析器。
// 将数据库返回的未找到、约束冲突等错误转换为对应的业务错误。
var repoParser = bizerr.NewRepoErrorParser(ErrUserNotFound, []bizerr.ConstraintMapping{
	{Constraint: "uk_username", ErrFactory: ErrDuplicateUsername},
	{Constraint: "uk_email", ErrFactory: ErrDuplicateEmail},
	{Constraint: "uk_phone", ErrFactory: ErrDuplicatePhone},
})
