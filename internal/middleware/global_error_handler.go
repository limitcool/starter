package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// GlobalErrorHandler 全局错误处理中间件
// 它会捕获所有通过 c.Error() 添加的错误，并统一处理
func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 检查是否已经有响应写入，避免重复响应
			if c.Writer.Written() {
				return
			}
			err := c.Errors.Last().Err
			handleError(c, err)
			c.Abort()
			return
		}

		// 检查响应状态码，如果是错误状态码，统一处理
		if c.Writer.Status() >= 400 && !c.Writer.Written() {
			handleError(c, fmt.Errorf("HTTP %d: %s", c.Writer.Status(), http.StatusText(c.Writer.Status())))
			c.Abort()
			return
		}
	}
}

// handleError 处理不同类型的错误
func handleError(c *gin.Context, err error) {
	// 获取请求上下文
	ctx := c.Request.Context()

	// 尝试使用错误链推导错误码
	err = errorx.WrapErrorWithContext(ctx, err, "")

	// 根据错误类型返回不同的响应
	var appErr *errorx.AppError
	switch e := err.(type) {
	case *errorx.AppError:
		appErr = e
	case error:
		// 如果是HTTP错误，转换为对应的应用错误
		if httpErr, ok := e.(*gin.Error); ok {
			switch httpErr.Type {
			case gin.ErrorTypeBind:
				appErr = errorx.ErrInvalidParams.WithError(e)
			case gin.ErrorTypeRender:
				appErr = errorx.ErrInternal.WithError(e)
			default:
				appErr = errorx.ErrInternal.WithError(e)
			}
		} else {
			appErr = errorx.ErrInternal.WithError(e)
		}
	default:
		appErr = errorx.ErrInternal.WithMsg(fmt.Sprintf("%v", err))
	}

	// 记录错误日志
	logger.ErrorContext(ctx, "Request error",
		"error", appErr,
		"path", c.Request.URL.Path,
		"method", c.Request.Method)

	// 使用统一的错误响应函数
	response.Error(c, appErr)
}

// ErrorHandlerFunc 是一个辅助函数，用于将返回 error 的控制器方法转换为 gin.HandlerFunc
// 使用方法: router.GET("/path", ErrorHandlerFunc(controller.Method))
func ErrorHandlerFunc(handler func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			// 如果已经有响应写入，不再处理错误
			if c.Writer.Written() {
				return
			}

			// 将错误添加到 gin 的错误链中，由 GlobalErrorHandler 统一处理
			_ = c.Error(err)
		}
	}
}
