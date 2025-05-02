package controller

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
)

// ProvideRoleController 提供角色控制器
func ProvideRoleController(
	config *configs.Config,
	roleService services.RoleServiceInterface,
	menuService services.MenuServiceInterface,
) RoleControllerInterface {
	// 记录初始化日志
	logger.InfoContext(context.Background(), "初始化角色控制器")

	// 创建角色控制器
	roleController := NewRoleController(roleService, menuService)
	return roleController
}

// ProvideMenuController 提供菜单控制器
func ProvideMenuController(
	config *configs.Config,
	menuService services.MenuServiceInterface,
) MenuControllerInterface {
	// 记录初始化日志
	logger.InfoContext(context.Background(), "初始化菜单控制器")

	// 创建菜单控制器
	menuController := NewMenuController(menuService)
	return menuController
}

// ProvidePermissionController 提供权限控制器
func ProvidePermissionController(
	config *configs.Config,
	permissionService services.PermissionServiceInterface,
) PermissionControllerInterface {
	// 记录初始化日志
	logger.InfoContext(context.Background(), "初始化权限控制器")

	// 创建权限控制器
	permissionController := NewPermissionController(permissionService)
	return permissionController
}
