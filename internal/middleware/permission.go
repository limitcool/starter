package middleware

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/storage/casbin"
)

// CasbinComponentMiddleware 基于路径和方法的权限控制中间件
func CasbinComponentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查权限系统是否启用
		if !core.Instance().Config().Permission.Enabled {
			// 权限系统未启用，直接放行
			c.Next()
			return
		}

		// 从上下文中获取用户ID
		userIDInterface, exists := c.Get("userID")
		if !exists {
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		// 将用户ID转换为字符串
		userID := strconv.FormatUint(uint64(userIDInterface.(float64)), 10)

		// 请求的路径
		obj := c.Request.URL.Path
		// 请求的方法
		act := c.Request.Method

		// 获取Casbin实例
		enforcer := casbin.GetEnforcer()
		if enforcer == nil {
			response.Error(c, errorx.ErrCasbinService)
			c.Abort()
			return
		}

		log.Debug("检查权限", "userID", userID, "object", obj, "action", act)

		// 检查权限
		pass, err := enforcer.Enforce(userID, obj, act)
		if err != nil {
			log.Error("权限检查错误", "error", err)
			response.Error(c, errorx.ErrCasbinService)
			c.Abort()
			return
		}

		if !pass {
			// 尝试获取用户角色
			roles, _ := enforcer.GetRolesForUser(userID)
			log.Debug("权限检查失败", "userID", userID, "roles", strings.Join(roles, ","))

			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		log.Debug("权限检查通过", "userID", userID)
		c.Next()
	}
}

// PermissionCodeMiddleware 基于权限编码的权限控制中间件
func PermissionCodeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查权限系统是否启用
		if !core.Instance().Config().Permission.Enabled {
			// 权限系统未启用，直接放行
			c.Next()
			return
		}

		// 获取需要的权限标识
		requiredPerm := c.GetHeader("X-Required-Permission")
		if requiredPerm == "" {
			// 如果没有设置所需权限，则默认通过
			c.Next()
			return
		}

		// 从上下文中获取用户ID
		userIDInterface, exists := c.Get("userID")
		if !exists {
			response.Error(c, errorx.ErrUserAuthFailed)
			c.Abort()
			return
		}

		userID := strconv.FormatUint(uint64(userIDInterface.(float64)), 10)

		// 获取Casbin实例
		enforcer := casbin.GetEnforcer()
		if enforcer == nil {
			response.Error(c, errorx.ErrCasbinService)
			c.Abort()
			return
		}

		// 获取用户角色
		roles, err := enforcer.GetRolesForUser(userID)
		if err != nil {
			log.Error("获取用户角色失败", "error", err)
			response.Error(c, errorx.ErrCasbinService)
			c.Abort()
			return
		}

		// 检查角色是否有所需权限
		hasPermission := false
		for _, role := range roles {
			// 检查是否为管理员
			if role == "admin" {
				hasPermission = true
				break
			}

			// 检查角色是否有权限
			pass, err := enforcer.Enforce(role, requiredPerm, "*")
			if err != nil {
				log.Error("权限检查错误", "error", err)
				continue
			}

			if pass {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}
