package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
)

// SimpleAdminCheck 简单模式下的管理员检查中间件
// 用于检查用户是否是管理员
func SimpleAdminCheck(userRepo *repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求上下文
		ctx := c.Request.Context()

		// 从上下文中获取用户ID
		userID := GetUserID(c)
		if userID == 0 {
			logger.WarnContext(ctx, "未登录用户尝试访问管理员接口")
			response.Error(c, errorx.ErrUserNoLogin)
			c.Abort()
			return
		}

		// 获取用户类型
		userType, exists := c.Get("user_type") // 使用统一的键名 user_type
		if !exists {
			logger.WarnContext(ctx, "用户类型不存在", "user_id", userID)
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 安全地检查用户类型是否是管理员
		userTypeStr, ok := userType.(string)
		if !ok {
			logger.ErrorContext(ctx, "用户类型格式错误", "user_id", userID, "user_type", userType)
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		if userTypeStr != enum.UserTypeAdminUser.String() {
			logger.WarnContext(ctx, "非管理员用户尝试访问管理员接口", "user_id", userID, "user_type", userTypeStr)
			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		// 获取用户信息
		user, err := userRepo.GetByID(ctx, int64(userID))
		if err != nil {
			logger.ErrorContext(ctx, "获取用户信息失败", "error", err, "user_id", userID)
			response.Error(c, errorx.ErrUserNotFound)
			c.Abort()
			return
		}

		// 检查用户是否是管理员
		if !user.IsAdmin {
			logger.WarnContext(ctx, "非管理员用户尝试访问管理员接口", "user_id", userID)
			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}
