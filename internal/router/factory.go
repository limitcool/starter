package router

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// ProvideRouteRegistrar 提供路由注册器
func ProvideRouteRegistrar(config *configs.Config) RouteRegistrarInterface {
	ctx := context.Background()
	logger.InfoContext(ctx, "初始化路由注册器")
	return &RouteRegistrar{}
}
