package usermode

import (
	"github.com/limitcool/starter/configs"
	"go.uber.org/fx"
)

// Module 用户模式模块
var Module = fx.Options(
	// 提供用户模式服务
	fx.Provide(NewService),
)

// ProvideUserModeService 提供用户模式服务
func ProvideUserModeService(config *configs.Config) *Service {
	return NewService(config)
}
