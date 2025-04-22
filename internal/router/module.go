package router

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"go.uber.org/fx"
)

// Module 路由模块
var Module = fx.Options(
	// 提供路由
	fx.Provide(NewRouter),
)

// RouterParams 路由参数
type RouterParams struct {
	fx.In

	Config   *configs.Config
	LC       fx.Lifecycle
	Enforcer *casbin.Enforcer `optional:"true"`
	Logger   *logger.Logger   `optional:"true"`

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

	// 获取用户模式
	userMode := enum.GetUserMode(params.Config.Admin.UserMode)

	// 在分离模式下添加Casbin中间件（如果启用）
	if userMode == enum.UserModeSeparate && params.Config.Casbin.Enabled && params.Enforcer != nil {
		r.Use(middleware.CasbinMiddleware(params.PermissionController.GetPermissionService(), params.Config))
	}

	// 注册路由
	registerRoutes(r, params)

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Router initialized successfully")

			// 打印路由信息
			logger.Info("==================================================")
			logger.Info("路由信息:")

			// 获取所有路由
			routes := r.Routes()
			for _, route := range routes {
				logger.Info(fmt.Sprintf("%-7s %s", route.Method, route.Path))
			}

			logger.Info("==================================================")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Router stopped")
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

	// 获取用户模式
	userMode := enum.GetUserMode(params.Config.Admin.UserMode)

	// API版本
	v1 := r.Group("/api/v1")
	{
		// 根据用户模式注册不同的路由
		if userMode == enum.UserModeSimple {
			// 简单模式 - 使用简化的路由
			registerSimpleRoutes(v1, params)
		} else {
			// 分离模式 - 使用完整的路由
			registerSeparateRoutes(v1, params)
		}

		// 其他路由可以在这里添加
	}
}
