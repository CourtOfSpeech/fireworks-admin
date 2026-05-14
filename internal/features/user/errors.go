// Package user 提供User功能，包括User的创建、查询、更新和删除操作。
// 本文件定义了User模块的错误码和错误处理函数，用于统一管理业务错误。
package user

import (
	"net/http"
	"strings"

	entgo "github.com/speech/fireworks-admin/internal/ent"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// User模块错误码常量定义。
const (
	ErrCodeUserNotFound      = 40401 // User不存在错误码
	ErrCodeDuplicateusername = 40901 // username重复错误码
	ErrCodeDuplicateemail    = 40902 // email重复错误码
	ErrCodeDuplicatephone    = 40903 // phone重复错误码
	ErrCodeLoginFailed       = 40101 // 登录失败错误码
	ErrCodeTenantMismatch    = 40301 // 租户不匹配错误码
	ErrCodeUserDisabled      = 40302 // 用户被禁用错误码
)

// ErrDuplicateusername 创建username重复错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 409 Conflict。
func ErrDuplicateusername(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicateusername, "username已被使用", http.StatusConflict)
}

// ErrDuplicateemail 创建email重复错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 409 Conflict。
func ErrDuplicateemail(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicateemail, "email已被使用", http.StatusConflict)
}

// ErrDuplicatephone 创建phone重复错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 409 Conflict。
func ErrDuplicatephone(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicatephone, "phone已被使用", http.StatusConflict)
}

// UserNotFound 创建User不存在错误。
// 参数 err 为原始错误，返回包装后的业务错误，HTTP 状态码为 404 Not Found。
func UserNotFound(err error) error {
	return bizerr.Wrap(err, ErrCodeUserNotFound, "User不存在", http.StatusNotFound)
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

// ParseRepoError 解析 Repository 层返回的错误并转换为业务错误。
// 支持处理：未找到错误、约束冲突错误（唯一键冲突）。
// 如果错误无法识别则返回原始错误。
func ParseRepoError(err error) error {
	if err == nil {
		return nil
	}

	if entgo.IsNotFound(err) {
		return UserNotFound(err)
	}

	if entgo.IsConstraintError(err) {
		return parseConstraintError(err)
	}

	return err
}

// parseConstraintError 解析数据库约束冲突错误。
// 根据错误信息中的约束名称判断具体冲突类型：
// - uk_username: username重复
// - uk_email: email重复
// - uk_phone: phone重复
func parseConstraintError(err error) error {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "uk_username"):
		return ErrDuplicateusername(err)
	case strings.Contains(errMsg, "uk_email"):
		return ErrDuplicateemail(err)
	case strings.Contains(errMsg, "uk_phone"):
		return ErrDuplicatephone(err)
	default:
		return err
	}
}
