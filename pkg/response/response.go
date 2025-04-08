package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/pkg/code"
	"github.com/limitcool/starter/pkg/errors"
)

// Response 标准响应结构
type Response struct {
	Code    int         `json:"code"`    // 错误码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code.Success,
		Message: code.GetMsg(code.Success),
		Data:    data,
	})
}

// SuccessWithMsg 带消息的成功响应
func SuccessWithMsg(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code.Success,
		Message: message,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, errorCode int, message string) {
	// 如果没有自定义消息，则使用错误码对应的默认消息
	if message == "" {
		message = code.GetMsg(errorCode)
	}
	c.JSON(http.StatusOK, Response{
		Code:    errorCode,
		Message: message,
		Data:    nil,
	})
}

// FailWithData 带数据的失败响应
func FailWithData(c *gin.Context, errorCode int, message string, data interface{}) {
	// 如果没有自定义消息，则使用错误码对应的默认消息
	if message == "" {
		message = code.GetMsg(errorCode)
	}
	c.JSON(http.StatusOK, Response{
		Code:    errorCode,
		Message: message,
		Data:    data,
	})
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.InvalidParams)
	}
	Fail(c, code.InvalidParams, message)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context) {
	Fail(c, code.ErrorUnknown, code.GetMsg(code.ErrorUnknown))
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.ErrorNotFound)
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    code.ErrorNotFound,
		Message: message,
		Data:    nil,
	})
}

// DbError 数据库错误响应
func DbError(c *gin.Context) {
	Fail(c, code.ErrorDatabase, code.GetMsg(code.ErrorDatabase))
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.UserNoLogin)
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    code.UserNoLogin,
		Message: message,
		Data:    nil,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.AccessDenied)
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    code.AccessDenied,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.InvalidParams)
	}
	c.JSON(http.StatusBadRequest, Response{
		Code:    code.InvalidParams,
		Message: message,
		Data:    nil,
	})
}

// InternalServerError 服务器内部错误
func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.ErrorInternal)
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    code.ErrorInternal,
		Message: message,
		Data:    nil,
	})
}

// ServiceUnavailable 服务不可用
func ServiceUnavailable(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.ErrorUnknown)
	}
	c.JSON(http.StatusServiceUnavailable, Response{
		Code:    code.ErrorUnknown,
		Message: message,
		Data:    nil,
	})
}

// ResponseWithCode 自定义状态码响应
func ResponseWithCode(c *gin.Context, statusCode int, errorCode int, message string, data interface{}) {
	if message == "" {
		message = code.GetMsg(errorCode)
	}
	c.JSON(statusCode, Response{
		Code:    errorCode,
		Message: message,
		Data:    data,
	})
}

// HandleError 处理错误并返回响应
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 使用errors包解析错误
	errCode, errMsg := errors.ParseError(err)
	Fail(c, errCode, errMsg)
}
