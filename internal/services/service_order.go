package services

import (
	"go.uber.org/fx"
)

// ServiceOrderGroup 定义服务初始化顺序
var ServiceOrderGroup = fx.Options(
	// 第一组：基础服务
	fx.Provide(
		NewAuthService,
		NewMenuAPIService,
		NewOperationLogService,
		NewFileService,
		NewAPIService,
		NewAdminService,
	),

	// 第二组：用户和角色服务
	fx.Provide(
		ProvideUserService,
		ProvideAdminUserService,
		ProvideRoleService,
	),

	// 第三组：菜单服务（依赖角色服务）
	fx.Provide(ProvideMenuService),

	// 第四组：权限服务（依赖菜单服务）
	fx.Provide(ProvidePermissionService),
)
