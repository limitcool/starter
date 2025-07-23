package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/errorx"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// UserCheck 用户检查中间件 - 确保用户已登录
func UserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !CheckUserLogin(c) {
			return
		}
		c.Next()
	}
}

// UserCheckWithDB 用户检查中间件 - 从数据库验证用户身份
// 适用于需要确保用户在数据库中仍然存在的场景
func UserCheckWithDB(userRepo *model.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求上下文
		ctx := c.Request.Context()

		// 从上下文获取用户ID
		userID := GetUserIDInt64(c)
		if userID == 0 {
			logger.WarnContext(ctx, "用户ID不存在")
			response.Error(c, errorx.ErrUserNotLogin.New(ctx))
			c.Abort()
			return
		}

		// 获取用户
		_, err := userRepo.GetByID(ctx, userID)
		if err != nil {
			logger.ErrorContext(ctx, "获取用户失败", "error", err, "user_id", userID)
			response.Error(c, errorx.ErrUserNotFound.New(ctx))
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// RegularUserCheck 普通用户检查中间件 - 确保用户不是管理员
// 适用于只允许普通用户访问的接口
func RegularUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先检查是否已登录
		if !CheckUserLogin(c) {
			return
		}

		// 检查用户是否为管理员
		isAdmin, ok := c.Get("is_admin")
		if ok && isAdmin.(bool) {
			ctx := c.Request.Context()
			logger.WarnContext(ctx, "管理员不能访问普通用户接口")
			response.Error(c, errorx.ErrAccessDenied.New(ctx))
			c.Abort()
			return
		}

		c.Next()
	}
}
