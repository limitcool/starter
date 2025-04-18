package routers

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/storage"
	"github.com/limitcool/starter/internal/repository"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/internal/storage/casbin"
	"github.com/limitcool/starter/internal/storage/database"
)

// NewRouter 初始化并返回一个配置完整的Gin路由引擎
func NewRouter(db database.DB) *gin.Engine {
	// 创建不带默认中间件的路由
	r := gin.New()

	// 添加全局中间件
	r.Use(middleware.ErrorHandler()) // 添加全局错误处理中间件
	r.Use(middleware.LoggerWithCharmbracelet())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.I18n())

	// 获取服务配置
	config := core.Instance().Config()

	// 创建Casbin组件
	var casbinComponent *casbin.Component
	if config.Casbin.Enabled {
		// 创建Casbin组件
		casbinComponent = casbin.NewComponent(config, db.GetDB())
		// 初始化
		if err := casbinComponent.Initialize(); err != nil {
			log.Warn("初始化Casbin组件失败", "err", err)
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
			log.Error("Failed to initialize storage service", "err", err)
		} else {
			log.Info("Storage service initialized successfully", "type", config.Storage.Type)
		}
	}

	// 配置静态文件服务
	if stg != nil && config.Storage.Type == storage.StorageTypeLocal {
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

		log.Info("Configuring local static file service", "path", config.Storage.Local.Path, "url_path", urlPath)
		// 使用StaticFS提供静态文件服务
		r.StaticFS(urlPath, http.Dir(config.Storage.Local.Path))
	}

	// 使用传入的数据库实例创建服务

	// 创建仓库实例
	gormDB := db.GetDB()
	menuRepo := repository.NewMenuRepo(gormDB)
	roleRepo := repository.NewRoleRepo(gormDB)
	sysUserRepo := repository.NewSysUserRepo(gormDB)
	permissionRepo := repository.NewPermissionRepo(gormDB)
	operationLogRepo := repository.NewOperationLogRepo(gormDB)
	userRepo := repository.NewUserRepository(gormDB)
	fileRepo := repository.NewFileRepo(gormDB)

	// 创建服务实例
	casbinService := services.NewCasbinService(db)
	roleService := services.NewRoleService(roleRepo, casbinService)
	menuService := services.NewMenuService(menuRepo, casbinService)
	sysUserService := services.NewSysUserService(sysUserRepo, userRepo, roleService)
	permissionService := services.NewPermissionService(permissionRepo)
	operationLogService := services.NewOperationLogService(operationLogRepo)
	systemService := services.NewSystemService(db)
	userService := services.NewUserService(userRepo)

	// 初始化控制器
	userController := controller.NewUserController(sysUserService, userService)
	adminController := controller.NewAdminController(sysUserService)
	roleController := controller.NewRoleController(roleService, menuService)
	menuController := controller.NewMenuController(menuService)
	permissionController := controller.NewPermissionController(permissionService)
	operationLogController := controller.NewOperationLogController(operationLogService)
	systemController := controller.NewSystemController(systemService)

	var fileController *controller.FileController
	if stg != nil {
		fileController = controller.NewFileController(stg, fileRepo)
	}

	// 设置API路由组
	apiV1 := r.Group("/api/v1")

	// 健康检查路由
	apiV1.GET("/ping", controller.Ping)

	// 令牌刷新（公共API）
	apiV1.POST("/auth/tokens/refresh", userController.RefreshToken)

	// 用户相关 - 公开路由
	publicUser := apiV1.Group("/user")
	{
		publicUser.POST("/register", userController.UserRegister)
		publicUser.POST("/login", userController.UserLogin)
	}

	// 管理员相关 - 公开路由
	publicAdmin := apiV1.Group("/admin")
	{
		publicAdmin.POST("/login", adminController.AdminLogin)
	}

	// ======= 文件下载API（无需认证）=======
	if fileController != nil {
		apiV1.GET("/files/download/:id", fileController.DownloadFile)
	}

	// ======= 需要认证的路由 =======
	auth := apiV1.Group("")
	auth.Use(middleware.JWTAuth())

	// 用户相关 - 需要认证
	authUser := auth.Group("/user")
	authUser.Use(middleware.AuthNormalUser())
	{
		authUser.GET("", userController.UserInfo)
		authUser.PUT("/password", userController.UserChangePassword)
		authUser.GET("/menus", menuController.GetUserMenus)
		authUser.GET("/permissions", menuController.GetUserMenuPerms)

		// 用户头像上传
		if fileController != nil {
			authUser.PUT("/avatar", fileController.UpdateUserAvatar)
		}
	}

	// 文件相关 - 需要认证
	if fileController != nil {
		authFiles := auth.Group("/files")
		{
			authFiles.POST("/upload", fileController.UploadFile)
			authFiles.GET("/:id", fileController.GetFile)
			authFiles.DELETE("/:id", fileController.DeleteFile)
		}
	}

	// 管理员路由 - 需要权限验证
	if config.Casbin.Enabled && casbinComponent != nil {
		adminGroup := auth.Group("/admin")
		adminGroup.Use(middleware.CasbinMiddleware(casbinService))
		{
			// 系统管理
			adminSystem := adminGroup.Group("/system")
			{
				adminSystem.GET("", systemController.GetSystemSettings)
				adminSystem.PUT("/permissions", permissionController.UpdatePermissionSettings)
			}

			// 菜单管理
			adminMenu := adminGroup.Group("/menu")
			{
				adminMenu.GET("", menuController.GetMenuTree)
				adminMenu.GET("/:id", menuController.GetMenu)
				adminMenu.POST("", menuController.CreateMenu)
				adminMenu.PUT("/:id", menuController.UpdateMenu)
				adminMenu.DELETE("/:id", menuController.DeleteMenu)
			}

			// 角色管理
			adminRole := adminGroup.Group("/role")
			{
				adminRole.GET("", roleController.GetRoles)
				adminRole.GET("/:id", roleController.GetRole)
				adminRole.POST("", roleController.CreateRole)
				adminRole.PUT("/:id", roleController.UpdateRole)
				adminRole.DELETE("/:id", roleController.DeleteRole)
				adminRole.POST("/:id/menus", roleController.AssignMenuToRole)
				adminRole.POST("/:id/permissions", roleController.SetRolePermission)
				adminRole.DELETE("/:id/permissions", roleController.DeleteRolePermission)
			}

			// 权限管理
			adminPermission := adminGroup.Group("/permission")
			{
				adminPermission.GET("", permissionController.GetPermissions)
				adminPermission.GET("/:id", permissionController.GetPermission)
				adminPermission.POST("", permissionController.CreatePermission)
				adminPermission.PUT("/:id", permissionController.UpdatePermission)
				adminPermission.DELETE("/:id", permissionController.DeletePermission)
			}

			// 操作日志管理
			adminLog := adminGroup.Group("/log")
			{
				adminLog.GET("", operationLogController.GetOperationLogs)
				adminLog.DELETE("/:id", operationLogController.DeleteOperationLog)
				adminLog.POST("/batch-delete", operationLogController.ClearOperationLogs)
			}

			// 系统用户头像管理
			if fileController != nil {
				adminUser := adminGroup.Group("/user")
				{
					adminUser.PUT("/avatar", fileController.UpdateSysUserAvatar)
					adminUser.PUT("/avatar/:id", fileController.UpdateSysUserAvatar)
				}
			}
		}
	}

	// 打印所有注册的路由
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

	return r
}
