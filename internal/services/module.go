package services

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"go.uber.org/fx"
)

// Module 服务模块
var Module = fx.Options(
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
		NewUserService,
		NewAdminUserService,
		NewRoleService,
	),

	// 第三组：菜单服务（依赖角色服务）
	fx.Provide(NewMenuService),

	// 第四组：权限服务（依赖菜单服务）
	fx.Provide(NewPermissionService),
)

// ServiceParams 服务参数
type ServiceParams struct {
	fx.In

	Config *configs.Config
	LC     fx.Lifecycle
	Logger *logger.Logger `optional:"true"`

	// 仓库
	UserRepo         *repository.UserRepo
	AdminUserRepo    *repository.AdminUserRepo
	RoleRepo         *repository.RoleRepo
	MenuRepo         *repository.MenuRepo
	MenuButtonRepo   *repository.MenuButtonRepo
	PermissionRepo   *repository.PermissionRepo
	OperationLogRepo *repository.OperationLogRepo
	FileRepo         *repository.FileRepo
	APIRepo          *repository.APIRepo
	AdminRepo        *repository.AdminRepo
}
