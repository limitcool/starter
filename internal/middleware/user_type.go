package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/errorx"
)

// AuthSystemUser 验证是否系统用户中间件
func AuthSystemUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("userType")
		if !exists {
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 验证是否系统用户
		if userType != "system" {
			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthNormalUser 验证是否普通用户中间件
func AuthNormalUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 验证是否普通用户
		if userType != "user" {
			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}
