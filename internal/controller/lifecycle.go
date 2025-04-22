package controller

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/enum"
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
	// 获取用户模式
	userMode := enum.GetUserMode(params.Config.Admin.UserMode)

	// 注册生命周期钩子
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if userMode == enum.UserModeSimple {
				logger.Info("简单模式控制器已注册")
			} else {
				logger.Info("分离模式控制器已注册")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if userMode == enum.UserModeSimple {
				logger.Info("简单模式控制器已停止")
			} else {
				logger.Info("分离模式控制器已停止")
			}
			return nil
		},
	})

	// 将控制器注册到全局变量（如果需要）
	if userMode == enum.UserModeSeparate {
		if params.RoleController != nil {
			Controllers.RoleController = params.RoleController.(*RoleController)
		}
		if params.MenuController != nil {
			Controllers.MenuController = params.MenuController.(*MenuController)
		}
		if params.PermissionController != nil {
			Controllers.PermissionController = params.PermissionController.(*PermissionController)
		}
	}
}
