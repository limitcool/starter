package services

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"go.uber.org/fx"
)

// Module 服务模块
var Module = fx.Options(
	// 提供所有服务
	fx.Provide(NewAuthService),
	fx.Provide(NewUserService),
	fx.Provide(NewSysUserService),
	fx.Provide(NewRoleService),
	fx.Provide(NewMenuService),
	fx.Provide(NewMenuAPIService),
	fx.Provide(NewPermissionService),
	fx.Provide(NewOperationLogService),
	fx.Provide(NewFileService),
	fx.Provide(NewAPIService),
	fx.Provide(NewSystemService),
)

// ServiceParams 服务参数
type ServiceParams struct {
	fx.In

	Config *configs.Config
	LC     fx.Lifecycle
	Logger *logger.Logger `optional:"true"`

	// 仓库
	UserRepo         *repository.UserRepo
	SysUserRepo      *repository.SysUserRepo
	RoleRepo         *repository.RoleRepo
	MenuRepo         *repository.MenuRepo
	MenuButtonRepo   *repository.MenuButtonRepo
	PermissionRepo   *repository.PermissionRepo
	OperationLogRepo *repository.OperationLogRepo
	FileRepo         *repository.FileRepo
	APIRepo          *repository.APIRepo
	SystemRepo       *repository.SystemRepo
}
