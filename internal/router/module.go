package router

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/usermode"
	"github.com/limitcool/starter/internal/repository"
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

	Config         *configs.Config
	LC             fx.Lifecycle
	Enforcer       *casbin.Enforcer `optional:"true"`
	Logger         *logger.Logger   `optional:"true"`
	UserModeService *usermode.Service

	// 路由注册器
	RouteRegistrar RouteRegistrarInterface

	// 控制器
	UserController         *controller.UserController
	AdminUserController    *controller.AdminUserController
	RoleController         controller.RoleControllerInterface       `optional:"true"`
	MenuController         controller.MenuControllerInterface       `optional:"true"`
	PermissionController   controller.PermissionControllerInterface `optional:"true"`
	OperationLogController *controller.OperationLogController       `optional:"true"`
	FileController         *controller.FileController
	APIController          *controller.APIController
	AdminController        *controller.AdminController

	// 仓库
	UserRepo *repository.UserRepo `optional:"true"`
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
	r.Use(middleware.RequestLogger())
	r.Use(middleware.RequestContext())
	r.Use(middleware.Cors())
	r.Use(middleware.ErrorHandler())

	// 使用用户模式服务
	ctx := context.Background()

	// 在分离模式下添加Casbin中间件（如果启用）
	if params.UserModeService.IsSeparateMode() && params.Config.Casbin.Enabled && params.Enforcer != nil {
		// 检查PermissionController是否为空
		if params.PermissionController != nil {
			r.Use(middleware.CasbinMiddleware(params.PermissionController.GetPermissionService(), params.Config))
			logger.InfoContext(ctx, "添加Casbin中间件")
		} else {
			logger.WarnContext(ctx, "权限控制器为空，无法添加Casbin中间件")
		}
	}

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
