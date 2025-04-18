package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/i18n"
	"github.com/limitcool/starter/internal/services"
	"gorm.io/gorm"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer+recover捕获所有可能的panic
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())
				log.Error("Panic recovered", "error", err, "stack", stack)

				// 根据不同类型的panic返回不同的错误
				var appErr *errorx.AppError
				switch e := err.(type) {
				case *errorx.AppError:
					appErr = e
				case error:
					appErr = errorx.ErrInternal.WithError(e)
				case string:
					appErr = errorx.ErrInternal.WithMsg(e)
				default:
					appErr = errorx.ErrInternal.WithMsg(fmt.Sprintf("%v", err))
				}

				// 返回错误响应
				response.Error(c, appErr)
				c.Abort()
			}
		}()

		// 处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
			c.Abort()
		}
	}
}

// handleError 处理不同类型的错误
func handleError(c *gin.Context, err error) {
	// 获取当前语言
	lang := getCurrentLanguage(c)

	// 处理不同类型的错误
	switch {
	case errorx.IsAppErr(err):
		// 已经是AppError类型，直接使用
		appErr := err.(*errorx.AppError)

		// 尝试翻译错误消息
		message := i18n.T(appErr.GetI18nKey(), lang)
		if message == appErr.GetI18nKey() {
			// 如果翻译失败，使用原始错误消息
			message = appErr.GetErrorMsg()
		}

		// 记录错误日志
		logAppError(appErr)

		// 返回错误响应
		c.JSON(appErr.GetHttpStatus(), gin.H{
			"code":    appErr.GetErrorCode(),
			"message": message,
			"data":    struct{}{},
		})

	case errors.Is(err, gorm.ErrRecordNotFound):
		// 数据库记录未找到错误
		appErr := errorx.ErrNotFound.WithError(err)
		logAppError(appErr)

		message := i18n.T(appErr.GetI18nKey(), lang)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    appErr.GetErrorCode(),
			"message": message,
			"data":    struct{}{},
		})

	case isGormError(err):
		// 其他GORM错误
		appErr := errorx.ErrDatabase.WithError(err)
		logAppError(appErr)

		message := i18n.T(appErr.GetI18nKey(), lang)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    appErr.GetErrorCode(),
			"message": message,
			"data":    struct{}{},
		})

	default:
		// 未知错误
		appErr := errorx.ErrUnknown.WithError(err)
		logAppError(appErr)

		message := i18n.T(appErr.GetI18nKey(), lang)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    appErr.GetErrorCode(),
			"message": message,
			"data":    struct{}{},
		})
	}
}

// getCurrentLanguage 获取当前语言
func getCurrentLanguage(c *gin.Context) string {
	lang, exists := c.Get("lang")
	if !exists {
		return i18n.GetDefaultLanguage()
	}
	return lang.(string)
}

// logAppError 记录AppError类型的错误
func logAppError(err *errorx.AppError) {
	// 记录错误详情
	log.Error("Application error",
		"code", err.GetErrorCode(),
		"message", err.GetErrorMsg(),
		"i18n_key", err.GetI18nKey(),
		"stack", err.GetStackTrace(),
	)

	// 记录错误到错误监控服务
	services.NewErrorMonitorService().RecordError(err)
}

// isGormError 判断是否为GORM错误
func isGormError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "gorm") ||
		strings.Contains(errMsg, "sql") ||
		strings.Contains(errMsg, "database") ||
		strings.Contains(errMsg, "constraint")
}
