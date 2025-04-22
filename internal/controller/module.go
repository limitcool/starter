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

	// 使用工厂函数提供根据用户模式创建的控制器
	fx.Provide(ProvideRoleController),
	fx.Provide(ProvideMenuController),
	fx.Provide(ProvidePermissionController),
	fx.Provide(NewOperationLogController),

	// 提供所有gRPC控制器
	fx.Provide(NewAdminGRPCController),

	// 注册gRPC控制器
	fx.Invoke(RegisterAdminGRPCController),

	// 注册生命周期钩子
	fx.Invoke(RegisterControllerLifecycle),
)

// ControllerParams 控制器参数
type ControllerParams struct {
	fx.In

	Config *configs.Config
	LC     fx.Lifecycle
	Logger *logger.Logger `optional:"true"`

	// 服务接口
	UserService         services.UserServiceInterface
	AdminUserService    services.AdminUserServiceInterface
	RoleService         services.RoleServiceInterface
	MenuService         services.MenuServiceInterface
	MenuAPIService      *services.MenuAPIService
	PermissionService   services.PermissionServiceInterface
	OperationLogService *services.OperationLogService
	FileService         *services.FileService
	APIService          *services.APIService
	AdminService        *services.AdminService

	// 控制器接口
	RoleController       RoleControllerInterface
	MenuController       MenuControllerInterface
	PermissionController PermissionControllerInterface
}
