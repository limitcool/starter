package repository

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module 仓库模块
var Module = fx.Options(
	// 提供所有仓库
	fx.Provide(NewUserRepo),
	fx.Provide(NewAdminUserRepo),
	fx.Provide(NewRoleRepo),
	fx.Provide(NewMenuRepo),
	fx.Provide(NewMenuButtonRepo),
	fx.Provide(NewPermissionRepo),
	fx.Provide(NewOperationLogRepo),
	fx.Provide(NewFileRepo),
	fx.Provide(NewAPIRepo),
	fx.Provide(NewAdminSystemRepo),
	// 不再提供旧的类型
)

// RepoParams 仓库参数
type RepoParams struct {
	fx.In

	DB     *gorm.DB
	Config *configs.Config
	LC     fx.Lifecycle
	Logger *logger.Logger `optional:"true"`
}
