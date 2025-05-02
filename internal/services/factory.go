package services

import (
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
)

// ProvideUserService 提供用户服务
func ProvideUserService(params ServiceParams, authService *AuthService) UserServiceInterface {
	logger.Info("初始化用户服务")
	return NewSeparateUserService(params, authService)
}

// ProvideAdminUserService 提供管理员用户服务
func ProvideAdminUserService(params ServiceParams, casbinService casbin.Service, authService *AuthService) AdminUserServiceInterface {
	logger.Info("初始化管理员用户服务")
	return NewSeparateAdminUserService(params, casbinService, authService)
}

// ProvideRoleService 提供角色服务
func ProvideRoleService(params ServiceParams, casbinService casbin.Service) RoleServiceInterface {
	logger.Info("初始化角色服务")
	return NewRoleService(params, casbinService)
}

// ProvideMenuService 提供菜单服务
func ProvideMenuService(params ServiceParams, casbinService casbin.Service) MenuServiceInterface {
	logger.Info("初始化菜单服务")
	return NewMenuService(params, casbinService)
}

// ProvidePermissionService 提供权限服务
func ProvidePermissionService(params ServiceParams, casbinService casbin.Service, menuService MenuServiceInterface) PermissionServiceInterface {
	logger.Info("初始化权限服务")
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
