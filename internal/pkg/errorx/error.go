package errorx

import "fmt"

var _ error = &CodeError{}

// CodeError 自定义带错误码的错误类型
type CodeError struct {
	errCode int
	errMsg  string
}

// GetErrCode 返回错误码
func (e *CodeError) GetErrCode() int {
	return e.errCode
}

// GetErrMsg 返回错误消息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

// Error 实现error接口
func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

// NewErrCodeMsg 创建带有自定义错误码和消息的错误
func NewErrCodeMsg(errCode int, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}

// NewErrCode 创建带有错误码的错误，自动获取对应的消息
func NewErrCode(errCode int) *CodeError {
	return &CodeError{errCode: errCode, errMsg: GetMsg(errCode)}
}

// NewErrMsg 创建带有自定义消息的通用错误
func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: ErrorUnknown, errMsg: errMsg}
}

// IsCodeErr 判断错误码是否存在
func IsCodeErr(errcode int) bool {
	_, ok := MsgFlags[errcode]
	return ok
}

// IsErrCode 判断错误是否为CodeError类型
func IsErrCode(err error) bool {
	_, ok := err.(*CodeError)
	return ok
}
