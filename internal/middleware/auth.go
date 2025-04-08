package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/pkg/code"
	"github.com/limitcool/starter/pkg/jwt"
	"github.com/limitcool/starter/pkg/response"
)

func AuthMiddleware() gin.HandlerFunc {
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
			response.Unauthorized(c, code.GetMsg(code.UserAuthFailed))
			c.Abort()
			return
		}

		log.Debug("Authentication token received", "token", token)
		claims, err := jwt.ParseToken(token, global.Config.JwtAuth.AccessSecret)
		if err != nil {
			response.Unauthorized(c, code.GetMsg(code.UserAuthFailed)+":"+err.Error())
			c.Abort()
			return
		}

		log.Debug("Token claims parsed", "claims", claims)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), global.Token, claims))

		// 将用户ID存入请求上下文
		if userId, exists := (*claims)["user_id"]; exists {
			c.Set("userID", userId)
		}

		// 将token存入请求上下文
		c.Set("token", token)

		// 继续处理该请求
		c.Next()
	}
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    code.UserNoLogin,
		Message: message,
		Data:    nil,
	})
	// Unauthorized 方法不会自动调用 c.Abort()，需要在调用处手动添加
}
