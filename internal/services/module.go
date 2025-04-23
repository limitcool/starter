package services

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"go.uber.org/fx"
)

// 注意: Module 变量已经移动到 service_order.go 文件中
// 请使用 ServiceOrderGroup 替代 Module

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
