package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/limitcool/starter/pkg/apiresponse"
)

// RequireSysUser 要求系统用户中间件
// 只允许用户类型为sys_user的用户访问
func RequireSysUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取令牌Claims
		claims, exists := c.Get("claims")
		if !exists {
			apiresponse.Unauthorized(c, "未授权访问")
			c.Abort()
			return
		}

		// 转换为MapClaims
		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			apiresponse.Unauthorized(c, "无效的令牌")
			c.Abort()
			return
		}

		// 获取用户类型
		userType, ok := mapClaims["user_type"].(string)
		if !ok || userType != "sys_user" {
			apiresponse.Forbidden(c, "访问被拒绝，需要系统用户权限")
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// RequireNormalUser 要求普通用户中间件
// 只允许用户类型为user的用户访问
func RequireNormalUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取令牌Claims
		claims, exists := c.Get("claims")
		if !exists {
			apiresponse.Unauthorized(c, "未授权访问")
			c.Abort()
			return
		}

		// 转换为MapClaims
		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			apiresponse.Unauthorized(c, "无效的令牌")
			c.Abort()
			return
		}

		// 获取用户类型
		userType, ok := mapClaims["user_type"].(string)
		if !ok || userType != "user" {
			apiresponse.Forbidden(c, "访问被拒绝，需要普通用户权限")
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}
