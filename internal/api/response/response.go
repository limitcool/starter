// Package response 提供API响应相关的功能
// 此包整合了以前 internal/pkg/apiresponse 包的功能，并增加了分页支持
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/errorx"
)

// Response API标准响应结构
type Response[T any] struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 提示信息
	Data    T      `json:"data"`    // 数据
}

// PageResult 分页结果
type PageResult[T any] struct {
	Total    int64 `json:"total"`     // 总记录数
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 每页大小
	List     T     `json:"list"`      // 数据列表
}

// NewPageResult 创建分页结果
func NewPageResult[T any](list T, total int64, page, pageSize int) *PageResult[T] {
	return &PageResult[T]{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}
}

// Success 返回成功响应
func Success[T any](c *gin.Context, data T, msg ...string) {
	message := errorx.ErrSuccess.GetErrorMsg()
	if len(msg) > 0 {
		message = msg[0]
	}

	c.JSON(http.StatusOK, Response[T]{
		Code:    0, // 成功码为0
		Message: message,
		Data:    data,
	})
}

// SuccessNoData 返回无数据的成功响应
func SuccessNoData(c *gin.Context, msg ...string) {
	message := errorx.ErrSuccess.GetErrorMsg()
	if len(msg) > 0 {
		message = msg[0]
	}

	c.JSON(http.StatusOK, Response[struct{}]{
		Code:    0, // 成功码为0
		Message: message,
		Data:    struct{}{},
	})
}

// Error 返回错误响应
func Error(c *gin.Context, err error, msg ...string) {
	var (
		httpStatus int
		errorCode  int
		message    string
	)

	// 类型断言获取错误信息
	if appErr, ok := err.(*errorx.AppError); ok {
		message = appErr.GetErrorMsg()
		httpStatus = getHttpStatus(appErr)
		errorCode = appErr.GetErrorCode()
	} else {
		message = err.Error()
		httpStatus = http.StatusInternalServerError
		errorCode = errorx.ErrorUnknownCode
	}

	// 允许调用方覆盖原始错误消息
	if len(msg) > 0 {
		message = msg[0]
	}

	// 统一响应结构
	c.JSON(httpStatus, Response[struct{}]{
		Code:    errorCode,
		Message: message,
		Data:    struct{}{},
	})
}

// getHttpStatus 获取HTTP状态码，如果AppError没有设置HttpStatus则返回500
func getHttpStatus(err *errorx.AppError) int {
	if err.GetHttpStatus() == 0 {
		return http.StatusInternalServerError
	}
	return err.GetHttpStatus()
}
