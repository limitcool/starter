package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	// 尝试转换为float64
	if id, ok := userID.(float64); ok {
		return uint(id)
	}

	// 尝试转换为uint
	if id, ok := userID.(uint); ok {
		return id
	}

	// 尝试转换为int64
	if id, ok := userID.(int64); ok {
		return uint(id)
	}

	return 0
}

// GetUserIDInt64 从上下文中获取用户ID（int64类型）
func GetUserIDInt64(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	// 尝试转换为float64
	if id, ok := userID.(float64); ok {
		return int64(id)
	}

	// 尝试转换为uint
	if id, ok := userID.(uint); ok {
		return int64(id)
	}

	// 尝试转换为int64
	if id, ok := userID.(int64); ok {
		return id
	}

	return 0
}

// GetUserIDString 从上下文中获取用户ID（字符串类型）
func GetUserIDString(c *gin.Context) string {
	id := GetUserID(c)
	if id == 0 {
		return ""
	}
	return fmt.Sprintf("%d", id)
}

// CheckUserLogin 检查用户是否已登录，如果未登录则返回错误响应
func CheckUserLogin(c *gin.Context) bool {
	ctx := c.Request.Context()

	_, exists := c.Get("user_id")
	if !exists {
		logger.WarnContext(ctx, "用户ID不存在")
		response.Error(c, errorx.ErrUserNotFound.New(ctx, errorx.None))
		c.Abort()
		return false
	}

	return true
}

// CheckAdminPermission 检查用户是否为管理员，如果不是则返回错误响应
func CheckAdminPermission(c *gin.Context) bool {
	ctx := c.Request.Context()

	// 先检查是否已登录
	if !CheckUserLogin(c) {
		return false
	}

	// 检查用户是否为管理员
	isAdmin, ok := c.Get("is_admin")
	if !ok || !isAdmin.(bool) {
		logger.WarnContext(ctx, "用户不是管理员", "is_admin", isAdmin)
		response.Error(c, errorx.ErrForbidden.New(ctx, errorx.None))
		c.Abort()
		return false
	}

	return true
}
