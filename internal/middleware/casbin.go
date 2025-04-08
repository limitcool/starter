package middleware

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/pkg/code"
	"github.com/limitcool/starter/pkg/response"
	"gorm.io/gorm"
)

// CasbinMiddleware Casbin权限控制中间件
func CasbinMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userIDInterface, exists := c.Get("userID")
		if !exists {
			response.Unauthorized(c, code.GetMsg(code.UserAuthFailed))
			c.Abort()
			return
		}

		// 将用户ID转换为字符串
		userID := strconv.FormatUint(uint64(userIDInterface.(float64)), 10)

		// 请求的路径
		obj := c.Request.URL.Path
		// 请求的方法
		act := c.Request.Method

		// 创建Casbin服务
		casbinService := services.NewCasbinService(db)
		if casbinService == nil {
			response.InternalServerError(c, "Casbin服务初始化失败")
			c.Abort()
			return
		}

		log.Debug("检查权限", "userID", userID, "object", obj, "action", act)

		// 检查权限
		pass, err := casbinService.CheckPermission(userID, obj, act)
		if err != nil {
			log.Error("权限检查错误", "error", err)
			response.InternalServerError(c, "权限检查失败")
			c.Abort()
			return
		}

		if !pass {
			// 尝试获取用户角色
			roles, _ := casbinService.GetRolesForUser(userID)
			log.Debug("权限检查失败", "userID", userID, "roles", strings.Join(roles, ","))

			response.Forbidden(c, code.GetMsg(code.AccessDenied))
			c.Abort()
			return
		}

		log.Debug("权限检查通过", "userID", userID)
		c.Next()
	}
}

// PermissionMiddleware 基于菜单权限标识的权限控制中间件
func PermissionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userIDInterface, exists := c.Get("userID")
		if !exists {
			response.Unauthorized(c, code.GetMsg(code.UserAuthFailed))
			c.Abort()
			return
		}

		// 获取需要的权限标识
		requiredPerm := c.GetHeader("X-Required-Permission")
		if requiredPerm == "" {
			// 如果没有设置所需权限，则默认通过
			c.Next()
			return
		}

		// 获取用户ID
		userID := uint(userIDInterface.(float64))

		// 创建菜单服务
		menuService := services.NewMenuService(db)

		// 获取用户所有权限标识
		perms, err := menuService.GetMenuPermsByUserID(userID)
		if err != nil {
			log.Error("获取用户权限标识失败", "error", err)
			response.InternalServerError(c, "权限检查失败")
			c.Abort()
			return
		}

		// 检查是否具有所需权限
		hasPermission := false
		for _, perm := range perms {
			if perm == requiredPerm || perm == "*" {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, code.GetMsg(code.AccessDenied))
			c.Abort()
			return
		}

		c.Next()
	}
}
