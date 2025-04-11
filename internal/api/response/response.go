// Package response 提供API响应相关的功能
// 此包整合了以前 internal/pkg/apiresponse 包的功能，并增加了分页支持
package response

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/i18n"
	"github.com/limitcool/starter/internal/pkg/logger"
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
		i18nKey    string
	)

	// 检查是否需要显示堆栈
	showStackTrace := logger.ShouldShowStackTrace(log.DebugLevel)

	// 类型断言获取错误信息
	if appErr, ok := err.(*errorx.AppError); ok {
		message = appErr.GetErrorMsg()
		httpStatus = getHttpStatus(appErr)
		errorCode = appErr.GetErrorCode()
		i18nKey = appErr.GetI18nKey()

		// 根据配置决定是否记录堆栈
		if showStackTrace {
			// 在开发环境中记录完整错误（包含堆栈）
			log.Debug("API错误详情", "err", fmt.Sprintf("%+v", appErr))
		} else {
			// 不记录堆栈，只记录基本错误信息
			log.Debug("API错误详情",
				"err_code", errorCode,
				"err_msg", message,
				"i18n_key", i18nKey)
		}
	} else {
		message = err.Error()
		httpStatus = http.StatusInternalServerError
		errorCode = errorx.ErrorUnknownCode
		i18nKey = "error.unknown"

		// 记录原始错误，对于非AppError类型
		logFields := []interface{}{"err_code", errorCode, "err_msg", message, "i18n_key", i18nKey}

		// 根据配置决定是否记录堆栈
		if showStackTrace {
			// 尝试使用Formatter接口获取堆栈
			if formatter, ok := err.(fmt.Formatter); ok {
				var buf bytes.Buffer
				fmt.Fprintf(&buf, "%+v", formatter)
				logFields = append(logFields, "stack", "\n"+buf.String())
			}
		}

		log.Debug("API错误详情", logFields...)
	}

	// 允许调用方覆盖原始错误消息
	if len(msg) > 0 {
		message = msg[0]
	} else {
		// 尝试获取国际化的消息
		lang, exists := c.Get("lang")
		if exists && i18nKey != "" {
			// 使用i18n翻译消息
			langStr := lang.(string)
			translatedMsg := i18n.T(i18nKey, langStr)
			if translatedMsg != i18nKey { // 确保翻译成功
				message = translatedMsg
			}
		}
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
