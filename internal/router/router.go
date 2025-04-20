package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/datastore/database"
	"github.com/limitcool/starter/internal/filestore"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/types"
)

// SetupRouter 初始化并返回一个配置完整的Gin路由引擎
func SetupRouter(db database.Database, config *configs.Config) *gin.Engine {
	// 创建不带默认中间件的路由
	r := gin.New()

	// 添加全局中间件
	r.Use(middleware.RequestContext()) // 添加请求上下文中间件
	r.Use(middleware.ErrorHandler())   // 添加全局错误处理中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.I18n(config))

	// 使用依赖注入传入的配置

	// 创建Casbin服务
	var casbinService casbin.Service
	if config.Casbin.Enabled {
		// 创建Casbin服务
		casbinService = casbin.NewService(db.DB(), config)
		// 初始化
		if err := casbinService.Initialize(); err != nil {
			logger.Warn("初始化Casbin服务失败", "err", err)
		}
	}

	// 初始化存储服务（如果启用）
	var stg *filestore.Storage
	if config.Storage.Enabled {
		storageConfig := filestore.Config{Type: config.Storage.Type}

		// 根据存储类型设置配置
		switch config.Storage.Type {
		case types.StorageTypeLocal:
			storageConfig.Path = config.Storage.Local.Path
			storageConfig.URL = config.Storage.Local.URL
		case types.StorageTypeS3:
			storageConfig.AccessKey = config.Storage.S3.AccessKey
			storageConfig.SecretKey = config.Storage.S3.SecretKey
			storageConfig.Region = config.Storage.S3.Region
			storageConfig.Bucket = config.Storage.S3.Bucket
			storageConfig.Endpoint = config.Storage.S3.Endpoint
		}

		var err error
		stg, err = filestore.New(storageConfig)
		if err != nil {
			logger.Error("Failed to initialize storage service", "err", err)
		} else {
			logger.Info("Storage service initialized successfully", "type", config.Storage.Type)
		}
	}

	// 配置静态文件服务
	if stg != nil && config.Storage.Type == types.StorageTypeLocal {
		// 从URL提取路径前缀
		urlPath := "/static" // 默认路径
		if config.Storage.Local.URL != "" {
			u := config.Storage.Local.URL
			// 如果URL包含http://或https://，则提取路径部分
			if strings.Contains(u, "://") {
				parts := strings.Split(u, "://")
				if len(parts) > 1 {
					hostPath := strings.Split(parts[1], "/")
					if len(hostPath) > 1 {
						urlPath = "/" + strings.Join(hostPath[1:], "/")
					}
				}
			} else if strings.HasPrefix(u, "/") {
				// 如果URL直接以/开头，则直接使用
				urlPath = u
			}
		}

		logger.Info("Configuring local static file service", "path", config.Storage.Local.Path, "url_path", urlPath)
		// 使用StaticFS提供静态文件服务
		r.StaticFS(urlPath, http.Dir(config.Storage.Local.Path))
	}

	// 使用传入的数据库实例创建服务

	// 获取数据库连接
	gormDB := db.DB()

	// 不再使用旧的Casbin服务
	// 我们已经创建了新的Casbin服务

	// 初始化仓库层
	repos := initRepositories(gormDB)

	// 初始化服务层
	svcs := initServices(repos, casbinService, db, config)

	// 初始化控制器层
	controllers := initControllers(svcs, repos, stg)

	// 设置API路由组
	apiV1 := r.Group("/api/v1")

	// 健康检查路由
	apiV1.GET("/ping", controller.Ping)

	// 令牌刷新（公共API）
	apiV1.POST("/auth/tokens/refresh", controllers.UserController.RefreshToken)

	// 用户相关 - 公开路由
	publicUser := apiV1.Group("/user")
	{
		publicUser.POST("/register", controllers.UserController.UserRegister)
		publicUser.POST("/login", controllers.UserController.UserLogin)
	}

	// 管理员相关 - 公开路由
	publicAdmin := apiV1.Group("/admin")
	{
		publicAdmin.POST("/login", controllers.SysUserController.SysUserLogin)
	}

	// ======= 文件下载API（无需认证）=======
	if controllers.FileController != nil {
		apiV1.GET("/files/download/:id", controllers.FileController.DownloadFile)
	}

	// ======= 需要认证的路由 =======
	auth := apiV1.Group("")
	auth.Use(middleware.JWTAuth(config))

	// 用户相关 - 需要认证
	authUser := auth.Group("/user")
	authUser.Use(middleware.AuthNormalUser())
	{
		authUser.GET("", controllers.UserController.UserInfo)
		authUser.PUT("/password", controllers.UserController.UserChangePassword)
		authUser.GET("/menus", controllers.PermissionController.GetUserMenus)
		authUser.GET("/permissions", controllers.PermissionController.GetUserPermissions)
		authUser.GET("/roles", controllers.PermissionController.GetUserRoles)

		// 用户头像上传
		if controllers.FileController != nil {
			authUser.PUT("/avatar", controllers.FileController.UpdateUserAvatar)
		}
	}

	// 文件相关 - 需要认证
	if controllers.FileController != nil {
		authFiles := auth.Group("/files")
		{
			authFiles.POST("/upload", controllers.FileController.UploadFile)
			authFiles.GET("/:id", controllers.FileController.GetFile)
			authFiles.DELETE("/:id", controllers.FileController.DeleteFile)
		}
	}

	// 管理员路由 - 需要权限验证
	if config.Casbin.Enabled && casbinService != nil {
		adminGroup := auth.Group("/admin")
		adminGroup.Use(middleware.CasbinMiddleware(svcs.PermissionService, config))
		{
			// 系统管理
			adminSystem := adminGroup.Group("/system")
			{
				adminSystem.GET("", controllers.SystemController.GetSystemSettings)
				adminSystem.PUT("/permissions", controllers.PermissionController.UpdatePermissionSettings)
			}

			// 菜单管理
			adminMenu := adminGroup.Group("/menu")
			{
				adminMenu.GET("", controllers.MenuController.GetMenuTree)
				adminMenu.GET("/:id", controllers.MenuController.GetMenu)
				adminMenu.POST("", controllers.MenuController.CreateMenu)
				adminMenu.PUT("/:id", controllers.MenuController.UpdateMenu)
				adminMenu.DELETE("/:id", controllers.MenuController.DeleteMenu)
			}

			// 角色管理
			adminRole := adminGroup.Group("/role")
			{
				adminRole.GET("", controllers.RoleController.GetRoles)
				adminRole.GET("/:id", controllers.RoleController.GetRole)
				adminRole.POST("", controllers.RoleController.CreateRole)
				adminRole.PUT("/:id", controllers.RoleController.UpdateRole)
				adminRole.DELETE("/:id", controllers.RoleController.DeleteRole)
				adminRole.POST("/:id/menus", controllers.RoleController.AssignMenuToRole)
				adminRole.POST("/:id/permissions", controllers.RoleController.SetRolePermission)
				adminRole.DELETE("/:id/permissions", controllers.RoleController.DeleteRolePermission)
			}

			// 权限管理
			adminPermission := adminGroup.Group("/permission")
			{
				adminPermission.GET("", controllers.PermissionController.GetPermissions)
				adminPermission.GET("/:id", controllers.PermissionController.GetPermission)
				adminPermission.POST("", controllers.PermissionController.CreatePermission)
				adminPermission.PUT("/:id", controllers.PermissionController.UpdatePermission)
				adminPermission.DELETE("/:id", controllers.PermissionController.DeletePermission)
				// 角色权限管理
				adminPermission.GET("/role/:id", controllers.PermissionController.GetRolePermissions)
				adminPermission.POST("/role/:id", controllers.PermissionController.AssignPermissionsToRole)
			}

			// API管理
			adminAPI := adminGroup.Group("/api")
			{
				adminAPI.GET("", controllers.APIController.GetAPIs)
				adminAPI.GET("/:id", controllers.APIController.GetAPI)
				adminAPI.POST("", controllers.APIController.CreateAPI)
				adminAPI.PUT("/:id", controllers.APIController.UpdateAPI)
				adminAPI.DELETE("/:id", controllers.APIController.DeleteAPI)
				// 菜单API关联
				adminAPI.GET("/menu/:id", controllers.APIController.GetMenuAPIs)
				adminAPI.POST("/menu/:id", controllers.APIController.AssignAPIsToMenu)
				// 同步菜单API权限
				adminAPI.POST("/sync", controllers.APIController.SyncMenuAPIPermissions)
			}

			// 操作日志管理
			adminLog := adminGroup.Group("/log")
			{
				adminLog.GET("", controllers.OperationLogController.GetOperationLogs)
				adminLog.DELETE("/:id", controllers.OperationLogController.DeleteOperationLog)
				adminLog.POST("/batch-delete", controllers.OperationLogController.ClearOperationLogs)
			}

			// 系统用户管理
			adminUser := adminGroup.Group("/user")
			{
				// 用户角色管理
				adminUser.POST("/:id/roles", controllers.PermissionController.AssignRolesToUser)

				// 系统用户头像管理
				if controllers.FileController != nil {
					adminUser.PUT("/avatar", controllers.FileController.UpdateSysUserAvatar)
					adminUser.PUT("/avatar/:id", controllers.FileController.UpdateSysUserAvatar)
				}
			}
		}
	}

	// 打印所有注册的路由
	routes := r.Routes()
	logger.Info("Registered routes:")
	for _, route := range routes {
		handlerName := route.Handler
		parts := strings.Split(handlerName, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			if dotIndex := strings.Index(lastPart, "."); dotIndex != -1 {
				handlerName = lastPart
			}
		}
		logger.Info("Route", "method", route.Method, "path", route.Path, "handler", handlerName)
	}

	return r
}
