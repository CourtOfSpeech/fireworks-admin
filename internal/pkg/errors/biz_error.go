package errors

import "net/http"

type BizError struct {
	Code       int
	Message    string
	HTTPStatus int
	Err        error
}

func (e *BizError) Error() string {
	return e.Message
}

func (e *BizError) Unwrap() error {
	return e.Err
}

func New(code int, message string, httpStatus int) *BizError {
	return &BizError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func Wrap(err error, code int, message string, httpStatus int) *BizError {
	return &BizError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

func InvalidParam(message string) *BizError {
	return New(ErrInvalidParam, message, http.StatusBadRequest)
}

func Internal(err error) *BizError {
	return Wrap(err, ErrInternal, "internal server error", http.StatusInternalServerError)
}
