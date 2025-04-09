package errorx

import "fmt"

var _ error = &AppError{}

// AppError 自定义应用错误类型
type AppError struct {
	errCode    int
	errMsg     string
	HttpStatus int
}

// GetErrCode 返回错误码
func (e *AppError) GetErrCode() int {
	return e.errCode
}

// GetErrMsg 返回错误消息
func (e *AppError) GetErrMsg() string {
	return e.errMsg
}

// Error 实现error接口
func (e *AppError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

// WithMsg 为错误添加额外的错误信息
func (e *AppError) WithMsg(msg string) error {
	e.errMsg = fmt.Sprintf("%s, %s", e.errMsg, msg)
	return e
}

// WithError 为错误添加额外的错误
func (e *AppError) WithError(err error) error {
	e.errMsg = fmt.Sprintf("%s, %s", e.errMsg, err.Error())
	return e
}

// NewErrCodeMsg 创建带有自定义错误码和消息的错误
func NewErrCodeMsg(errCode int, errMsg string) *AppError {
	return &AppError{errCode: errCode, errMsg: errMsg}
}

// NewErrMsg 创建带有自定义消息的通用错误
func NewErrMsg(errMsg string) *AppError {
	return &AppError{errCode: ErrorUnknownCode, errMsg: errMsg}
}

// IsAppErr 判断错误是否为AppError类型
func IsAppErr(err error) bool {
	_, ok := err.(*AppError)
	return ok
}
