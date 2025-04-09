package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/code"
	"github.com/limitcool/starter/internal/services"
	"gorm.io/gorm"
)

// CasbinMiddleware Casbin权限控制中间件
func CasbinMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的URI
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := c.GetString("role")

		// 创建Casbin服务
		casbinService := services.NewCasbinService(db)
		if casbinService == nil {
			response.ServerError(c)
			c.Abort()
			return
		}

		// 检查权限
		ok, err := casbinService.CheckPermission(sub, obj, act)
		if err != nil {
			response.ServerError(c)
			c.Abort()
			return
		}

		if !ok {
			response.Forbidden(c, code.GetMsg(code.AccessDenied))
			c.Abort()
			return
		}

		c.Next()
	}
}

// PermissionMiddleware 基于菜单权限标识的权限控制中间件
func PermissionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的URI
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := c.GetString("role")

		// 创建Casbin服务
		casbinService := services.NewCasbinService(db)
		if casbinService == nil {
			response.ServerError(c)
			c.Abort()
			return
		}

		// 检查权限
		ok, err := casbinService.CheckPermission(sub, obj, act)
		if err != nil {
			response.ServerError(c)
			c.Abort()
			return
		}

		if !ok {
			response.Forbidden(c, code.GetMsg(code.AccessDenied))
			c.Abort()
			return
		}

		c.Next()
	}
}
