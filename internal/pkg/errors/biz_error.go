package errors

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
)

type BizError struct {
	Code       int
	Message    string
	HTTPStatus int
	Err        error
	Stack      []Frame
}

type Frame struct {
	File string
	Line int
	Func string
}

type StackTrace []Frame

func (s StackTrace) LogValue() slog.Value {
	parts := make([]string, 0, len(s))
	for _, f := range s {
		file := trimFile(f.File)
		parts = append(parts, fmt.Sprintf("%s:%d", file, f.Line))
	}
	return slog.StringValue(strings.Join(parts, " -> "))
}

func (e *BizError) Error() string {
	return e.Message
}

func (e *BizError) Unwrap() error {
	return e.Err
}

func (e *BizError) StackValue() slog.Value {
	return StackTrace(e.Stack).LogValue()
}

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

func New(code int, message string, httpStatus int) *BizError {
	return newBizError(code, message, httpStatus, 3)
}

func Wrap(err error, code int, message string, httpStatus int) *BizError {
	// 如果已经是 BizError，直接返回，避免重复包装
	if biz, ok := errors.AsType[*BizError](err); ok {
		return biz
	}
	biz := newBizError(code, message, httpStatus, 3)
	biz.Err = err
	return biz
}

func InvalidParam(message string) *BizError {
	return New(ErrInvalidParam, message, http.StatusBadRequest)
}

func InvalidParamWrap(err error, message string) *BizError {
	return Wrap(err, ErrInvalidParam, message, http.StatusBadRequest)
}

func Internal(err error) *BizError {
	return Wrap(err, ErrInternal, "internal server error", http.StatusInternalServerError)
}

func shouldSkip(file string) bool {
	return strings.Contains(file, "runtime/") ||
		strings.Contains(file, "net/http") ||
		strings.Contains(file, "github.com/labstack/echo")
}

func trimFile(file string) string {
	_, after, ok := strings.Cut(file, "fireworks-admin/")
	if !ok {
		return file
	}
	return after
}
