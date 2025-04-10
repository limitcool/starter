package routers

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/storage"
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
	config := core.Instance().Config()

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
	AdminControllerInstance := controller.NewAdminController()
	// 设置API路由组
	apiV1 := r.Group("/api/v1")

	// 公共路由
	{
		apiV1.GET("/ping", controller.Ping)

		// 认证相关
		auth := apiV1.Group("/auth")
		{
			auth.POST("/admin/login", AdminControllerInstance.AdminLogin)
			auth.POST("/tokens/refresh", UserControllerInstance.RefreshToken)
			// 普通用户认证
			auth.POST("/users/register", UserControllerInstance.UserRegister)
			auth.POST("/users/login", UserControllerInstance.UserLogin)
		}
	}

	// 需要认证的路由
	authRequired := apiV1.Group("")
	authRequired.Use(middleware.JWTAuth())
	{
		// 用户相关
		users := authRequired.Group("/users")
		{
			users.GET("/menus", MenuControllerInstance.GetUserMenus)
			users.GET("/permissions", MenuControllerInstance.GetUserMenuPerms)                                     // 更改为复数形式
			users.GET("/current", middleware.AuthNormalUser(), UserControllerInstance.UserInfo)                    // 使用/current表示当前用户
			users.PUT("/current/password", middleware.AuthNormalUser(), UserControllerInstance.UserChangePassword) // 更改为更符合RESTful的形式
		}

		// 管理员权限路由（如果启用了权限系统）
		if config.Permission.Enabled && casbinComponent != nil {
			admin := authRequired.Group("")
			admin.Use(middleware.CasbinComponentMiddleware())
			{
				// 系统设置
				systems := admin.Group("/systems")
				{
					systems.GET("", SystemControllerInstance.GetSystemSettings)
					systems.PUT("/permissions", PermissionControllerInstance.UpdatePermissionSettings)
				}

				// 菜单管理
				menus := admin.Group("/menus")
				{
					menus.POST("", MenuControllerInstance.CreateMenu)
					menus.PUT("/:id", MenuControllerInstance.UpdateMenu)
					menus.DELETE("/:id", MenuControllerInstance.DeleteMenu)
					menus.GET("/:id", MenuControllerInstance.GetMenu)
					menus.GET("", MenuControllerInstance.GetMenuTree)
				}

				// 角色管理
				roles := admin.Group("/roles")
				{
					roles.POST("", RoleControllerInstance.CreateRole)
					roles.PUT("/:id", RoleControllerInstance.UpdateRole)
					roles.DELETE("/:id", RoleControllerInstance.DeleteRole)
					roles.GET("/:id", RoleControllerInstance.GetRole)
					roles.GET("", RoleControllerInstance.GetRoles)
					roles.POST("/:id/menus", RoleControllerInstance.AssignMenuToRole)             // 使用子资源路径
					roles.POST("/:id/permissions", RoleControllerInstance.SetRolePermission)      // 使用子资源路径
					roles.DELETE("/:id/permissions", RoleControllerInstance.DeleteRolePermission) // 使用子资源路径
				}

				// 权限管理
				permissions := admin.Group("/permissions")
				{
					permissions.GET("", PermissionControllerInstance.GetPermissions)
					permissions.GET("/:id", PermissionControllerInstance.GetPermission)
					permissions.POST("", PermissionControllerInstance.CreatePermission)
					permissions.PUT("/:id", PermissionControllerInstance.UpdatePermission)
					permissions.DELETE("/:id", PermissionControllerInstance.DeletePermission)
				}

				// 操作日志管理 - 更改为复数形式
				operationLogs := admin.Group("/operation-logs")
				{
					operationLogs.GET("", OperationLogControllerInstance.GetOperationLogs)
					operationLogs.DELETE("/:id", OperationLogControllerInstance.DeleteOperationLog)
					// 使用POST方法进行批量删除，更符合请求体传递ID列表的设计
					operationLogs.POST("/batch-delete", OperationLogControllerInstance.ClearOperationLogs)
				}
			}
		}
	}

	// 如果存储服务可用，设置文件相关路由
	if stg != nil {
		fileController := controller.NewFileController(stg)

		// 文件资源
		files := apiV1.Group("/files")
		{
			// 上传文件需要登录
			files.POST("", middleware.JWTAuth(), fileController.UploadFile)

			// 获取文件信息
			files.GET("/:id", fileController.GetFile)

			// 下载文件
			files.GET("/:id/download", fileController.DownloadFile)

			// 删除文件需要登录
			files.DELETE("/:id", middleware.JWTAuth(), fileController.DeleteFile)
		}

		// 用户头像
		apiV1.PUT("/users/:id/avatar", middleware.JWTAuth(), fileController.UpdateUserAvatar)
		apiV1.PUT("/users/current/avatar", middleware.JWTAuth(), fileController.UpdateUserAvatar)

		// 系统用户头像
		apiV1.PUT("/system-users/:id/avatar", middleware.JWTAuth(), fileController.UpdateSysUserAvatar)
		apiV1.PUT("/system-users/current/avatar", middleware.JWTAuth(), fileController.UpdateSysUserAvatar)
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
