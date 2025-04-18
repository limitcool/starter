package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
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
	err := pc.permissionService.UpdatePermissionSettings(req.Enabled, req.DefaultAllow)
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
	permissions, err := pc.permissionService.GetPermissions()
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

	permission, err := pc.permissionService.GetPermission(id)
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

	if err := pc.permissionService.CreatePermission(&permission); err != nil {
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

	if err := pc.permissionService.UpdatePermission(id, &permission); err != nil {
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
	if err := pc.permissionService.DeletePermission(id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}
