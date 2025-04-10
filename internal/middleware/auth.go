package middleware

import (
	"context"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/jwt"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
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
			log.Error("No authentication token provided")
			response.Error(c, errorx.ErrUserNoLogin)
			c.Abort()
			return
		}

		// 解析token
		claims, err := jwt.ParseToken(token, core.Instance().Config().JwtAuth.AccessSecret)
		if err != nil {
			log.Error("Authentication token parse failed", "error", err)
			response.Error(c, errorx.ErrUserTokenError)
			c.Abort()
			return
		}

		// 将claims存入请求上下文
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), global.Token, claims))

		// 将用户ID存入请求上下文
		if userId, exists := (*claims)["user_id"]; exists {
			c.Set("user_id", userId)
		}
		if userType, exists := (*claims)["user_type"]; exists {
			c.Set("user_type", userType)
		}
		// 将token存入请求上下文
		c.Set("token", token)

		// 添加用户信息到上下文
		// TODO: 在此处获取用户/系统用户信息并添加到上下文中
		// 这里需要调用 userService 或 sysUserService 来获取用户信息

		// 继续处理该请求
		c.Next()
	}
}
