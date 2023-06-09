package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/common/jwtx"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/pkg/code"
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

		// 如果没有token,返回错误
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, code.ApiResponse{
				ErrorCode: code.UserAuthFailed,
				Message:   code.GetMsg(code.UserAuthFailed),
				Data:      nil,
			})
			return
		}
		fmt.Printf("token: %v\n", token)
		claims, err := jwtx.ParseToken(token, global.Config.JwtAuth.AccessSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, code.ApiResponse{
				ErrorCode: code.UserAuthFailed,
				Message:   code.GetMsg(code.UserAuthFailed) + ":" + err.Error(),
				Data:      nil,
			})
		}
		fmt.Printf("claims: %v\n", claims)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), global.Token, claims))
		// 将token存入请求上下文
		c.Set("token", token)

		// 继续处理该请求
		c.Next()
	}
}
