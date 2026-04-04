package validator

import (
	"fmt"
	"net/http"

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
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				errorMessages[e.Field()] = fmt.Sprintf("验证失败: %s", e.Tag())
			}
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%v", errorMessages))
		}
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}
