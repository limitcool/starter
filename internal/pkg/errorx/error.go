package errorx

import "fmt"

var _ error = &AppError{}

// AppError 自定义应用错误类型
type AppError struct {
	errorCode  int
	errorMsg   string
	httpStatus int
	causeErr   error // 添加存储被包装的原始错误
}

// GetErrorCode 返回错误码
func (e *AppError) GetErrorCode() int {
	return e.errorCode
}

// GetErrorMsg 返回错误消息
func (e *AppError) GetErrorMsg() string {
	return e.errorMsg
}

// GetHttpStatus 返回HTTP状态码
func (e *AppError) GetHttpStatus() int {
	return e.httpStatus
}

// Error 实现error接口
func (e *AppError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errorCode, e.errorMsg)
}

// WithMsg 为错误添加额外的错误信息
func (e *AppError) WithMsg(msg string) error {
	e.errorMsg = fmt.Sprintf("%s, %s", e.errorMsg, msg)
	return e
}

// WithError 为错误添加额外的错误
func (e *AppError) WithError(err error) error {
	clone := &AppError{
		errorCode:  e.errorCode,
		errorMsg:   fmt.Sprintf("%s, %s", e.errorMsg, err.Error()),
		httpStatus: e.httpStatus,
		causeErr:   err,
	}
	return clone
}

// WithNewMsgAndError 同时覆盖错误消息并追加上层错误
func (e *AppError) WithNewMsgAndError(newMsg string, err error) error {
	clone := &AppError{
		errorCode:  e.errorCode,
		errorMsg:   newMsg,
		httpStatus: e.httpStatus,
		causeErr:   err,
	}
	return clone
}

// Unwrap 实现errors.Unwrap接口，这样可以获取原始错误
func (e *AppError) Unwrap() error {
	return e.causeErr
}

// NewAppError 创建带有自定义错误码和消息的错误
func NewAppError(errCode int, errMsg string, httpStatus int) *AppError {
	return &AppError{errorCode: errCode, errorMsg: errMsg, httpStatus: httpStatus}
}

// IsAppErr 判断错误是否为AppError类型
func IsAppErr(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppErr 将普通错误转换为AppError，如果不是AppError则返回ErrUnknown
func GetAppErr(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return ErrUnknown
}
