package router

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
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
	SysUserController      *controller.SysUserController
	RoleController         *controller.RoleController
	MenuController         *controller.MenuController
	PermissionController   *controller.PermissionController
	OperationLogController *controller.OperationLogController
	FileController         *controller.FileController
	APIController          *controller.APIController
	SystemController       *controller.SystemController
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

	// 添加Casbin中间件（如果启用）
	if params.Config.Casbin.Enabled && params.Enforcer != nil {
		r.Use(middleware.CasbinMiddleware(params.PermissionController.GetPermissionService(), params.Config))
	}

	// 注册路由
	registerRoutes(r, params)

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Router initialized successfully")
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

	// API版本
	v1 := r.Group("/api/v1")
	{
		// 用户路由
		userRoutes(v1, params)

		// 其他路由可以在这里添加
	}
}

// userRoutes 用户路由
func userRoutes(r *gin.RouterGroup, params RouterParams) {
	// 用户路由
	users := r.Group("/users")
	{
		// 用户登录
		users.POST("/login", params.UserController.UserLogin)

		// 用户注册
		users.POST("/register", params.UserController.UserRegister)

		// 刷新令牌
		users.POST("/refresh-token", params.UserController.RefreshToken)

		// 获取用户信息
		users.GET("/info", params.UserController.UserInfo)

		// 修改密码
		users.POST("/change-password", params.UserController.UserChangePassword)
	}
}
