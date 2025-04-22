package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
