package services

import (
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
)

// ProvideUserService 根据用户模式提供用户服务
func ProvideUserService(params ServiceParams, authService *AuthService) UserServiceInterface {
	// 获取用户模式
	userMode := enum.GetUserMode(params.Config.Admin.UserMode)
	logger.Info("初始化用户服务", "user_mode", userMode)

	// 根据用户模式创建对应的服务实现
	if userMode == enum.UserModeSimple {
		return NewSimpleUserService(params, authService)
	} else {
		return NewSeparateUserService(params, authService)
	}
}

// ProvideAdminUserService 根据用户模式提供管理员用户服务
func ProvideAdminUserService(params ServiceParams, casbinService casbin.Service, authService *AuthService) AdminUserServiceInterface {
	// 获取用户模式
	userMode := enum.GetUserMode(params.Config.Admin.UserMode)
	logger.Info("初始化管理员用户服务", "user_mode", userMode)

	// 根据用户模式创建对应的服务实现
	if userMode == enum.UserModeSimple {
		return NewSimpleAdminUserService(params, authService)
	} else {
		return NewSeparateAdminUserService(params, casbinService, authService)
	}
}

// ProvideRoleService 根据用户模式提供角色服务
func ProvideRoleService(params ServiceParams, casbinService casbin.Service) RoleServiceInterface {
	// 获取用户模式
	userMode := enum.GetUserMode(params.Config.Admin.UserMode)
	logger.Info("初始化角色服务", "user_mode", userMode)

	// 使用现有的RoleService，它已经根据用户模式提供了不同的实现
	return NewRoleService(params, casbinService)
}

// ProvideMenuService 根据用户模式提供菜单服务
func ProvideMenuService(params ServiceParams, casbinService casbin.Service) MenuServiceInterface {
	// 获取用户模式
	userMode := enum.GetUserMode(params.Config.Admin.UserMode)
	logger.Info("初始化菜单服务", "user_mode", userMode)

	// 使用现有的MenuService，它已经根据用户模式提供了不同的实现
	return NewMenuService(params, casbinService)
}

// ProvidePermissionService 根据用户模式提供权限服务
func ProvidePermissionService(params ServiceParams, casbinService casbin.Service, menuService MenuServiceInterface) PermissionServiceInterface {
	// 获取用户模式
	userMode := enum.GetUserMode(params.Config.Admin.UserMode)
	logger.Info("初始化权限服务", "user_mode", userMode)

	// 使用现有的PermissionService，它已经根据用户模式提供了不同的实现
	return NewPermissionService(params, casbinService, menuService)
}

// RegisterServices 注册所有服务
func RegisterServices() fx.Option {
	return fx.Options(
		// 提供通用服务
		fx.Provide(NewAuthService),

		// 提供根据用户模式创建的服务
		fx.Provide(ProvideUserService),
		fx.Provide(ProvideAdminUserService),
		fx.Provide(ProvideRoleService),
		fx.Provide(ProvideMenuService),
		fx.Provide(ProvidePermissionService),
	)
}
