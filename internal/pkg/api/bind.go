package api

import (
	"github.com/labstack/echo/v5"
	bizerr "github.com/speech/fireworks-admin/internal/pkg/errors"
)

// BindAndValidate 绑定并验证请求数据。
// 该函数将 Echo 上下文中的请求数据绑定到 target 指针，
// 然后执行结构体标签定义的验证规则。
// 绑定失败返回参数无效错误，验证失败返回验证错误。
// c 是 Echo 上下文，target 是请求数据绑定目标（必须为指针），
// errMsg 是绑定失败时的错误消息。
func BindAndValidate(c *echo.Context, target any, errMsg string) error {
	if err := c.Bind(target); err != nil {
		return bizerr.InvalidParamWrap(err, errMsg)
	}
	if err := c.Validate(target); err != nil {
		return err
	}
	return nil
}
