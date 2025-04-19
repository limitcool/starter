// Package response 提供API响应相关的功能
package response

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/errorx"
)

// Response API标准响应结构
type Response[T any] struct {
	Code    int    `json:"code"`                 // 错误码
	Msg     string `json:"message"`              // 提示信息
	Data    T      `json:"data"`                 // 数据
	ReqID   string `json:"request_id,omitempty"` // 请求ID
	Time    int64  `json:"timestamp,omitempty"`  // 时间戳
	TraceID string `json:"trace_id,omitempty"`   // 链路追踪ID
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
	message := "success"
	if len(msg) > 0 {
		message = msg[0]
	}

	// 获取请求ID
	requestID := getRequestID(c)

	c.JSON(http.StatusOK, Response[T]{
		Code:  0, // 成功码为0
		Msg:   message,
		Data:  data,
		ReqID: requestID,
		Time:  time.Now().Unix(),
	})
}

// SuccessNoData 返回无数据的成功响应
func SuccessNoData(c *gin.Context, msg ...string) {
	message := "success"
	if len(msg) > 0 {
		message = msg[0]
	}

	// 获取请求ID
	requestID := getRequestID(c)

	c.JSON(http.StatusOK, Response[struct{}]{
		Code:  0, // 成功码为0
		Msg:   message,
		Data:  struct{}{},
		ReqID: requestID,
		Time:  time.Now().Unix(),
	})
}

// Error 返回错误响应
func Error(c *gin.Context, err error, msg ...string) {
	var (
		httpStatus int
		errorCode  int
		message    string
	)

	// 尝试使用错误链推导错误码
	err = errorx.WrapError(err, "")

	// 获取错误信息
	if appErr, ok := err.(*errorx.AppError); ok {
		// 如果是 AppError类型，直接使用其属性
		message = appErr.GetErrorMsg()
		httpStatus = getHttpStatus(appErr)
		errorCode = appErr.GetErrorCode()
	} else {
		// 如果不是 AppError类型，使用默认值
		message = err.Error()
		httpStatus = http.StatusInternalServerError
		errorCode = errorx.ErrorUnknownCode
	}

	// 允许调用方覆盖原始错误消息
	if len(msg) > 0 {
		message = msg[0]
	}

	// 获取请求ID
	requestID := getRequestID(c)

	// 获取链路追踪ID
	traceID := getTraceIDFromContext(c)

	// 记录错误到日志
	log.Error("API error occurred",
		"code", errorCode,
		"message", message,
		"trace_id", traceID,
		"request_id", requestID,
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"client_ip", c.ClientIP(),
		"error_chain", errorx.FormatErrorChain(err),
	)

	// 统一响应结构
	c.JSON(httpStatus, Response[struct{}]{
		Code:    errorCode,
		Msg:     message,
		Data:    struct{}{},
		ReqID:   requestID,
		Time:    time.Now().Unix(),
		TraceID: traceID,
	})
}

// getHttpStatus 获取HTTP状态码，如果AppError没有设置HttpStatus则返回500
func getHttpStatus(err *errorx.AppError) int {
	if err.GetHttpStatus() == 0 {
		return http.StatusInternalServerError
	}
	return err.GetHttpStatus()
}

// getRequestID 获取请求ID，如果不存在则生成新的
func getRequestID(c *gin.Context) string {
	// 先从请求头部获取
	reqID := c.GetHeader("X-Request-ID")
	if reqID != "" {
		return reqID
	}

	// 如果上下文中已经有请求ID，则使用它
	if id, exists := c.Get("request_id"); exists {
		if strID, ok := id.(string); ok && strID != "" {
			return strID
		}
	}

	// 生成新的请求ID
	newID := time.Now().UnixNano()
	// 将请求ID存储到上下文中
	c.Set("request_id", newID)
	return time.Now().Format("20060102150405") + "-" + c.ClientIP()
}

// getTraceIDFromContext 从上下文中获取链路追踪ID
func getTraceIDFromContext(c *gin.Context) string {
	// 先从上下文中获取
	if traceID, exists := c.Get("trace_id"); exists {
		if strID, ok := traceID.(string); ok && strID != "" {
			return strID
		}
	}

	// 如果上下文中没有，尝试从请求头中获取
	traceID := c.GetHeader("X-Trace-ID")
	if traceID != "" {
		return traceID
	}

	return ""
}
