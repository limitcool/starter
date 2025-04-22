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
		// 从上下文中获取用户ID
		userID := GetUserID(c)
		if userID == 0 {
			response.Error(c, errorx.ErrUserNoLogin)
			c.Abort()
			return
		}

		// 获取用户类型
		userType, exists := c.Get("userType")
		if !exists {
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 检查用户类型是否是管理员
		if userType.(string) != enum.UserTypeAdminUser.String() {
			logger.Warn("非管理员用户尝试访问管理员接口", "userID", userID, "userType", userType)
			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		// 获取用户信息
		user, err := userRepo.GetByID(c.Request.Context(), int64(userID))
		if err != nil {
			logger.Error("获取用户信息失败", "error", err)
			response.Error(c, errorx.ErrUserNotFound)
			c.Abort()
			return
		}

		// 检查用户是否是管理员
		if !user.IsAdmin {
			logger.Warn("非管理员用户尝试访问管理员接口", "userID", userID)
			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}
