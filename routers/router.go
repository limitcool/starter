package routers

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/internal/storage/casbin"
)

// Load loads the middlewares, routes, handlers.
func NewRouter() *gin.Engine {
	// 创建不带默认中间件的路由
	r := gin.New()

	// 使用我们的结构化日志中间件
	r.Use(middleware.LoggerWithCharmbracelet())

	// 使用恢复中间件
	r.Use(gin.Recovery())
	// 使用CORS中间件
	r.Use(middleware.Cors())

	// 只有在启用权限系统时才初始化Casbin组件
	config := services.Instance().GetConfig()
	if config.Permission.Enabled {
		// 初始化Casbin组件
		casbinComponent := casbin.NewComponent(config)
		if err := casbinComponent.Initialize(); err != nil {
			panic("Casbin组件初始化失败: " + err.Error())
		}
	}

	// v1 router
	apiV1 := r.Group("/api/v1")

	// 公共路由
	{
		apiV1.GET("/ping", controller.Ping)
		apiV1.POST("/admin/login", controller.AdminLogin) // 管理员登录接口
		apiV1.POST("/refresh", controller.RefreshToken)   // 刷新访问令牌接口

		// 普通用户公共路由
		apiV1.POST("/user/register", controller.UserRegister) // 用户注册
		apiV1.POST("/user/login", controller.UserLogin)       // 用户登录
	}

	// 需要认证的路由
	auth := apiV1.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		// 获取当前用户菜单
		auth.GET("/user/menus", controller.GetUserMenus)
		// 获取当前用户权限标识
		auth.GET("/user/perms", controller.GetUserMenuPerms)

		// 普通用户需要认证的路由
		user := auth.Group("/user")
		{
			// 普通用户信息
			user.GET("/info", middleware.RequireNormalUser(), controller.UserInfo)
			// 修改密码
			user.POST("/change-password", middleware.RequireNormalUser(), controller.UserChangePassword)
		}
	}

	// 只有在启用权限系统时才注册需要权限控制的路由
	if config.Permission.Enabled {
		// 需要权限控制的路由
		admin := auth.Group("/admin")
		admin.Use(middleware.CasbinComponentMiddleware())
		{
			// 系统设置
			system := admin.Group("/system")
			{
				system.GET("/settings", controller.GetSystemSettings)
				system.PUT("/permission", controller.UpdatePermissionSettings)
			}

			// 菜单管理
			menu := admin.Group("/menu")
			{
				menu.POST("", controller.CreateMenu)
				menu.PUT("/:id", controller.UpdateMenu)
				menu.DELETE("/:id", controller.DeleteMenu)
				menu.GET("/:id", controller.GetMenu)
				menu.GET("", controller.GetMenuTree)
			}

			// 角色管理
			role := admin.Group("/role")
			{
				role.POST("", controller.CreateRole)
				role.PUT("/:id", controller.UpdateRole)
				role.DELETE("/:id", controller.DeleteRole)
				role.GET("/:id", controller.GetRole)
				role.GET("", controller.GetRoles)
				// 为角色分配菜单
				role.POST("/menu", controller.AssignMenuToRole)
				// 为角色设置权限
				role.POST("/permission", controller.SetRolePermission)
				// 删除角色权限
				role.DELETE("/permission", controller.DeleteRolePermission)
			}

			// 权限管理
			permission := admin.Group("/permission")
			{
				permission.GET("", controller.GetPermissions)
				permission.GET("/:id", controller.GetPermission)
				permission.POST("", controller.CreatePermission)
				permission.PUT("/:id", controller.UpdatePermission)
				permission.DELETE("/:id", controller.DeletePermission)
			}

			// 操作日志管理
			oplog := admin.Group("/operation-logs")
			{
				oplog.GET("", controller.GetOperationLogs)
				oplog.DELETE("/:id", controller.DeleteOperationLog)
				oplog.DELETE("/batch", controller.BatchDeleteOperationLogs)
			}
		}
	}

	// 打印所有注册的路由
	printRegisteredRoutes(r)

	return r
}

// printRegisteredRoutes 打印所有注册的路由
func printRegisteredRoutes(r *gin.Engine) {
	routes := r.Routes()
	log.Info("Registered routes:")
	for _, route := range routes {
		// 从完整路径中提取包名和处理函数名称
		handlerName := route.Handler
		parts := strings.Split(handlerName, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			if dotIndex := strings.Index(lastPart, "."); dotIndex != -1 {
				handlerName = lastPart
			}
		}
		log.Info("Route", "method", route.Method, "path", route.Path, "handler", handlerName)
	}
}
