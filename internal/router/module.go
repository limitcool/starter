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
	RoleController         *controller.RoleController
	MenuController         *controller.MenuController
	PermissionController   *controller.PermissionController
	OperationLogController *controller.OperationLogController
	FileController         *controller.FileController
	APIController          *controller.APIController
	AdminController        *controller.AdminController
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

	// API版本
	v1 := r.Group("/api/v1")
	{
		// 用户路由
		userRoutes(v1, params)

		// 管理员路由
		adminRoutes(v1, params)

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

// adminRoutes 管理员路由
func adminRoutes(r *gin.RouterGroup, params RouterParams) {
	// 管理员路由
	admin := r.Group("/admin-api")
	{
		// 管理员登录
		admin.POST("/login", params.AdminUserController.AdminUserLogin)

		// 角色管理
		roles := admin.Group("/roles")
		{
			// 获取角色列表
			roles.GET("", params.RoleController.GetRoles)

			// 获取角色详情
			roles.GET("/:id", params.RoleController.GetRole)

			// 创建角色
			roles.POST("", params.RoleController.CreateRole)

			// 更新角色
			roles.PUT("/:id", params.RoleController.UpdateRole)

			// 删除角色
			roles.DELETE("/:id", params.RoleController.DeleteRole)

			// 为角色分配菜单
			roles.POST("/menu", params.RoleController.AssignMenuToRole)

			// 为角色设置权限
			roles.POST("/permission", params.RoleController.SetRolePermission)

			// 删除角色权限
			roles.DELETE("/permission", params.RoleController.DeleteRolePermission)

			// 获取角色权限
			roles.GET("/:id/permissions", params.PermissionController.GetRolePermissions)

			// 为角色分配权限
			roles.POST("/:id/permissions", params.PermissionController.AssignPermissionsToRole)
		}

		// 菜单管理
		menus := admin.Group("/menus")
		{
			// 获取菜单详情
			menus.GET("/:id", params.MenuController.GetMenu)

			// 获取菜单树
			menus.GET("/tree", params.MenuController.GetMenuTree)

			// 获取用户菜单
			menus.GET("/user", params.MenuController.GetUserMenus)

			// 获取用户菜单权限
			menus.GET("/user/perms", params.MenuController.GetUserMenuPerms)

			// 创建菜单
			menus.POST("", params.MenuController.CreateMenu)

			// 更新菜单
			menus.PUT("/:id", params.MenuController.UpdateMenu)

			// 删除菜单
			menus.DELETE("/:id", params.MenuController.DeleteMenu)
		}

		// 权限管理
		permissions := admin.Group("/permissions")
		{
			// 获取权限列表
			permissions.GET("", params.PermissionController.GetPermissions)

			// 获取权限详情
			permissions.GET("/:id", params.PermissionController.GetPermission)

			// 获取用户权限
			permissions.GET("/user", params.PermissionController.GetUserPermissions)

			// 获取用户菜单
			permissions.GET("/user/menus", params.PermissionController.GetUserMenus)

			// 获取用户角色
			permissions.GET("/user/roles", params.PermissionController.GetUserRoles)

			// 为用户分配角色
			permissions.POST("/user/:id/roles", params.PermissionController.AssignRolesToUser)

			// 创建权限
			permissions.POST("", params.PermissionController.CreatePermission)

			// 更新权限
			permissions.PUT("/:id", params.PermissionController.UpdatePermission)

			// 删除权限
			permissions.DELETE("/:id", params.PermissionController.DeletePermission)

			// 更新权限系统设置
			permissions.PUT("/settings", params.PermissionController.UpdatePermissionSettings)
		}

		// 系统设置
		system := admin.Group("/system")
		{
			// 获取系统设置
			system.GET("/settings", params.AdminController.GetSystemSettings)
		}
	}
}
