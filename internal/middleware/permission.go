package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/jwt"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/permission"
)

// PermissionMiddleware 权限验证中间件
type PermissionMiddleware struct {
	permissionService *permission.Service
}

// NewPermissionMiddleware 创建权限验证中间件
func NewPermissionMiddleware(permissionService *permission.Service) *PermissionMiddleware {
	return &PermissionMiddleware{
		permissionService: permissionService,
	}
}

// RequirePermission 需要特定权限的中间件
func (m *PermissionMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果权限服务未启用，跳过权限检查
		if m.permissionService == nil {
			c.Next()
			return
		}

		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			logger.WarnContext(c.Request.Context(), "权限检查失败：未找到用户ID")
			response.Error(c, errorx.ErrUnauthorized)
			c.Abort()
			return
		}

		// 转换用户ID类型
		var uid int64
		switch v := userID.(type) {
		case int64:
			uid = v
		case string:
			var err error
			uid, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				logger.ErrorContext(c.Request.Context(), "权限检查失败：用户ID类型转换错误", "user_id", userID, "error", err)
				response.Error(c, errorx.ErrUnauthorized)
				c.Abort()
				return
			}
		default:
			logger.ErrorContext(c.Request.Context(), "权限检查失败：用户ID类型不支持", "user_id", userID, "type", v)
			response.Error(c, errorx.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查权限
		hasPermission, err := m.permissionService.CheckPermission(c.Request.Context(), uid, resource, action)
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "权限检查失败", "user_id", uid, "resource", resource, "action", action, "error", err)
			response.Error(c, errorx.ErrInternal)
			c.Abort()
			return
		}

		if !hasPermission {
			logger.WarnContext(c.Request.Context(), "权限不足", "user_id", uid, "resource", resource, "action", action)
			response.Error(c, errorx.ErrForbidden)
			c.Abort()
			return
		}

		logger.DebugContext(c.Request.Context(), "权限检查通过", "user_id", uid, "resource", resource, "action", action)
		c.Next()
	}
}

// RequireAdmin 需要管理员权限的中间件
func (m *PermissionMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果权限服务未启用，跳过权限检查
		if m.permissionService == nil {
			c.Next()
			return
		}

		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			logger.WarnContext(c.Request.Context(), "管理员权限检查失败：未找到用户ID")
			response.Error(c, errorx.ErrUnauthorized)
			c.Abort()
			return
		}

		// 转换用户ID类型
		var uid int64
		switch v := userID.(type) {
		case int64:
			uid = v
		case string:
			var err error
			uid, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				logger.ErrorContext(c.Request.Context(), "管理员权限检查失败：用户ID类型转换错误", "user_id", userID, "error", err)
				response.Error(c, errorx.ErrUnauthorized)
				c.Abort()
				return
			}
		default:
			logger.ErrorContext(c.Request.Context(), "管理员权限检查失败：用户ID类型不支持", "user_id", userID)
			response.Error(c, errorx.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查是否是管理员
		isAdmin, err := m.permissionService.IsAdmin(c.Request.Context(), uid)
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "管理员权限检查失败", "user_id", uid, "error", err)
			response.Error(c, errorx.ErrInternal)
			c.Abort()
			return
		}

		if !isAdmin {
			logger.WarnContext(c.Request.Context(), "非管理员用户访问管理员接口", "user_id", uid)
			response.Error(c, errorx.ErrForbidden)
			c.Abort()
			return
		}

		logger.DebugContext(c.Request.Context(), "管理员权限检查通过", "user_id", uid)
		c.Next()
	}
}

// RequireRole 需要特定角色的中间件
func (m *PermissionMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果权限服务未启用，跳过权限检查
		if m.permissionService == nil {
			c.Next()
			return
		}

		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			logger.WarnContext(c.Request.Context(), "角色权限检查失败：未找到用户ID")
			response.Error(c, errorx.ErrUnauthorized)
			c.Abort()
			return
		}

		// 转换用户ID类型
		var uid int64
		switch v := userID.(type) {
		case int64:
			uid = v
		case string:
			var err error
			uid, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				logger.ErrorContext(c.Request.Context(), "角色权限检查失败：用户ID类型转换错误", "user_id", userID, "error", err)
				response.Error(c, errorx.ErrUnauthorized)
				c.Abort()
				return
			}
		default:
			logger.ErrorContext(c.Request.Context(), "角色权限检查失败：用户ID类型不支持", "user_id", userID)
			response.Error(c, errorx.ErrUnauthorized)
			c.Abort()
			return
		}

		// 获取用户角色
		roles, err := m.permissionService.GetUserRoles(c.Request.Context(), uid)
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "角色权限检查失败", "user_id", uid, "error", err)
			response.Error(c, errorx.ErrInternal)
			c.Abort()
			return
		}

		// 检查是否有指定角色
		hasRole := false
		for _, role := range roles {
			if role.Name == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			logger.WarnContext(c.Request.Context(), "角色权限不足", "user_id", uid, "required_role", roleName)
			response.Error(c, errorx.ErrForbidden)
			c.Abort()
			return
		}

		logger.DebugContext(c.Request.Context(), "角色权限检查通过", "user_id", uid, "role", roleName)
		c.Next()
	}
}

