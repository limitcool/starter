package errorx

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// AppError 是应用程序错误的基本结构
type AppError struct {
	code     int    // 错误码
	msg      string // 错误消息
	httpCode int    // HTTP状态码
	traceID  string // 链路追踪 ID
	orig     error  // 原始错误
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	return e.msg
}

// GetErrorCode 获取错误码
func (e *AppError) GetErrorCode() int {
	return e.code
}

// GetErrorMsg 获取错误消息
func (e *AppError) GetErrorMsg() string {
	return e.msg
}

// GetHttpStatus 获取HTTP状态码
func (e *AppError) GetHttpStatus() int {
	return e.httpCode
}

// WithMsg 添加错误消息
func (e *AppError) WithMsg(msg string) *AppError {
	e.msg = msg
	return e
}

// WithError 添加原始错误
func (e *AppError) WithError(err error) *AppError {
	// 创建新的 AppError 实例，避免修改原始错误
	newErr := &AppError{
		code:     e.code,
		msg:      e.msg,
		httpCode: e.httpCode,
		orig:     err,
	}

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
	newErr.orig = errors.WithMessage(err, fmt.Sprintf("[%s:%d]", file, line))

	return newErr
}

// GetTraceID 获取链路追踪 ID
func (e *AppError) GetTraceID() string {
	return e.traceID
}

// SetTraceID 设置链路追踪 ID
func (e *AppError) SetTraceID(traceID string) *AppError {
	e.traceID = traceID
	return e
}

// Unwrap 获取原始错误
func (e *AppError) Unwrap() error {
	return e.orig
}

// NewAppError 创建一个新的应用程序错误
func NewAppError(code int, msg string, httpCode int) *AppError {
	return &AppError{
		code:     code,
		msg:      msg,
		httpCode: httpCode,
	}
}

// Errorf 创建一个格式化的应用程序错误
// 用法: Errorf(ErrNotFound, "用户ID %d 不存在", id)
func Errorf(baseErr *AppError, format string, args ...any) error {
	return &AppError{
		code:     baseErr.code,
		msg:      fmt.Sprintf(format, args...),
		httpCode: baseErr.httpCode,
	}
}

// WrapError 包装错误并添加上下文信息和位置信息
// 用法: WrapError(err, "查询用户失败")
// 如果不需要添加消息，可以传入空字符串: WrapError(err, "")
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

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

	// 位置信息
	location := fmt.Sprintf("[%s:%d]", file, line)

	// 如果没有提供消息，只添加位置信息
	if message == "" {
		return errors.WithMessage(err, location)
	}

	// 添加位置信息和消息
	return errors.WithMessage(err, fmt.Sprintf("%s %s", message, location))
}

// FormatErrorChain 格式化错误链，包括位置信息
// 用法: FormatErrorChain(err)
func FormatErrorChain(err error) string {
	if err == nil {
		return ""
	}

	// 使用 %+v 格式化错误，包含堆栈跟踪
	return fmt.Sprintf("%+v", err)
}

// GetUserMessage 获取用户友好的错误消息
// 用法: GetUserMessage(err)
func GetUserMessage(err error) string {
	if err == nil {
		return ""
	}

	// 如果是 AppError，返回错误消息
	if appErr, ok := err.(*AppError); ok {
		return appErr.msg
	}

	// 默认返回通用错误消息
	return "系统发生错误，请稍后再试"
}

// GetErrorCode 获取错误码
// 用法: GetErrorCode(err)
func GetErrorCode(err error) int {
	if err == nil {
		return 0
	}

	// 如果是 AppError，返回错误码
	if appErr, ok := err.(*AppError); ok {
		return appErr.code
	}

	// 默认返回未知错误码
	return 500000 // ErrorUnknownCode
}

// GetHttpStatus 获取HTTP状态码
// 用法: GetHttpStatus(err)
func GetHttpStatus(err error) int {
	if err == nil {
		return 200
	}

	// 如果是 AppError，返回HTTP状态码
	if appErr, ok := err.(*AppError); ok {
		return appErr.httpCode
	}

	// 默认返回500
	return 500
}

// IsAppErr 检查错误是否是应用程序错误
func IsAppErr(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// Is 检查错误是否是特定的应用程序错误
func Is(err error, target *AppError) bool {
	if err == nil || target == nil {
		return false
	}

	// 如果是 AppError，检查错误码
	if appErr, ok := err.(*AppError); ok {
		return appErr.code == target.code
	}

	return false
}

// 默认最大堆栈帧数
var maxStackFrames = 32

// SetMaxStackFrames 设置最大堆栈帧数
func SetMaxStackFrames(frames int) {
	if frames > 0 {
		maxStackFrames = frames
	}
}

// GetError 根据错误码获取错误
// 用法: GetError(10001)
func GetError(code int) *AppError {
	// 默认返回未知错误
	switch code {
	case SuccessCode:
		return Success
	case ErrorUnknownCode:
		return ErrUnknown
	case ErrorInvalidParamsCode:
		return ErrInvalidParams
	case ErrorInternalCode:
		return ErrInternal
	case ErrorUnauthorizedCode:
		return ErrUnauthorized
	case ErrorForbiddenCode:
		return ErrForbidden
	case ErrorNotFoundCode:
		return ErrNotFound
	case ErrorTimeoutCode:
		return ErrTimeout
	case ErrorTooManyRequestsCode:
		return ErrTooManyRequests
	case ErrorAccessDeniedCode:
		return ErrAccessDenied
	case ErrorUserAuthFailedCode:
		return ErrUserAuthFailed
	case ErrorCasbinServiceCode:
		return ErrCasbinService
	case ErrorFileStorageCode:
		return ErrFileStorage
	case ErrorUserNotFoundCode:
		return ErrUserNotFound
	case ErrorInvalidCredentialsCode:
		return ErrInvalidCredentials
	case ErrorUserDisabledCode:
		return ErrUserDisabled
	case ErrorUserExistsCode:
		return ErrUserExists
	case ErrorPasswordExpiredCode:
		return ErrPasswordExpired
	case ErrorUserPasswordErrorCode:
		return ErrUserPasswordError
	case ErrorUserTokenErrorCode:
		return ErrUserTokenError
	case ErrorUserNoLoginCode:
		return ErrUserNoLogin
	case ErrorDatabaseCode:
		return ErrDatabase
	case ErrorDatabaseQueryCode:
		return ErrDatabaseQuery
	case ErrorDatabaseInsertCode:
		return ErrDatabaseInsert
	case ErrorDatabaseUpdateCode:
		return ErrDatabaseUpdate
	case ErrorDatabaseDeleteCode:
		return ErrDatabaseDelete
	case ErrorDatabaseConnectionCode:
		return ErrDatabaseConnection
	case ErrorDatabaseTransactionCode:
		return ErrDatabaseTransaction
	default:
		return ErrUnknown
	}
}
