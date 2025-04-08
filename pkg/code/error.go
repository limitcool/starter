package code

import "fmt"

var _ error = &CodeError{}

type CodeError struct {
	errCode int
	errMsg  string
}

// 返回给前端的错误码
func (e *CodeError) GetErrCode() int {
	return e.errCode
}

// 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

func NewErrCodeMsg(errCode int, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}
func NewErrCode(errCode int) *CodeError {
	return &CodeError{errCode: errCode, errMsg: GetMsg(errCode)}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: ErrorUnknown, errMsg: errMsg}
}

// IsCodeErr 判断错误码是否存在
func IsCodeErr(errcode int) bool {
	if _, ok := MsgFlags[errcode]; ok {
		return true
	} else {
		return false
	}
}

// IsErrCode 判断错误是否为CodeError类型
func IsErrCode(err error) bool {
	_, ok := err.(*CodeError)
	return ok
}
