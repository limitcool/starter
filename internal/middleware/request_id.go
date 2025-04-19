package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestID 中间件，用于生成和传递请求ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取请求ID
		requestID := c.GetHeader("X-Request-ID")
		
		// 如果请求头中没有请求ID，则生成一个新的
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		
		// 将请求ID存储到上下文中
		c.Set("request_id", requestID)
		
		// 将请求ID添加到响应头中
		c.Header("X-Request-ID", requestID)
		
		c.Next()
	}
}

// TraceID 中间件，用于生成和传递链路追踪ID
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取链路追踪ID
		traceID := c.GetHeader("X-Trace-ID")
		
		// 如果请求头中没有链路追踪ID，则生成一个新的
		if traceID == "" {
			traceID = fmt.Sprintf("trace-%d", time.Now().UnixNano())
		}
		
		// 将链路追踪ID存储到上下文中
		c.Set("trace_id", traceID)
		
		// 将链路追踪ID添加到响应头中
		c.Header("X-Trace-ID", traceID)
		
		c.Next()
	}
}

// RequestContext 中间件，同时处理请求ID和链路追踪ID
func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		
		c.Next()
	}
}
