package controller

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
)

// LifecycleParams 生命周期参数
type LifecycleParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *configs.Config

	// 控制器接口
	RoleController       RoleControllerInterface       `optional:"true"`
	MenuController       MenuControllerInterface       `optional:"true"`
	PermissionController PermissionControllerInterface `optional:"true"`
}

// RegisterControllerLifecycle 注册控制器生命周期钩子
func RegisterControllerLifecycle(params LifecycleParams) {
	// 注册生命周期钩子
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "控制器已注册")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "控制器已停止")
			return nil
		},
	})

	// 不再使用全局变量注册控制器
	// 所有控制器都通过依赖注入获取
}
