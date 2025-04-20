package middleware

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/jwt"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// 上下文键类型
type contextKey string

// 上下文键常量
const (
	TokenKey contextKey = "token"
)

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) uint64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	// 尝试转换为uint64
	switch v := userID.(type) {
	case float64:
		return uint64(v)
	case int:
		return uint64(v)
	case int64:
		return uint64(v)
	case uint:
		return uint64(v)
	case uint64:
		return v
	case string:
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0
		}
		return id
	default:
		return 0
	}
}

// JWTAuth JWT认证中间件
func JWTAuth(config *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authorization := c.GetHeader("Authorization")

		// 检查前缀并提取 token
		token := ""
		if strings.HasPrefix(authorization, "Bearer ") {
			token = strings.Split(authorization, " ")[1]
		}

		// 如果没有token,返回错误并中止
		if token == "" {
			ctx := c.Request.Context()
			logger.ErrorContext(ctx, "No authentication token provided")
			response.Error(c, errorx.ErrUserNoLogin)
			c.Abort()
			return
		}

		// 解析token
		ctx := c.Request.Context()
		claims, err := jwt.ParseTokenWithContext(ctx, token, config.JwtAuth.AccessSecret)
		if err != nil {
			logger.ErrorContext(ctx, "Authentication token parse failed", "error", err)
			response.Error(c, errorx.ErrUserTokenError)
			c.Abort()
			return
		}

		// 将claims存入请求上下文
		ctx = context.WithValue(ctx, TokenKey, claims)

		// 将用户ID存入请求上下文
		if userId, exists := (*claims)["user_id"]; exists {
			c.Set("user_id", userId)
			ctx = context.WithValue(ctx, "user_id", userId)
		}
		if userType, exists := (*claims)["user_type"]; exists {
			c.Set("user_type", userType)
			ctx = context.WithValue(ctx, "user_type", userType)
		}
		// 将token存入请求上下文
		c.Set("token", token)
		ctx = context.WithValue(ctx, "token", token)

		// 更新请求上下文
		c.Request = c.Request.WithContext(ctx)

		// 添加用户信息到上下文
		// TODO: 在此处获取用户/系统用户信息并添加到上下文中
		// 这里需要调用 userService 或 sysUserService 来获取用户信息

		// 继续处理该请求
		c.Next()
	}
}
