package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

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

		// 将请求ID和链路追踪ID添加到context.Context中
		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		ctx = context.WithValue(ctx, "trace_id", traceID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
