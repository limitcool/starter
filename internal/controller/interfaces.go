package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/services"
)

// RoleControllerInterface 角色控制器接口
type RoleControllerInterface interface {
	CreateRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
	GetRole(c *gin.Context)
	GetRoles(c *gin.Context)
	AssignMenuToRole(c *gin.Context)
	SetRolePermission(c *gin.Context)
	DeleteRolePermission(c *gin.Context)
}

// MenuControllerInterface 菜单控制器接口
type MenuControllerInterface interface {
	CreateMenu(c *gin.Context)
	UpdateMenu(c *gin.Context)
	DeleteMenu(c *gin.Context)
	GetMenu(c *gin.Context)
	GetMenuTree(c *gin.Context)
	GetUserMenus(c *gin.Context)
	GetUserMenuPerms(c *gin.Context)
}

// PermissionControllerInterface 权限控制器接口
type PermissionControllerInterface interface {
	UpdatePermissionSettings(c *gin.Context)
	GetPermissions(c *gin.Context)
	GetPermission(c *gin.Context)
	CreatePermission(c *gin.Context)
	UpdatePermission(c *gin.Context)
	DeletePermission(c *gin.Context)
	GetUserPermissions(c *gin.Context)
	GetUserMenus(c *gin.Context)
	GetUserRoles(c *gin.Context)
	AssignRolesToUser(c *gin.Context)
	AssignPermissionsToRole(c *gin.Context)
	GetRolePermissions(c *gin.Context)
	GetPermissionService() services.PermissionServiceInterface
}
