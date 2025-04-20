package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
)

// CasbinMiddleware 基于路径和方法的权限控制中间件
// 用于检查用户是否有权限访问特定的API路径和方法
func CasbinMiddleware(permissionService *services.PermissionService, config *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查权限系统是否启用
		if !config.Casbin.Enabled {
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

		ctx := c.Request.Context()
		logger.DebugContext(ctx, "检查权限", "userID", userID, "object", obj, "action", act)

		// 检查权限
		pass, err := permissionService.CheckPermission(ctx, userID, obj, act)
		if err != nil {
			logger.ErrorContext(ctx, "权限检查错误", "error", err)
			response.Error(c, errorx.ErrCasbinService)
			c.Abort()
			return
		}

		if !pass {
			// 尝试获取用户角色
			roles, err := permissionService.GetUserRoles(c.Request.Context(), userID)
			if err == nil {
				var roleNames []string
				for _, role := range roles {
					roleNames = append(roleNames, role.Name)
				}
				ctx := c.Request.Context()
				logger.DebugContext(ctx, "权限检查失败", "userID", userID, "roles", strings.Join(roleNames, ","))
			}

			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		logger.DebugContext(ctx, "权限检查通过", "userID", userID)
		c.Next()
	}
}

// PermissionCodeMiddleware 基于权限编码的权限控制中间件
// 用于检查用户是否有权限访问特定的权限编码
// 权限编码通过请求头 X-Required-Permission 指定
func PermissionCodeMiddleware(permissionService *services.PermissionService, config *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查权限系统是否启用
		if !config.Casbin.Enabled {
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

		// 获取用户角色
		roles, err := permissionService.GetUserRoles(c.Request.Context(), userID)
		if err != nil {
			ctx := c.Request.Context()
			logger.ErrorContext(ctx, "获取用户角色失败", "error", err)
			response.Error(c, errorx.ErrCasbinService)
			c.Abort()
			return
		}

		// 检查角色是否有所需权限
		hasPermission := false
		for _, role := range roles {
			// 检查是否为管理员
			if role.Code == "admin" {
				hasPermission = true
				break
			}

			// 检查角色是否有权限
			pass, err := permissionService.CheckPermission(c.Request.Context(), role.Code, requiredPerm, "*")
			if err != nil {
				ctx := c.Request.Context()
				logger.ErrorContext(ctx, "权限检查错误", "error", err)
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

// RequirePermission 创建一个需要特定权限的中间件
// 用于检查用户是否有权限访问特定的权限编码
// 与 PermissionCodeMiddleware 不同的是，权限编码在创建中间件时指定
func RequirePermission(permCode string, permissionService *services.PermissionService, config *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查权限系统是否启用
		if !config.Casbin.Enabled {
			// 权限系统未启用，直接放行
			c.Next()
			return
		}

		// 从上下文中获取用户ID
		userID := GetUserID(c)
		if userID == 0 {
			response.Error(c, errorx.ErrUserNoLogin)
			c.Abort()
			return
		}

		// 获取用户角色
		roles, err := permissionService.GetUserRoles(c.Request.Context(), strconv.FormatUint(userID, 10))
		if err != nil {
			ctx := c.Request.Context()
			logger.ErrorContext(ctx, "获取用户角色失败", "error", err)
			response.Error(c, errorx.ErrCasbinService)
			c.Abort()
			return
		}

		// 检查角色是否有所需权限
		hasPermission := false
		for _, role := range roles {
			// 检查是否为管理员
			if role.Code == "admin" {
				hasPermission = true
				break
			}

			// 检查角色是否有权限
			pass, err := permissionService.CheckPermission(c.Request.Context(), role.Code, permCode, "*")
			if err != nil {
				ctx := c.Request.Context()
				logger.ErrorContext(ctx, "权限检查错误", "error", err)
				continue
			}

			if pass {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			ctx := c.Request.Context()
			logger.WarnContext(ctx, "权限检查失败", "userID", userID, "permCode", permCode)
			response.Error(c, errorx.ErrAccessDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}
