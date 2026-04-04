package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// BizError 业务错误基础结构体。
// 封装了错误码、HTTP 状态码和用户友好的消息，
// 用于在业务层统一传递和处理可预期的错误信息。
type BizError struct {
	Code       int    // 业务错误码
	HTTPStatus int    // 对应的 HTTP 状态码
	Message    string // 用户友好的错误描述
	Err        error  // 原始错误（可选，用于日志记录）
}

// Error 实现 error 接口，返回用户友好的错误消息。
func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 支持 errors.Is 和 errors.As 的错误链解包。
func (e *BizError) Unwrap() error {
	return e.Err
}

// NewBizError 创建一个新的业务错误实例。
// code 为自定义业务错误码，httpStatus 为对应的 HTTP 劶态码，
// message 为面向用户的友好提示信息。
func NewBizError(code, httpStatus int, message string) *BizError {
	return &BizError{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
	}
}

// WrapBizError 创建一个包装了原始错误的业务错误实例。
// 用于将底层错误（如数据库错误）转换为带有上下文的业务错误，
// 同时保留原始错误以便日志记录和调试。
func WrapBizError(code, httpStatus int, message string, err error) *BizError {
	return &BizError{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
		Err:        err,
	}
}

// NotFoundError 资源未找到错误。
// 当查询的资源在系统中不存在时使用此错误类型。
type NotFoundError struct {
	*BizError
}

// NewNotFoundError 创建资源未找到错误。
// resource 表示资源名称（如"租户"、"用户"），key 表示查询条件。
func NewNotFoundError(resource, key string) *NotFoundError {
	return &NotFoundError{
		BizError: NewBizError(
			http.StatusNotFound,
			http.StatusNotFound,
			fmt.Sprintf("%s不存在: %s", resource, key),
		),
	}
}

// ConflictError 资源冲突错误。
// 当操作违反唯一约束或数据状态冲突时使用此错误类型。
type ConflictError struct {
	*BizError
}

// NewConflictError 创建资源冲突错误。
// reason 描述冲突的具体原因（如"证件号已存在"、"邮箱已被注册"）。
func NewConflictError(reason string) *ConflictError {
	return &ConflictError{
		BizError: NewBizError(
			http.StatusConflict,
			http.StatusConflict,
			reason,
		),
	}
}

// InvalidArgumentError 参数无效错误。
// 当请求参数不符合业务规则时使用此错误类型。
type InvalidArgumentError struct {
	*BizError
}

// NewInvalidArgumentError 创建参数无效错误。
// reason 描述参数无效的具体原因。
func NewInvalidArgumentError(reason string) *InvalidArgumentError {
	return &InvalidArgumentError{
		BizError: NewBizError(
			http.StatusBadRequest,
			http.StatusBadRequest,
			reason,
		),
	}
}

// IsBizError 判断目标错误是否为 BizError 或其子类型。
// 使用 errors.As 进行类型断言，支持错误链中的嵌套判断。
func IsBizError(err error) bool {
	var bizErr *BizError
	return errors.As(err, &bizErr)
}

// IsNotFoundError 判断目标错误是否为 NotFoundError 类型。
func IsNotFoundError(err error) bool {
	var notFound *NotFoundError
	return errors.As(err, &notFound)
}

// IsConflictError 判断目标错误是否为 ConflictError 类型。
func IsConflictError(err error) bool {
	var conflict *ConflictError
	return errors.As(err, &conflict)
}

// HTTPStatus 从错误中提取 HTTP 状态码。
// 如果是 BizError 及其子类型则返回其 HTTPStatus 字段，
// 否则默认返回 Internal Server Error (500)。
func HTTPStatus(err error) int {
	var bizErr *BizError
	if errors.As(err, &bizErr) {
		return bizErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

// Wrap 为现有错误添加上下文消息并保留原始错误链。
// 类似于 fmt.Errorf 的 %w 动词，但提供更语义化的封装方式。
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
