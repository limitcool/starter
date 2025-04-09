package middleware

import (
	"context"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/pkg/apiresponse"
	"github.com/limitcool/starter/pkg/code"
	"github.com/limitcool/starter/pkg/jwt"
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
			apiresponse.Unauthorized(c, code.GetMsg(code.UserAuthFailed))
			c.Abort()
			return
		}

		log.Debug("Authentication token received", "token", token)
		claims, err := jwt.ParseToken(token, global.Config.JwtAuth.AccessSecret)
		if err != nil {
			apiresponse.Unauthorized(c, code.GetMsg(code.UserAuthFailed)+":"+err.Error())
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
