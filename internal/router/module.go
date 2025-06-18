package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/handler"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// RouterParams 路由参数
type RouterParams struct {
	Config *configs.Config

	// 处理器
	UserHandler       *handler.UserHandler
	FileHandler       *handler.FileHandler
	AdminHandler      *handler.AdminHandler
	PermissionHandler *handler.PermissionHandler
}

// NewRouter 创建路由器（不依赖fx）
func NewRouter(config *configs.Config, userHandler *handler.UserHandler, fileHandler *handler.FileHandler, adminHandler *handler.AdminHandler, permissionHandler *handler.PermissionHandler) (*gin.Engine, error) {
	// 设置Gin模式
	gin.SetMode(config.App.Mode)

	// 创建路由器
	r := gin.New()

	// 添加中间件
	r.Use(middleware.RequestLoggerMiddleware())
	r.Use(middleware.Cors())

	// 添加错误处理中间件（替换gin.Recovery()）
	r.Use(middleware.PanicRecovery())
	r.Use(middleware.GlobalErrorHandler())

	// 创建路由注册器
	registrar := &RouteRegistrar{}

	// 注册路由
	registrar.RegisterRoutes(r, config, userHandler, fileHandler, adminHandler, permissionHandler)

	// 打印路由信息
	logger.Info("==================================================")
	logger.Info("Route information:")

	// 获取所有路由
	routes := r.Routes()
	for _, route := range routes {
		logger.Info(fmt.Sprintf("%-7s %s", route.Method, route.Path))
	}

	logger.Info("==================================================")

	return r, nil
}
