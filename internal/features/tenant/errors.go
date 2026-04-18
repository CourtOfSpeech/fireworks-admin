package tenant

import (
	"net/http"
	"strings"

	entgo "github.com/speech/fireworks-admin/internal/ent"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

const (
	ErrCodeTenantNotFound  = 40401
	ErrCodeDuplicateCertNo = 40901
	ErrCodeDuplicateEmail  = 40902
	ErrCodeDuplicatePhone  = 40903
	ErrCodeInvalidStatus   = 40001
	ErrCodeTenantExpired   = 40002
)

func ErrDuplicateCertNo(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicateCertNo, "证件号已被使用", http.StatusConflict)
}

func ErrDuplicateEmail(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicateEmail, "邮箱已被注册", http.StatusConflict)
}

func ErrDuplicatePhone(err error) error {
	return bizerr.Wrap(err, ErrCodeDuplicatePhone, "电话号码已被使用", http.StatusConflict)
}

func ErrInvalidStatus() error {
	return bizerr.New(ErrCodeInvalidStatus, "无效的租户状态值", http.StatusBadRequest)
}

func TenantNotFound(err error) error {
	return bizerr.Wrap(err, ErrCodeTenantNotFound, "租户不存在", http.StatusNotFound)
}

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
