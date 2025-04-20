package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
)

func NewPermissionController(permissionService *services.PermissionService) *PermissionController {
	return &PermissionController{
		permissionService: permissionService,
	}
}

type PermissionController struct {
	permissionService *services.PermissionService
}

// 更新权限系统设置
func (pc *PermissionController) UpdatePermissionSettings(c *gin.Context) {
	var req struct {
		Enabled      bool `json:"enabled"`
		DefaultAllow bool `json:"default_allow"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 使用服务层更新权限设置
	err := pc.permissionService.UpdatePermissionSettings(c.Request.Context(), req.Enabled, req.DefaultAllow)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, map[string]any{
		"message":       "权限系统设置已更新",
		"enabled":       req.Enabled,
		"default_allow": req.DefaultAllow,
	})
}

// GetPermissions 获取权限列表
func (pc *PermissionController) GetPermissions(c *gin.Context) {
	permissions, err := pc.permissionService.GetPermissions(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, permissions)
}

// GetPermission 获取权限详情
func (pc *PermissionController) GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	permission, err := pc.permissionService.GetPermission(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, permission)
}

// CreatePermission 创建权限
func (pc *PermissionController) CreatePermission(c *gin.Context) {
	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.Error(c, err)
		return
	}

	if err := pc.permissionService.CreatePermission(c.Request.Context(), &permission); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, permission)
}

// UpdatePermission 更新权限
func (pc *PermissionController) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.Error(c, err)
		return
	}

	if err := pc.permissionService.UpdatePermission(c.Request.Context(), id, &permission); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// DeletePermission 删除权限
func (pc *PermissionController) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 删除权限
	if err := pc.permissionService.DeletePermission(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// GetUserPermissions 获取当前用户的权限列表
func (pc *PermissionController) GetUserPermissions(c *gin.Context) {
	// 从上下文中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	// 获取用户权限
	permissions, err := pc.permissionService.GetPermissionsByUserID(c.Request.Context(), uint(userID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, permissions)
}

// GetUserMenus 获取当前用户的菜单列表（包括按钮权限）
func (pc *PermissionController) GetUserMenus(c *gin.Context) {
	// 从上下文中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	// 获取用户菜单
	menus, err := pc.permissionService.GetUserMenus(c.Request.Context(), strconv.FormatUint(userID, 10))
	if err != nil {
		logger.Error("获取用户菜单失败", "error", err)
		response.Error(c, err)
		return
	}

	response.Success(c, menus)
}

// GetUserRoles 获取当前用户的角色列表
func (pc *PermissionController) GetUserRoles(c *gin.Context) {
	// 从上下文中获取用户ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	// 获取用户角色
	roles, err := pc.permissionService.GetUserRoles(c.Request.Context(), strconv.FormatUint(userID, 10))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, roles)
}

// AssignRolesToUser 为用户分配角色
func (pc *PermissionController) AssignRolesToUser(c *gin.Context) {
	// 获取用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取角色ID列表
	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 分配角色
	if err := pc.permissionService.AssignRolesToUser(c.Request.Context(), strconv.FormatUint(userID, 10), req.RoleIDs); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// AssignPermissionsToRole 为角色分配权限
func (pc *PermissionController) AssignPermissionsToRole(c *gin.Context) {
	// 获取角色ID
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取权限ID列表
	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 分配权限
	if err := pc.permissionService.AssignPermissionToRole(c.Request.Context(), uint(roleID), req.PermissionIDs); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// GetRolePermissions 获取角色的权限列表
func (pc *PermissionController) GetRolePermissions(c *gin.Context) {
	// 获取角色ID
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取角色权限
	permissions, err := pc.permissionService.GetPermissionsByRoleID(c.Request.Context(), uint(roleID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, permissions)
}
