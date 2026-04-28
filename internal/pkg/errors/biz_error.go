// Package errors 提供业务错误处理功能，支持错误码、HTTP状态码和堆栈跟踪。
// 该包定义了统一的业务错误类型，便于在应用层进行错误处理和日志记录。
package errors

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
)

// BizError 表示业务逻辑错误，包含错误码、消息、HTTP状态码和堆栈信息。
// 它实现了 error 接口，并支持错误包装和日志结构化输出。
type BizError struct {
	// Code 业务错误码，用于标识具体的错误类型
	Code int
	// Message 错误消息，描述错误的具体内容
	Message string
	// HTTPStatus HTTP 响应状态码，用于 API 响应
	HTTPStatus int
	// Err 被包装的原始错误，可为 nil
	Err error
	// Stack 错误发生时的堆栈跟踪信息
	Stack []Frame
}

// Frame 表示堆栈跟踪中的单个调用帧，包含文件名、行号和函数名。
type Frame struct {
	// File 源代码文件路径
	File string
	// Line 源代码行号
	Line int
	// Func 函数名称
	Func string
}

// StackTrace 是 Frame 的切片类型，表示完整的堆栈跟踪信息。
// 该类型实现了 slog.LogValuer 接口，支持结构化日志输出。
type StackTrace []Frame

// LogValue 实现 slog.LogValuer 接口，将堆栈跟踪转换为结构化日志值。
// 返回的日志值为文件路径和行号的链式表示，格式为 "file1:line1 -> file2:line2"。
func (s StackTrace) LogValue() slog.Value {
	parts := make([]string, 0, len(s))
	for _, f := range s {
		file := trimFile(f.File)
		parts = append(parts, fmt.Sprintf("%s:%d", file, f.Line))
	}
	return slog.StringValue(strings.Join(parts, " -> "))
}

// Error 实现 error 接口，返回错误消息。
// 如果有原始错误，则包含原始错误信息。
func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 返回被包装的原始错误，支持 errors.Is 和 errors.As 进行错误链检查。
func (e *BizError) Unwrap() error {
	return e.Err
}

// StackValue 返回堆栈跟踪的结构化日志值，便于日志输出。
func (e *BizError) StackValue() slog.Value {
	return StackTrace(e.Stack).LogValue()
}

// LogValue 实现 slog.LogValuer 接口，自定义日志输出格式。
// 返回包含错误消息和原始错误的结构化日志值。
func (e *BizError) LogValue() slog.Value {
	if e.Err != nil {
		return slog.GroupValue(
			slog.String("message", e.Message),
			slog.String("cause", e.Err.Error()),
		)
	}
	return slog.StringValue(e.Message)
}

// newBizError 创建一个新的 BizError 实例，并自动捕获调用堆栈。
// skip 参数用于跳过指定数量的调用帧，通常用于排除错误创建函数本身。
func newBizError(code int, message string, httpStatus int, skip int) *BizError {

	const depth = 32

	var pcs [depth]uintptr

	n := runtime.Callers(skip, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])

	stack := make([]Frame, 0, n)

	for {

		frame, more := frames.Next()

		if shouldSkip(frame.File) {

			if !more {
				break
			}

			continue
		}

		stack = append(stack, Frame{
			File: frame.File,
			Line: frame.Line,
			Func: frame.Function,
		})

		if !more {
			break
		}
	}

	return &BizError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Stack:      stack,
	}
}

// New 创建一个新的 BizError 实例，包含指定的错误码、消息和 HTTP 状态码。
// 该函数会自动捕获调用堆栈，用于错误追踪和日志记录。
func New(code int, message string, httpStatus int) *BizError {
	return newBizError(code, message, httpStatus, 3)
}

// Wrap 将一个现有错误包装为 BizError，添加错误码、消息和 HTTP 状态码。
// 如果传入的错误为 nil，则返回 nil。
// 如果传入的错误已经是 BizError 类型，则直接返回原错误，避免重复包装。
func Wrap(err error, code int, message string, httpStatus int) *BizError {
	if err == nil {
		return nil
	}
	// 如果已经是 BizError，直接返回，避免重复包装
	if biz, ok := errors.AsType[*BizError](err); ok {
		return biz
	}
	biz := newBizError(code, message, httpStatus, 3)
	biz.Err = err
	return biz
}

// InvalidParam 创建一个表示参数无效的 BizError。
// 使用 ErrInvalidParam 错误码和 HTTP 400 Bad Request 状态码。
func InvalidParam(message string) *BizError {
	return New(ErrInvalidParam, message, http.StatusBadRequest)
}

// InvalidParamWrap 将错误包装为参数无效的 BizError。
// 使用 ErrInvalidParam 错误码和 HTTP 400 Bad Request 状态码。
func InvalidParamWrap(err error, message string) *BizError {
	return Wrap(err, ErrInvalidParam, message, http.StatusBadRequest)
}

// Internal 将错误包装为内部服务器错误的 BizError。
// 使用 ErrInternal 错误码和 HTTP 500 Internal Server Error 状态码。
func Internal(err error) *BizError {
	return Wrap(err, ErrInternal, "internal server error", http.StatusInternalServerError)
}

// shouldSkip 判断给定的文件路径是否应该从堆栈跟踪中跳过。
// 跳过运行时、标准库和框架相关的文件，只保留业务代码的堆栈信息。
func shouldSkip(file string) bool {
	return strings.Contains(file, "runtime/") ||
		strings.Contains(file, "net/http") ||
		strings.Contains(file, "github.com/labstack/echo")
}

// trimFile 裁剪文件路径，移除项目根目录前缀，使日志输出更加简洁。
// 如果文件路径不包含项目名称，则返回原始路径。
func trimFile(file string) string {
	_, after, ok := strings.Cut(file, "fireworks-admin/")
	if !ok {
		return file
	}
	return after
}
