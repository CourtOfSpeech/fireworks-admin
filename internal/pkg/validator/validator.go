// Package validator 提供了基于 go-playground/validator 的自定义验证器实现。
// 该包实现了 echo.Validator 接口，用于验证请求数据并返回友好的中文错误信息。
package validator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// CustomValidator 是自定义验证器结构体，实现了 echo.Validator 接口。
// 它封装了 go-playground/validator 库的 Validate 实例，提供数据验证功能。
type CustomValidator struct {
	validate *validator.Validate
}

// NewValidator 创建并返回一个新的 CustomValidator 实例。
// 该函数初始化内部的 validator.Validate 实例，用于后续的数据验证操作。
// 返回值实现了 echo.Validator 接口，可直接用于 Echo 框架的验证器设置。
func NewValidator() echo.Validator {
	return &CustomValidator{
		validate: validator.New(),
	}
}

// Validate 验证给定的数据结构是否符合其定义的验证规则。
// 该方法会检查结构体字段的标签（如 required、email、min 等），
// 并在验证失败时返回包含中文错误信息的 HTTP 错误。
// 参数 i 为需要验证的数据结构，通常为请求绑定的结构体指针。
// 返回 nil 表示验证通过，否则返回相应的错误信息。
func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validate.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var msgs []string
			for _, e := range validationErrors {
				msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field(), cv.tagToMessage(e.Tag())))
			}
			return echo.NewHTTPError(http.StatusBadRequest, strings.Join(msgs, "; "))
		}
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

// tagToMessage 将验证标签转换为用户友好的中文错误信息。
// 该方法支持常见的验证标签，如 required、email、min、max、oneof 等。
// 参数 tag 为验证器返回的错误标签名称。
// 返回对应的中文错误描述，如果标签未定义则返回原始标签名。
func (cv *CustomValidator) tagToMessage(tag string) string {
	switch tag {
	case "required":
		return "必填"
	case "email":
		return "邮箱格式不正确"
	case "min":
		return "长度不足"
	case "max":
		return "超出最大长度"
	case "oneof":
		return "值不在允许范围内"
	default:
		return tag
	}
}
