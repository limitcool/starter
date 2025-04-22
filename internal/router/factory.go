package router

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/usermode"
)

// ProvideRouteRegistrar 根据用户模式提供路由注册器
func ProvideRouteRegistrar(config *configs.Config, userModeService *usermode.Service) RouteRegistrarInterface {
	// 使用用户模式服务获取用户模式
	ctx := context.Background()
	logger.InfoContext(ctx, "初始化路由注册器", "user_mode", userModeService.GetMode())

	// 根据用户模式返回不同的路由注册器
	if userModeService.IsSimpleMode() {
		return &SimpleRouteRegistrar{}
	} else {
		return &SeparateRouteRegistrar{}
	}
}
