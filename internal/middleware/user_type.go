package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// AuthSystemUser 验证是否系统用户中间件
func AuthSystemUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求上下文
		ctx := c.Request.Context()

		// 从上下文中获取用户ID
		userID := GetUserID(c)

		// 获取用户类型
		userType, exists := c.Get("user_type") // 使用统一的键名 user_type
		if !exists {
			logger.WarnContext(ctx, "用户类型不存在", "user_id", userID)
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 安全地检查用户类型
		userTypeStr, ok := userType.(string)
		if !ok {
			logger.ErrorContext(ctx, "用户类型格式错误", "user_id", userID, "user_type", userType)
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 验证是否系统用户
		if userTypeStr != enum.UserTypeSysUser.String() {
			logger.WarnContext(ctx, "非系统用户尝试访问系统接口", "user_id", userID, "user_type", userTypeStr)
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
		// 获取请求上下文
		ctx := c.Request.Context()

		// 从上下文中获取用户ID
		userID := GetUserID(c)

		// 获取用户类型
		userType, exists := c.Get("user_type")
		if !exists {
			logger.WarnContext(ctx, "用户类型不存在", "user_id", userID)
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 安全地检查用户类型
		userTypeStr, ok := userType.(string)
		if !ok {
			logger.ErrorContext(ctx, "用户类型格式错误", "user_id", userID, "user_type", userType)
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 验证是否普通用户
		if userTypeStr != enum.UserTypeUser.String() {
			logger.WarnContext(ctx, "非普通用户尝试访问普通用户接口", "user_id", userID, "user_type", userTypeStr)
			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}
