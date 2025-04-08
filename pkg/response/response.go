package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/pkg/code"
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
		Message: "操作成功",
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
func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// FailWithData 带数据的失败响应
func FailWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, message string) {
	Fail(c, code.InvalidParams, message)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context) {
	Fail(c, code.ErrorUnknown, "服务器内部错误")
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	Fail(c, code.ErrorNotFound, message)
}

// DbError 数据库错误响应
func DbError(c *gin.Context) {
	Fail(c, code.ErrorDatabase, "数据库操作失败")
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    code.UserNoLogin,
		Message: message,
		Data:    nil,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    code.UserNoPermission,
		Message: message,
		Data:    nil,
	})
}
