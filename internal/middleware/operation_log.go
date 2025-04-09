package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/limitcool/starter/internal/services"
)

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware(module, action, description string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 记录请求体，仅适用于POST/PUT/PATCH请求
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			// 恢复Body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// 保存请求体数据供后续使用
			c.Set("requestBody", string(bodyBytes))
		}

		// 处理请求
		c.Next()

		// 获取用户信息
		claims, exists := c.Get("claims")
		if !exists {
			// 未登录用户，不记录操作日志
			return
		}

		// 创建操作日志服务
		logService := services.NewOperationLogService()

		// 提取用户信息
		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			return
		}

		// 提取用户ID、用户名和用户类型
		var userID uint
		var username string
		var userType string

		if uid, ok := mapClaims["user_id"].(float64); ok {
			userID = uint(uid)
		}
		if un, ok := mapClaims["username"].(string); ok {
			username = un
		}
		if ut, ok := mapClaims["user_type"].(string); ok {
			userType = ut
		} else {
			userType = "sys_user" // 默认为系统用户
		}

		// 根据用户类型记录不同的日志
		if userType == "sys_user" {
			logService.CreateSysUserLog(c, userID, username, module, action, description, startTime)
		} else {
			logService.CreateUserLog(c, userID, username, module, action, description, startTime)
		}
	}
}
