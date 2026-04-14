package validator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type CustomValidator struct {
	validate *validator.Validate
}

func NewValidator() echo.Validator {
	return &CustomValidator{
		validate: validator.New(),
	}
}

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
