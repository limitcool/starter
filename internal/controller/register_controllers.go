package controller

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
	"go.uber.org/fx"
)

// RegisterControllersParams 注册控制器参数
type RegisterControllersParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *configs.Config

	// 服务
	RoleService         *services.RoleService
	MenuService         *services.MenuService
	PermissionService   *services.PermissionService
	OperationLogService *services.OperationLogService
}

// RegisterControllers 根据用户模式注册控制器
func RegisterControllers(params RegisterControllersParams) {
	// 获取用户模式
	userMode := enum.GetUserMode(params.Config.Admin.UserMode)

	// 如果是简单模式，不注册角色和菜单相关的控制器
	if userMode == enum.UserModeSimple {
		logger.Info("简单模式：不注册角色和菜单相关的控制器")
		return
	}

	// 分离模式，注册所有控制器
	logger.Info("分离模式：注册所有控制器")

	// 创建角色控制器
	roleController := NewRoleController(params.RoleService, params.MenuService)

	// 创建菜单控制器
	menuController := NewMenuController(params.MenuService)

	// 创建权限控制器
	permissionController := NewPermissionController(params.PermissionService)

	// 创建操作日志控制器
	operationLogController := NewOperationLogController(params.OperationLogService)

	// 注册生命周期钩子
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("角色和菜单相关控制器已注册")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("角色和菜单相关控制器已停止")
			return nil
		},
	})

	// 将控制器注册到全局变量
	Controllers.RoleController = roleController
	Controllers.MenuController = menuController
	Controllers.PermissionController = permissionController
	Controllers.OperationLogController = operationLogController
}
