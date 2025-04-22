package controller

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/usermode"
	"go.uber.org/fx"
)

// LifecycleParams 生命周期参数
type LifecycleParams struct {
	fx.In

	Lifecycle      fx.Lifecycle
	Config         *configs.Config
	UserModeService *usermode.Service

	// 控制器接口
	RoleController       RoleControllerInterface       `optional:"true"`
	MenuController       MenuControllerInterface       `optional:"true"`
	PermissionController PermissionControllerInterface `optional:"true"`
}

// RegisterControllerLifecycle 注册控制器生命周期钩子
func RegisterControllerLifecycle(params LifecycleParams) {
	// 使用用户模式服务获取用户模式
	userMode := params.UserModeService.GetMode()

	// 注册生命周期钩子
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.UserModeService.IsSimpleMode() {
				logger.InfoContext(ctx, "简单模式控制器已注册", "user_mode", userMode)
			} else {
				logger.InfoContext(ctx, "分离模式控制器已注册", "user_mode", userMode)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if params.UserModeService.IsSimpleMode() {
				logger.InfoContext(ctx, "简单模式控制器已停止", "user_mode", userMode)
			} else {
				logger.InfoContext(ctx, "分离模式控制器已停止", "user_mode", userMode)
			}
			return nil
		},
	})

	// 不再使用全局变量注册控制器
	// 所有控制器都通过依赖注入获取
}
