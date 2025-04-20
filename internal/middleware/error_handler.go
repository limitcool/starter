package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer+recover捕获所有可能的panic
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())

				// 获取请求ID和链路追踪ID
				requestID, _ := c.Get("request_id")
				traceID, _ := c.Get("trace_id")

				logger.Error("Panic recovered",
					"error", err,
					"stack", stack,
					"request_id", requestID,
					"trace_id", traceID)

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
	// 使用统一的错误响应函数
	response.Error(c, err)
}
