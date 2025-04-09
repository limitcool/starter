package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
)

// AuthSystemUser 验证是否系统用户中间件
func AuthSystemUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("userType")
		if !exists {
			response.Unauthorized(c, "未授权访问")
			c.Abort()
			return
		}

		// 验证是否系统用户
		if userType != "system" {
			response.Forbidden(c, "访问被拒绝，需要系统用户权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthNormalUser 验证是否普通用户中间件
func AuthNormalUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("userType")
		if !exists {
			response.Unauthorized(c, "未授权访问")
			c.Abort()
			return
		}

		// 验证是否普通用户
		if userType != "user" {
			response.Forbidden(c, "访问被拒绝，需要普通用户权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
