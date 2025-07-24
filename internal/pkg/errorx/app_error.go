package errorx

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// AppError 是应用程序错误的基本结构
type AppError struct {
	error
	message    string
	code       int    // 错误码
	httpStatus int    // HTTP状态码
	traceID    string // 链路追踪 ID
	cause      error  // 原始错误
}

func NewAppError(code int, msg any, httpStatus int) *AppError {

	appErr := &AppError{
		code:       code,
		httpStatus: httpStatus,
	}

	if e, ok := msg.(error); ok {
		appErr.error = e
		return appErr
	}

	if message, ok := msg.(string); ok {
		appErr.message = message
		return appErr
	}

	appErr.message = fmt.Sprintf("%v", msg)

	return appErr
}

func (e *AppError) Error() string {
	if e.message == "" {
		return e.error.Error()
	}

	return e.message
}

// Code 获取错误码
func (e *AppError) Code() int {
	return e.code
}

// HttpStatus 获取HTTP状态码
func (e *AppError) HttpStatus() int {
	return e.httpStatus
}

func (e *AppError) WithMessage(msg string) *AppError {
	e.message = msg
	return e
}

// TraceID 获取链路追踪 ID
func (e *AppError) TraceID() string {
	return e.traceID
}

// WithTraceID 设置链路追踪 ID
func (e *AppError) WithTraceID(traceID string) *AppError {
	e.traceID = traceID
	return e
}

// Is 检查错误是否是特定的应用程序错误
// 实现 errors.Is 接口
func (e *AppError) Is(target error) bool {
	if target == nil {
		return false
	}

	// 如果目标错误也是 AppError，检查错误码
	if targetErr, ok := target.(*AppError); ok {
		return e.code == targetErr.code
	}

	return false
}

// As 实现 errors.As 接口
func (e *AppError) As(target any) bool {
	if target == nil {
		return false
	}

	// 尝试将目标转换为 *AppError 指针
	appErrPtr, ok := target.(**AppError)
	if !ok {
		return false
	}

	// 设置目标为当前错误
	*appErrPtr = e
	return true
}

// Wrap 添加原始错误
func (e *AppError) Wrap(err error) *AppError {

	// 获取调用者的文件和行号
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	// 简化文件路径，只保留最后几个部分
	parts := strings.Split(file, "/")
	if len(parts) > 3 {
		file = strings.Join(parts[len(parts)-3:], "/")
	}

	// 使用 errors.WithMessage 添加位置信息
	e.cause = errors.WithMessage(err, fmt.Sprintf("[%s:%d]", file, line))

	return e
}

// Unwrap 获取原始错误
func (e *AppError) Unwrap() error {
	return e.cause
}

func FormatErrorChain(err error) string {
	if err == nil {
		return ""
	}

	// 使用 %+v 格式化错误，包含堆栈跟踪
	return fmt.Sprintf("%+v", err)
}
