package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/permission"
)

// UserMenuHandler 用户菜单处理器
type UserMenuHandler struct {
	permissionService *permission.Service
}

// NewUserMenuHandler 创建用户菜单处理器
func NewUserMenuHandler(permissionService *permission.Service) *UserMenuHandler {
	return &UserMenuHandler{
		permissionService: permissionService,
	}
}

// GetMyMenus 获取当前用户的菜单
func (h *UserMenuHandler) GetMyMenus(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.WarnContext(c.Request.Context(), "获取用户菜单失败：未找到用户ID")
		response.Error(c, errorx.ErrUnauthorized.WithMsg("用户未登录"))
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
			logger.ErrorContext(c.Request.Context(), "获取用户菜单失败：用户ID类型转换错误", "user_id", userID, "error", err)
			response.Error(c, errorx.ErrUnauthorized.WithMsg("用户ID无效"))
			return
		}
	default:
		logger.ErrorContext(c.Request.Context(), "获取用户菜单失败：用户ID类型不支持", "user_id", userID)
		response.Error(c, errorx.ErrUnauthorized.WithMsg("用户ID类型错误"))
		return
	}

	// 获取平台参数，默认为admin
	platform := c.DefaultQuery("platform", "admin")

	// 验证平台参数
	if platform != "admin" && platform != "coach_mp" {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("不支持的平台类型"))
		return
	}

	// 获取用户菜单
	menus, err := h.permissionService.GetUserMenusByPlatform(c.Request.Context(), uid, platform)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取用户菜单失败", "user_id", uid, "platform", platform, "error", err)
		response.Error(c, err)
		return
	}

	// 转换为响应格式
	menuResponses := h.convertToMenuResponses(menus)

	logger.InfoContext(c.Request.Context(), "获取用户菜单成功", "user_id", uid, "platform", platform, "menu_count", len(menuResponses))
	response.Success(c, menuResponses)
}

// GetMyPermissions 获取当前用户的权限列表
func (h *UserMenuHandler) GetMyPermissions(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.WarnContext(c.Request.Context(), "获取用户权限失败：未找到用户ID")
		response.Error(c, errorx.ErrUnauthorized.WithMsg("用户未登录"))
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
			logger.ErrorContext(c.Request.Context(), "获取用户权限失败：用户ID类型转换错误", "user_id", userID, "error", err)
			response.Error(c, errorx.ErrUnauthorized.WithMsg("用户ID无效"))
			return
		}
	default:
		logger.ErrorContext(c.Request.Context(), "获取用户权限失败：用户ID类型不支持", "user_id", userID)
		response.Error(c, errorx.ErrUnauthorized.WithMsg("用户ID类型错误"))
		return
	}

	// 获取用户权限
	permissions, err := h.permissionService.GetUserPermissions(c.Request.Context(), uid)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取用户权限失败", "user_id", uid, "error", err)
		response.Error(c, err)
		return
	}

	// 转换为响应格式
	permissionResponses := h.convertToPermissionResponses(permissions)

	logger.InfoContext(c.Request.Context(), "获取用户权限成功", "user_id", uid, "permission_count", len(permissionResponses))
	response.Success(c, permissionResponses)
}

// GetMyRoles 获取当前用户的角色列表
func (h *UserMenuHandler) GetMyRoles(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.WarnContext(c.Request.Context(), "获取用户角色失败：未找到用户ID")
		response.Error(c, errorx.ErrUnauthorized.WithMsg("用户未登录"))
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
			logger.ErrorContext(c.Request.Context(), "获取用户角色失败：用户ID类型转换错误", "user_id", userID, "error", err)
			response.Error(c, errorx.ErrUnauthorized.WithMsg("用户ID无效"))
			return
		}
	default:
		logger.ErrorContext(c.Request.Context(), "获取用户角色失败：用户ID类型不支持", "user_id", userID)
		response.Error(c, errorx.ErrUnauthorized.WithMsg("用户ID类型错误"))
		return
	}

	// 获取用户角色
	roles, err := h.permissionService.GetUserRoles(c.Request.Context(), uid)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取用户角色失败", "user_id", uid, "error", err)
		response.Error(c, err)
		return
	}

	// 转换为响应格式
	roleResponses := h.convertToRoleResponses(roles)

	logger.InfoContext(c.Request.Context(), "获取用户角色成功", "user_id", uid, "role_count", len(roleResponses))
	response.Success(c, roleResponses)
}

// CheckMyPermission 检查当前用户是否有指定权限
func (h *UserMenuHandler) CheckMyPermission(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.WarnContext(c.Request.Context(), "检查用户权限失败：未找到用户ID")
		response.Error(c, errorx.ErrUnauthorized.WithMsg("用户未登录"))
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
			logger.ErrorContext(c.Request.Context(), "检查用户权限失败：用户ID类型转换错误", "user_id", userID, "error", err)
			response.Error(c, errorx.ErrUnauthorized.WithMsg("用户ID无效"))
			return
		}
	default:
		logger.ErrorContext(c.Request.Context(), "检查用户权限失败：用户ID类型不支持", "user_id", userID)
		response.Error(c, errorx.ErrUnauthorized.WithMsg("用户ID类型错误"))
		return
	}

	var req dto.CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	// 检查权限
	hasPermission, err := h.permissionService.CheckPermission(c.Request.Context(), uid, req.Resource, req.Action)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "检查用户权限失败", "user_id", uid, "resource", req.Resource, "action", req.Action, "error", err)
		response.Error(c, err)
		return
	}

	result := dto.CheckPermissionResponse{
		HasPermission: hasPermission,
	}

	logger.InfoContext(c.Request.Context(), "检查用户权限完成", "user_id", uid, "resource", req.Resource, "action", req.Action, "has_permission", hasPermission)
	response.Success(c, result)
}

// convertToMenuResponses 转换菜单为响应格式
func (h *UserMenuHandler) convertToMenuResponses(menus []model.Menu) []dto.MenuResponse {
	var responses []dto.MenuResponse
	for _, menu := range menus {
		response := dto.MenuResponse{
			ID:            menu.ID,
			ParentID:      menu.ParentID,
			Name:          menu.Name,
			Path:          menu.Path,
			Component:     menu.Component,
			Icon:          menu.Icon,
			SortOrder:     menu.SortOrder,
			IsVisible:     menu.IsVisible,
			PermissionKey: menu.PermissionKey,
			Platform:      menu.Platform,
			CreatedAt:     menu.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     menu.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// 递归处理子菜单
		if len(menu.Children) > 0 {
			response.Children = h.convertToMenuResponses(menu.Children)
		}

		responses = append(responses, response)
	}
	return responses
}

// convertToPermissionResponses 转换权限为响应格式
func (h *UserMenuHandler) convertToPermissionResponses(permissions []model.Permission) []dto.PermissionResponse {
	var responses []dto.PermissionResponse
	for _, permission := range permissions {
		response := dto.PermissionResponse{
			ID:        permission.ID,
			ParentID:  permission.ParentID,
			Name:      permission.Name,
			Key:       permission.Key,
			Type:      permission.Type,
			CreatedAt: permission.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: permission.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		responses = append(responses, response)
	}
	return responses
}

// convertToRoleResponses 转换角色为响应格式
func (h *UserMenuHandler) convertToRoleResponses(roles []model.Role) []dto.RoleResponse {
	var responses []dto.RoleResponse
	for _, role := range roles {
		response := dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Key:         role.Key,
			Description: role.Description,
			Status:      int(role.Status),
			CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		responses = append(responses, response)
	}
	return responses
}
