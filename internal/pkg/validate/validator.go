package validate

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() echo.Validator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// 构建详细的错误信息
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				// 可以根据需要自定义错误信息格式
				errorMessages[e.Field()] = fmt.Sprintf("验证失败: %s", e.Tag())
			}
			// 返回带有详细信息的错误
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%v", errorMessages))
		}
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}