// RequirePermissionKey 需要特定权限Key的中间件
// permissionKey: 权限标识，如 "user:create", "course:list"
func (m *PermissionMiddleware) RequirePermissionKey(permissionKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果权限服务未启用，跳过权限检查
		if m.permissionService == nil {
			logger.WarnContext(c.Request.Context(), "权限服务未启用，跳过权限检查", "permission_key", permissionKey)
			c.Next()
			return
		}

		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			logger.WarnContext(c.Request.Context(), "权限检查失败：未找到用户ID")
			response.Error(c, errorx.ErrUnauthorized)
			c.Abort()
			return
		}

		// 转换用户ID类型
		var uid int64
		switch v := userID.(type) {
		case int64:
			uid = v
		case string:
			var err error
			uid, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				logger.ErrorContext(c.Request.Context(), "权限检查失败：用户ID类型转换错误", "user_id", userID, "error", err)
				response.Error(c, errorx.ErrUnauthorized)
				c.Abort()
				return
			}
		default:
			logger.ErrorContext(c.Request.Context(), "权限检查失败：用户ID类型不支持", "user_id", userID, "type", fmt.Sprintf("%T", v))
			response.Error(c, errorx.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查用户是否是管理员（管理员拥有所有权限）
		isAdminInterface, exists := c.Get("is_admin")
		if exists {
			if isAdmin, ok := isAdminInterface.(bool); ok && isAdmin {
				logger.InfoContext(c.Request.Context(), "管理员用户，跳过权限检查", "user_id", uid, "permission_key", permissionKey)
				c.Next()
				return
			}
		}

		// 解析权限Key，格式如 "user:create"
		resource, action := parsePermissionKey(permissionKey)
		if resource == "" || action == "" {
			logger.ErrorContext(c.Request.Context(), "权限Key格式错误", "permission_key", permissionKey)
			response.Error(c, errorx.ErrInternal.WithMsg("权限配置错误"))
			c.Abort()
			return
		}

		// 检查权限
		hasPermission, err := m.permissionService.CheckPermission(c.Request.Context(), uid, resource, action)
		if err != nil {
			logger.ErrorContext(c.Request.Context(), "权限检查失败",
				"user_id", uid,
				"permission_key", permissionKey,
				"resource", resource,
				"action", action,
				"error", err)
			response.Error(c, errorx.ErrInternal)
			c.Abort()
			return
		}

		if !hasPermission {
			logger.WarnContext(c.Request.Context(), "权限不足",
				"user_id", uid,
				"permission_key", permissionKey,
				"resource", resource,
				"action", action)
			response.Error(c, errorx.ErrForbidden.WithMsg(fmt.Sprintf("缺少权限: %s", permissionKey)))
			c.Abort()
			return
		}

		logger.InfoContext(c.Request.Context(), "权限检查通过",
			"user_id", uid,
			"permission_key", permissionKey,
			"resource", resource,
			"action", action)

		c.Next()
	}
}

// parsePermissionKey 解析权限Key，格式如 "user:create" -> ("user", "create")
func parsePermissionKey(permissionKey string) (resource, action string) {
	parts := strings.Split(permissionKey, ":")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// GetUserIDFromToken 从token中获取用户ID的辅助函数
func GetUserIDFromToken(c *gin.Context, secret string) (int64, error) {
	// 从Authorization header获取token
	token := c.GetHeader("Authorization")
	if token == "" {
		return 0, errorx.ErrUnauthorized.WithMsg("缺少Authorization header")
	}

	// 移除Bearer前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 解析token
	claims, err := jwt.ParseTokenWithCustomClaims(token, secret)
	if err != nil {
		return 0, errorx.ErrUnauthorized.WithMsg("无效的token")
	}

	return claims.UserID, nil
}
