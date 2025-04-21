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
	fx.Provide(NewAdminUserController),
	fx.Provide(NewFileController),
	fx.Provide(NewAPIController),
	fx.Provide(NewAdminController),

	// 根据用户模式决定是否提供角色和菜单相关的控制器
	fx.Invoke(RegisterControllers),

	// 提供所有gRPC控制器
	fx.Provide(NewAdminGRPCController),

	// 注册gRPC控制器
	fx.Invoke(RegisterAdminGRPCController),
)

// ControllerParams 控制器参数
type ControllerParams struct {
	fx.In

	Config *configs.Config
	LC     fx.Lifecycle
	Logger *logger.Logger `optional:"true"`

	// 服务
	UserService         *services.UserService
	AdminUserService    *services.AdminUserService
	RoleService         *services.RoleService
	MenuService         *services.MenuService
	MenuAPIService      *services.MenuAPIService
	PermissionService   *services.PermissionService
	OperationLogService *services.OperationLogService
	FileService         *services.FileService
	APIService          *services.APIService
	AdminService        *services.AdminService
}
