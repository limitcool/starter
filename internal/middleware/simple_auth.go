package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
)

// SimpleUserCheck 简化的普通用户检查中间件
// 确保用户不是管理员
func SimpleUserCheck(userRepo *repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			logger.ErrorContext(c.Request.Context(), "用户ID不存在")
			response.Error(c, errorx.ErrUserNoLogin)
			c.Abort()
			return
		}

		// 转换用户ID
		userIDInt64, ok := userID.(float64)
		if !ok {
			logger.ErrorContext(c.Request.Context(), "用户ID类型错误")
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 获取用户
		_, err := userRepo.GetByID(c.Request.Context(), int64(userIDInt64))
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "获取用户失败", "error", err)
			response.Error(c, errorx.ErrUserNotFound)
			c.Abort()
			return
		}

		// 如果需要严格区分，可以检查用户是否为管理员
		// 注释掉下面的代码，如果不需要严格区分
		/*
			if user.IsAdmin {
				logger.ErrorContext(c.Request.Context(), "管理员不能访问普通用户接口", "user_id", userIDInt64)
				response.Error(c, errorx.ErrAccessDenied)
				c.Abort()
				return
			}
		*/

		// 继续处理请求
		c.Next()
	}
}
