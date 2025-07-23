package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// PanicRecovery 中间件用于捕获 panic 并返回友好的错误响应
// 这个中间件只处理 panic，其他错误由 GlobalErrorHandler 处理
func PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer+recover捕获所有可能的panic
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())

				// 获取请求上下文
				ctx := c.Request.Context()

				// 获取请求ID和链路追踪ID
				requestID, _ := c.Get("request_id")
				traceID, _ := c.Get("trace_id")

				logger.ErrorContext(ctx, "Panic recovered",
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
					appErr = errorx.ErrInternal.New(ctx, errorx.None).Wrap(e)
				case string:
					appErr = errorx.ErrInternal.New(ctx, errorx.None).WithMessage(e)
				default:
					appErr = errorx.ErrInternal.New(ctx, errorx.None).WithMessage(fmt.Sprintf("%v", err))
				}

				// 检查是否已经有响应写入，避免重复响应
				if !c.Writer.Written() {
					// 返回错误响应
					response.Error(c, appErr)
				}
				c.Abort()
			}
		}()

		// 处理请求
		c.Next()
	}
}
