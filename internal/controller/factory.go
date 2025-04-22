package controller

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/usermode"
	"github.com/limitcool/starter/internal/services"
)

// ProvideRoleController 根据用户模式提供角色控制器
func ProvideRoleController(
	config *configs.Config,
	roleService services.RoleServiceInterface,
	menuService services.MenuServiceInterface,
	userModeService *usermode.Service,
) RoleControllerInterface {
	// 使用用户模式服务获取用户模式
	logger.InfoContext(context.Background(), "初始化角色控制器", "user_mode", userModeService.GetMode())

	// 创建角色控制器
	roleController := NewRoleController(roleService, menuService)

	// 在简单模式下，我们可以返回一个空实现或者标准实现
	// 这里我们返回标准实现，因为在路由注册时会根据模式决定是否注册相关路由
	return roleController
}

// ProvideMenuController 根据用户模式提供菜单控制器
func ProvideMenuController(
	config *configs.Config,
	menuService services.MenuServiceInterface,
	userModeService *usermode.Service,
) MenuControllerInterface {
	// 使用用户模式服务获取用户模式
	logger.InfoContext(context.Background(), "初始化菜单控制器", "user_mode", userModeService.GetMode())

	// 创建菜单控制器
	menuController := NewMenuController(menuService)

	// 在简单模式下，我们可以返回一个空实现或者标准实现
	// 这里我们返回标准实现，因为在路由注册时会根据模式决定是否注册相关路由
	return menuController
}

// ProvidePermissionController 根据用户模式提供权限控制器
func ProvidePermissionController(
	config *configs.Config,
	permissionService services.PermissionServiceInterface,
	userModeService *usermode.Service,
) PermissionControllerInterface {
	// 使用用户模式服务获取用户模式
	logger.InfoContext(context.Background(), "初始化权限控制器", "user_mode", userModeService.GetMode())

	// 创建权限控制器
	permissionController := NewPermissionController(permissionService)

	// 在简单模式下，我们可以返回一个空实现或者标准实现
	// 这里我们返回标准实现，因为在路由注册时会根据模式决定是否注册相关路由
	return permissionController
}
