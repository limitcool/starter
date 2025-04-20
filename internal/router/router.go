package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/filestore"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/types"
)

// SetupRouter 初始化并返回一个配置完整的Gin路由引擎
// 注意：这个函数已经被 NewRouter 函数替代，保留它是为了兼容旧代码
func SetupRouter(params RouterParams) *gin.Engine {
	// 创建不带默认中间件的路由
	r := gin.New()

	// 添加全局中间件
	r.Use(middleware.RequestContext()) // 添加请求上下文中间件
	r.Use(middleware.ErrorHandler())   // 添加全局错误处理中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.I18n(params.Config))

	// 使用依赖注入传入的配置
	config := params.Config

	// 创建Casbin服务
	if config.Casbin.Enabled && params.Enforcer != nil {
		logger.Info("Casbin服务已启用")
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

	// 注册路由
	registerRoutes(r, params)

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
