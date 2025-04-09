package routers

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/storage"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/internal/storage/casbin"
)

// NewRouter 初始化并返回一个配置完整的Gin路由引擎
func NewRouter() *gin.Engine {
	// 创建不带默认中间件的路由
	r := gin.New()

	// 添加全局中间件
	r.Use(middleware.LoggerWithCharmbracelet())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// 获取服务配置
	config := services.Instance().GetConfig()

	// 初始化Casbin权限系统（如果启用）
	var casbinComponent *casbin.Component
	if config.Permission.Enabled {
		casbinComponent = casbin.NewComponent(config)
		if err := casbinComponent.Initialize(); err != nil {
			panic("Casbin组件初始化失败: " + err.Error())
		}
	}

	// 初始化存储服务（如果启用）
	var stg *storage.Storage
	if config.Storage.Enabled {
		storageConfig := storage.Config{Type: config.Storage.Type}

		// 根据存储类型设置配置
		switch config.Storage.Type {
		case storage.StorageTypeLocal:
			storageConfig.Path = config.Storage.Local.Path
			storageConfig.URL = config.Storage.Local.URL
		case storage.StorageTypeS3:
			storageConfig.AccessKey = config.Storage.S3.AccessKey
			storageConfig.SecretKey = config.Storage.S3.SecretKey
			storageConfig.Region = config.Storage.S3.Region
			storageConfig.Bucket = config.Storage.S3.Bucket
			storageConfig.Endpoint = config.Storage.S3.Endpoint
		}

		var err error
		stg, err = storage.New(storageConfig)
		if err != nil {
			log.Error("初始化存储服务失败", "err", err)
		} else {
			log.Info("存储服务初始化成功", "type", config.Storage.Type)
		}
	}
	UserControllerInstance := controller.NewUserController()
	PermissionControllerInstance := controller.NewPermissionController()
	OperationLogControllerInstance := controller.NewOperationLogController()
	MenuControllerInstance := controller.NewMenuController()
	RoleControllerInstance := controller.NewRoleController()
	SystemControllerInstance := controller.NewSystemController()
	// 设置API路由组
	api := r.Group("/api")
	apiV1 := api.Group("/v1")

	// 公共路由
	{
		apiV1.GET("/ping", controller.Ping)
		apiV1.POST("/admin/login", controller.AdminControllerInstance.AdminLogin)
		apiV1.POST("/refresh", UserControllerInstance.RefreshToken)

		// 普通用户公共路由
		apiV1.POST("/user/register", UserControllerInstance.UserRegister)
		apiV1.POST("/user/login", UserControllerInstance.UserLogin)
	}

	// 需要认证的路由
	auth := apiV1.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user/menus", MenuControllerInstance.GetUserMenus)
		auth.GET("/user/perms", MenuControllerInstance.GetUserMenuPerms)

		// 普通用户认证路由
		user := auth.Group("/user")
		{
			user.GET("/info", middleware.AuthNormalUser(), UserControllerInstance.UserInfo)
			user.POST("/change-password", middleware.AuthNormalUser(), UserControllerInstance.UserChangePassword)
		}

		// 管理员权限路由（如果启用了权限系统）
		if config.Permission.Enabled && casbinComponent != nil {
			admin := auth.Group("/admin")
			admin.Use(middleware.CasbinComponentMiddleware())
			{
				// 系统设置
				system := admin.Group("/system")
				{
					system.GET("/settings", SystemControllerInstance.GetSystemSettings)
					system.PUT("/permission", PermissionControllerInstance.UpdatePermissionSettings)
				}

				// 菜单管理
				menu := admin.Group("/menu")
				{
					menu.POST("", MenuControllerInstance.CreateMenu)
					menu.PUT("/:id", MenuControllerInstance.UpdateMenu)
					menu.DELETE("/:id", MenuControllerInstance.DeleteMenu)
					menu.GET("/:id", MenuControllerInstance.GetMenu)
					menu.GET("", MenuControllerInstance.GetMenuTree)
				}

				// 角色管理
				role := admin.Group("/role")
				{
					role.POST("", RoleControllerInstance.CreateRole)
					role.PUT("/:id", RoleControllerInstance.UpdateRole)
					role.DELETE("/:id", RoleControllerInstance.DeleteRole)
					role.GET("/:id", RoleControllerInstance.GetRole)
					role.GET("", RoleControllerInstance.GetRoles)
					role.POST("/menu", RoleControllerInstance.AssignMenuToRole)
					role.POST("/permission", RoleControllerInstance.SetRolePermission)
					role.DELETE("/permission", RoleControllerInstance.DeleteRolePermission)
				}

				// 权限管理
				permission := admin.Group("/permission")
				{
					permission.GET("", PermissionControllerInstance.GetPermissions)
					permission.GET("/:id", PermissionControllerInstance.GetPermission)
					permission.POST("", PermissionControllerInstance.CreatePermission)
					permission.PUT("/:id", PermissionControllerInstance.UpdatePermission)
					permission.DELETE("/:id", PermissionControllerInstance.DeletePermission)
				}

				// 操作日志管理
				oplog := admin.Group("/operation-logs")
				{
					oplog.GET("", OperationLogControllerInstance.GetOperationLogs)
					oplog.DELETE("/:id", OperationLogControllerInstance.DeleteOperationLog)
					oplog.DELETE("/batch", OperationLogControllerInstance.ClearOperationLogs)
				}
			}
		}
	}

	// 用户管理
	userAuthGroup := api.Group("/users")
	{
		userAuthGroup.GET("/info", services.GetUserInfo)
		userAuthGroup.PUT("/", middleware.AuthNormalUser(), services.UserRegister)
		userAuthGroup.POST("/register", middleware.AuthNormalUser(), services.UserRegister)
	}

	// 如果存储服务可用，设置文件相关路由
	if stg != nil {
		SetupFileRoutes(api, stg)
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
