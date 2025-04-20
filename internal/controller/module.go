package controller

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
	"go.uber.org/fx"
)

// Module 控制器模块
var Module = fx.Options(
	// 提供所有HTTP控制器
	fx.Provide(NewUserController),
	fx.Provide(NewSysUserController),
	fx.Provide(NewRoleController),
	fx.Provide(NewMenuController),
	fx.Provide(NewPermissionController),
	fx.Provide(NewOperationLogController),
	fx.Provide(NewFileController),
	fx.Provide(NewAPIController),
	fx.Provide(NewSystemController),

	// 提供所有gRPC控制器
	fx.Provide(NewSystemGRPCController),

	// 注册gRPC控制器
	fx.Invoke(RegisterSystemGRPCController),
)

// ControllerParams 控制器参数
type ControllerParams struct {
	fx.In

	Config *configs.Config
	LC     fx.Lifecycle
	Logger *logger.Logger `optional:"true"`

	// 服务
	UserService         *services.UserService
	SysUserService      *services.SysUserService
	RoleService         *services.RoleService
	MenuService         *services.MenuService
	MenuAPIService      *services.MenuAPIService
	PermissionService   *services.PermissionService
	OperationLogService *services.OperationLogService
	FileService         *services.FileService
	APIService          *services.APIService
	SystemService       *services.SystemService
}
