package apiresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/code"
	"github.com/limitcool/starter/internal/pkg/errors"
)

// Response API标准响应结构
type Response[T any] struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 提示信息
	Data    T      `json:"data"`    // 数据
}

// Send 发送响应
func Send[T any](c *gin.Context, httpStatus int, errCode int, message string, data T) {
	// 如果没有自定义消息，则使用错误码对应的默认消息
	if message == "" {
		message = code.GetMsg(errCode)
	}

	c.JSON(httpStatus, Response[T]{
		Code:    errCode,
		Message: message,
		Data:    data,
	})
}

// Success 成功响应
func Success[T any](c *gin.Context, data T) {
	Send(c, http.StatusOK, code.Success, code.GetMsg(code.Success), data)
}

// SuccessWithMsg 带消息的成功响应
func SuccessWithMsg[T any](c *gin.Context, message string, data T) {
	Send(c, http.StatusOK, code.Success, message, data)
}

// Fail 失败响应
func Fail(c *gin.Context, errorCode int, message string) {
	Send[any](c, http.StatusOK, errorCode, message, nil)
}

// Error 错误响应
func Error(c *gin.Context, httpStatus int, errorCode int, message string) {
	Send[any](c, httpStatus, errorCode, message, nil)
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.InvalidParams)
	}
	Send[any](c, http.StatusBadRequest, code.InvalidParams, message, nil)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message ...string) {
	msg := code.GetMsg(code.ErrorInternal)
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	Send[any](c, http.StatusInternalServerError, code.ErrorInternal, msg, nil)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.UserNoLogin)
	}
	Send[any](c, http.StatusUnauthorized, code.UserNoLogin, message, nil)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.ErrorNotFound)
	}
	Send[any](c, http.StatusNotFound, code.ErrorNotFound, message, nil)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = code.GetMsg(code.AccessDenied)
	}
	Send[any](c, http.StatusForbidden, code.AccessDenied, message, nil)
}

// HandleError 统一处理错误并返回响应
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 使用errors包解析错误
	errCode, errMsg := errors.ParseError(err)

	// 根据错误类型决定HTTP状态码
	httpStatus := http.StatusOK

	if errors.IsAuthenticationFailed(err) {
		httpStatus = http.StatusUnauthorized
	} else if errors.IsPermissionDenied(err) {
		httpStatus = http.StatusForbidden
	} else if errors.IsValidationError(err) {
		httpStatus = http.StatusBadRequest
	} else if errors.IsNotFound(err) {
		httpStatus = http.StatusNotFound
	} else if errors.IsDBError(err) || errors.IsCacheError(err) {
		httpStatus = http.StatusInternalServerError
	}

	Send[any](c, httpStatus, errCode, errMsg, nil)
}
