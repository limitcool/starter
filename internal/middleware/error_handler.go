package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
)

// ErrorHandlerParams 错误处理中间件参数
type ErrorHandlerParams struct {
	// 错误监控服务，可选
	ErrorMonitorService *services.ErrorMonitorService
}

// ErrorHandler 全局错误处理中间件
// 统一处理所有控制器、服务和仓库层上抛的错误
func ErrorHandler(params ...ErrorHandlerParams) gin.HandlerFunc {
	var errorMonitorService *services.ErrorMonitorService
	if len(params) > 0 && params[0].ErrorMonitorService != nil {
		errorMonitorService = params[0].ErrorMonitorService
	}

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

				// 记录详细的panic日志
				logger.ErrorContext(ctx, "Panic recovered",
					"error", err,
					"stack", stack,
					"request_id", requestID,
					"trace_id", traceID,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"client_ip", c.ClientIP())

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

				// 记录错误到监控服务
				if errorMonitorService != nil {
					errorMonitorService.RecordError(ctx, appErr)
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
			handleError(c, err, errorMonitorService)
			c.Abort()
		}
	}
}

// handleError 处理不同类型的错误
func handleError(c *gin.Context, err error, errorMonitorService *services.ErrorMonitorService) {
	// 获取请求上下文
	ctx := c.Request.Context()

	// 记录错误到监控服务
	if errorMonitorService != nil && errorx.IsAppErr(err) {
		errorMonitorService.RecordError(ctx, err)
	}

	// 使用统一的错误响应函数
	// response.Error 内部会记录错误日志，所以这里不需要重复记录
	response.Error(c, err)
}

// ErrorAware 错误感知中间件
// 用于在控制器方法中使用 c.Error(err) 上抛错误
func ErrorAware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
