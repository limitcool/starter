package errorx

import (
	"fmt"
	"runtime"
	"strings"
)

var _ error = &AppError{}

// AppError 自定义应用错误类型
type AppError struct {
	errorCode  int
	errorMsg   string
	httpStatus int
	causeErr   error     // 添加存储被包装的原始错误
	stackTrace []uintptr // 存储错误发生时的堆栈信息
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
		errorMsg:   e.errorMsg,
		httpStatus: e.httpStatus,
		causeErr:   err,
		stackTrace: captureStackTrace(), // 捕获当前堆栈
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
		stackTrace: captureStackTrace(), // 捕获当前堆栈
	}
	return clone
}

// Unwrap 实现errors.Unwrap接口，这样可以获取原始错误
func (e *AppError) Unwrap() error {
	return e.causeErr
}

// GetStackTrace 获取格式化后的堆栈信息
func (e *AppError) GetStackTrace() string {
	return formatStackTrace(e.stackTrace)
}

// NewAppError 创建带有自定义错误码和消息的错误
func NewAppError(errCode int, errMsg string, httpStatus int) *AppError {
	return &AppError{
		errorCode:  errCode,
		errorMsg:   errMsg,
		httpStatus: httpStatus,
		stackTrace: captureStackTrace(), // 捕获创建时的堆栈
	}
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

// Format 实现fmt.Formatter接口，用于格式化错误输出
func (e *AppError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// 详细模式，打印包括堆栈信息在内的所有内容
			fmt.Fprintf(s, "错误: %s (错误码: %d)\n", e.errorMsg, e.errorCode)
			if e.causeErr != nil {
				fmt.Fprintf(s, "原因: %s\n", e.causeErr.Error())
			}

			// 根据日志级别决定是否输出堆栈
			// 移除对特定配置的依赖，始终在详细模式下显示堆栈
			fmt.Fprintf(s, "堆栈信息:\n%s", e.GetStackTrace())
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%s", e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// 捕获当前堆栈信息
func captureStackTrace() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	// 跳过这个函数本身和调用者
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

// 格式化堆栈信息
func formatStackTrace(stackTrace []uintptr) string {
	frames := runtime.CallersFrames(stackTrace)
	var result string

	// 默认最大堆栈帧数，可以考虑通过环境变量设置
	maxFrames := 10

	i := 1
	for {
		if i > maxFrames {
			result += fmt.Sprintf("... 更多堆栈被省略 ...\n")
			break
		}

		frame, more := frames.Next()
		// 过滤掉不需要的框架函数
		if !isInternalFrame(frame.Function) {
			// 添加序号和更美观的格式
			result += fmt.Sprintf("%d. %s\n   at %s:%d\n", i, frame.Function, frame.File, frame.Line)
			i++
		}
		if !more {
			break
		}
	}
	return result
}

// 判断是否是内部框架函数，过滤掉不需要显示的框架函数
func isInternalFrame(function string) bool {
	// 可以根据需要过滤掉一些函数
	prefixesToFilter := []string{
		"runtime.",                               // 过滤运行时函数
		"testing.",                               // 过滤测试函数
		"github.com/gin-gonic/gin.handlerNative", // 过滤一些gin内部函数
	}

	for _, prefix := range prefixesToFilter {
		if strings.HasPrefix(function, prefix) {
			return true
		}
	}
	return false
}
