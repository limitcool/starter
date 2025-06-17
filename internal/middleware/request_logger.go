package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// RequestLoggerMiddleware 是一个记录请求日志的中间件，同时处理请求ID和链路追踪ID
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 处理请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 处理链路追踪ID
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = fmt.Sprintf("trace-%d", time.Now().UnixNano())
		}
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)

		// 将请求ID和链路追踪ID添加到context.Context中
		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		ctx = context.WithValue(ctx, "trace_id", traceID)
		c.Request = c.Request.WithContext(ctx)

		// 处理请求
		c.Next()

		// 获取请求上下文
		reqCtx := c.Request.Context()

		// 计算延迟
		latency := time.Since(start)

		// 根据状态码选择日志级别
		status := c.Writer.Status()

		// 准备日志字段
		fields := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"request_id", requestID,
			"user_agent", c.Request.UserAgent(),
			"referer", c.Request.Referer(),
			"body_size", c.Writer.Size(),
		}

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			fields = append(fields, "errors", c.Errors.String())
		}

		// 根据状态码选择日志级别
		if status >= 500 {
			logger.ErrorContext(reqCtx, "Server error", fields...)
		} else if status >= 400 {
			logger.WarnContext(reqCtx, "Client error", fields...)
		} else {
			logger.InfoContext(reqCtx, "Request completed", fields...)
		}
	}
}
