package middleware

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()

		// 执行时间
		latency := end.Sub(start)

		// 请求方法
		method := c.Request.Method

		// 请求路由
		path := c.Request.URL.Path

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Info("HTTP Request",
			"method", method,
			"path", path,
			"status", statusCode,
			"latency", latency,
			"ip", clientIP,
		)
	}
}

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Casbin 权限中间件
func Casbin(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户
		// 这里需要根据实际情况获取用户信息
		user := "anonymous"

		// 获取请求方法和路径
		method := c.Request.Method
		path := c.Request.URL.Path

		// 检查权限
		if enforcer != nil {
			// 检查用户是否有权限访问
			if ok, err := enforcer.Enforce(user, path, method); err != nil {
				logger.Error("Casbin enforce error", "error", err)
				c.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "Internal server error",
				})
				return
			} else if !ok {
				// 没有权限
				c.AbortWithStatusJSON(403, gin.H{
					"code":    403,
					"message": "Forbidden",
				})
				return
			}
		}

		c.Next()
	}
}
