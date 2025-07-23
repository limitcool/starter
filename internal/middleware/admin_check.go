package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/errorx"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// AdminCheck 管理员检查中间件 - 基于JWT中的is_admin字段
func AdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !CheckAdminPermission(c) {
			return
		}
		c.Next()
	}
}

// AdminCheckWithDB 管理员检查中间件 - 从数据库验证管理员身份
// 适用于需要确保用户在数据库中仍然是管理员的场景
func AdminCheckWithDB(userRepo *model.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求上下文
		ctx := c.Request.Context()

		// 从上下文中获取用户ID
		userID := GetUserIDInt64(c)
		if userID == 0 {
			logger.WarnContext(ctx, "未登录用户尝试访问管理员接口")
			response.Error(c, errorx.ErrUserNotLogin.New(ctx))
			c.Abort()
			return
		}

		// 获取用户信息
		user, err := userRepo.GetByID(ctx, userID)
		if err != nil {
			logger.ErrorContext(ctx, "获取用户信息失败", "error", err, "user_id", userID)
			response.Error(c, errorx.ErrUserNotFound.New(ctx))
			c.Abort()
			return
		}

		// 检查用户是否是管理员
		if !user.IsAdmin {
			logger.WarnContext(ctx, "非管理员用户尝试访问管理员接口", "user_id", userID)
			response.Error(c, errorx.ErrAccessDenied.New(ctx))
			c.Abort()
			return
		}

		c.Next()
	}
}
