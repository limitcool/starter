package router

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/handler"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
)

// Module 路由模块
var Module = fx.Options(
	// 提供路由注册器
	fx.Provide(ProvideRouteRegistrar),

	// 提供路由
	fx.Provide(NewRouter),
)

// RouterParams 路由参数
type RouterParams struct {
	fx.In

	Config *configs.Config
	LC     fx.Lifecycle
	Logger *logger.Logger `optional:"true"`

	// 路由注册器
	RouteRegistrar RouteRegistrarInterface

	// 处理器
	UserHandler  *handler.UserHandler
	FileHandler  *handler.FileHandler
	AdminHandler *handler.AdminHandler
}

// RouterResult 路由结果
type RouterResult struct {
	fx.Out

	Router *gin.Engine
}

// NewRouter 创建路由
func NewRouter(params RouterParams) RouterResult {
	// 设置Gin模式
	if params.Config.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if params.Config.App.Mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建路由
	r := gin.New()

	// 添加中间件
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLoggerMiddleware())
	r.Use(middleware.RequestContext())
	r.Use(middleware.Cors())

	// 添加错误处理中间件
	// PanicRecovery: 用于捕获 panic 并返回友好的错误响应
	r.Use(middleware.PanicRecovery())

	// GlobalErrorHandler: 用于处理控制器方法通过 c.Error() 返回的错误
	// 配合 ErrorHandlerFunc 使用，可以简化控制器代码
	r.Use(middleware.GlobalErrorHandler())

	// 使用用户模式服务
	ctx := context.Background()

	// 在lite版本中，不使用Casbin中间件
	logger.InfoContext(ctx, "Lite版本不使用Casbin中间件")

	// 注册路由
	registerRoutes(r, params)

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "Router initialized successfully")

			// 打印路由信息
			logger.InfoContext(ctx, "==================================================")
			logger.InfoContext(ctx, "路由信息:")

			// 获取所有路由
			routes := r.Routes()
			for _, route := range routes {
				logger.InfoContext(ctx, fmt.Sprintf("%-7s %s", route.Method, route.Path))
			}

			logger.InfoContext(ctx, "==================================================")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "Router stopped")
			return nil
		},
	})

	return RouterResult{Router: r}
}

// registerRoutes 注册路由
func registerRoutes(r *gin.Engine, params RouterParams) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API版本
	v1 := r.Group("/api/v1")
	{
		// 使用路由注册器注册路由
		params.RouteRegistrar.RegisterRoutes(v1, params)

		// 其他路由可以在这里添加
	}
}
