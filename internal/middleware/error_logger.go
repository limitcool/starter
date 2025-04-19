package middleware

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
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
			log.Error("API error occurred",
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
			requestID = time.Now().Format("20060102150405") + "-" + c.ClientIP()
			c.Request.Header.Set("X-Request-ID", requestID)
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
		if status >= 500 {
			log.Error("Server error", fields...)
		} else if status >= 400 {
			log.Warn("Client error", fields...)
		} else {
			log.Info("Request completed", fields...)
		}
	}
}
