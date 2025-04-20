package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// ErrorLogger 是一个记录错误日志的中间件
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last().Err

			// 记录请求信息
			ctx := c.Request.Context()
			logger.ErrorContext(ctx, "API error occurred",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
				"latency_ms", time.Since(start).Milliseconds(),
				"status", c.Writer.Status(),
				"error", err.Error(),
			)

			// 返回错误响应
			response.Error(c, err)

			// 中止后续处理
			c.Abort()
		}
	}
}

// RequestLogger 是一个记录请求日志的中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 获取请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Request.Header.Set("X-Request-ID", requestID)
			c.Set("request_id", requestID) // 同时存入上下文

			// 将请求ID添加到context.Context中
			ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
			c.Request = c.Request.WithContext(ctx)
		}

		// 处理请求
		c.Next()

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
		ctx := c.Request.Context()
		if status >= 500 {
			logger.ErrorContext(ctx, "Server error", fields...)
		} else if status >= 400 {
			logger.WarnContext(ctx, "Client error", fields...)
		} else {
			logger.InfoContext(ctx, "Request completed", fields...)
		}
	}
}
