package app

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/datastore/sqldb"
	"github.com/limitcool/starter/internal/filestore"
	"github.com/limitcool/starter/internal/handler"
	"github.com/limitcool/starter/internal/pkg/cache"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/router"
	"gorm.io/gorm"
)

// App 应用容器
type App struct {
	config      *configs.Config
	db          *gorm.DB
	redis       *redis.Client
	cache       cache.Cache
	storage     *filestore.Storage
	handlers    *Handlers
	router      *gin.Engine
	server      *http.Server
	pprofServer *http.Server // pprof服务器
}

// Handlers 处理器集合
type Handlers struct {
	User  *handler.UserHandler
	File  *handler.FileHandler
	Admin *handler.AdminHandler
}

// New 创建新的应用实例
func New(config *configs.Config) (*App, error) {
	app := &App{
		config: config,
	}

	// 按顺序初始化各个组件
	if err := app.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := app.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	if err := app.initStorage(); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	if err := app.initHandlers(); err != nil {
		return nil, fmt.Errorf("failed to initialize handlers: %w", err)
	}

	if err := app.initRouter(); err != nil {
		return nil, fmt.Errorf("failed to initialize router: %w", err)
	}

	if err := app.initServer(); err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	if err := app.initPprof(); err != nil {
		return nil, fmt.Errorf("failed to initialize pprof: %w", err)
	}

	return app, nil
}

// initDatabase 初始化数据库连接
func (a *App) initDatabase() error {
	if !a.config.Database.Enabled {
		logger.Info("Database disabled")
		return nil
	}

	logger.Info("Connecting to database", "driver", a.config.Driver)

	db := sqldb.NewDBWithConfig(*a.config)
	if db == nil {
		return fmt.Errorf("failed to create database connection")
	}

	// 检查数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	a.db = db
	logger.Info("Database connected successfully")
	return nil
}

// initRedis 初始化Redis连接
func (a *App) initRedis() error {
	// 检查Redis配置
	if len(a.config.Redis.Instances) == 0 {
		logger.Info("Redis disabled - no instances configured")
		return nil
	}

	// 获取默认Redis实例配置
	var redisConfig configs.RedisInstance
	var exists bool
	if redisConfig, exists = a.config.Redis.Instances["default"]; !exists || !redisConfig.Enabled {
		logger.Info("Redis disabled - default instance not enabled")
		return nil
	}

	logger.Info("Connecting to Redis", "addr", redisConfig.Addr)

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         redisConfig.Addr,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		MinIdleConns: redisConfig.MinIdleConn,
		DialTimeout:  redisConfig.DialTimeout,
		ReadTimeout:  redisConfig.ReadTimeout,
		WriteTimeout: redisConfig.WriteTimeout,
		PoolSize:     redisConfig.PoolSize,
		PoolTimeout:  redisConfig.PoolTimeout,
	})

	// 检查Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.Info("Redis connected successfully")

	// 创建Redis缓存
	redisCache := cache.NewRedisCache(
		client,
		cache.WithExpiration(a.config.Redis.Cache.DefaultTTL),
		cache.WithKeyPrefix(a.config.Redis.Cache.KeyPrefix),
	)

	a.redis = client
	a.cache = redisCache
	return nil
}

// initStorage 初始化文件存储
func (a *App) initStorage() error {
	storage, err := filestore.NewWithConfig(*a.config)
	if err != nil {
		return fmt.Errorf("failed to create storage service: %w", err)
	}

	a.storage = storage
	logger.Info("Storage service initialized successfully")
	return nil
}

// initHandlers 初始化处理器
func (a *App) initHandlers() error {
	a.handlers = &Handlers{
		User:  handler.NewUserHandler(a.db, a.config),
		File:  handler.NewFileHandler(a.db, a.config, a.storage),
		Admin: handler.NewAdminHandler(a.db, a.config),
	}

	logger.Info("Handlers initialized successfully")
	return nil
}

// initRouter 初始化路由
func (a *App) initRouter() error {
	r, err := router.NewRouter(a.config, a.handlers.User, a.handlers.File, a.handlers.Admin)
	if err != nil {
		return fmt.Errorf("failed to create router: %w", err)
	}

	a.router = r
	logger.Info("Router initialized successfully")
	return nil
}

// initServer 初始化HTTP服务器
func (a *App) initServer() error {
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.App.Port),
		Handler: a.router,
	}

	logger.Info("HTTP server initialized successfully")
	return nil
}

// initPprof 初始化pprof服务器
func (a *App) initPprof() error {
	if !a.config.Pprof.Enabled {
		logger.Info("Pprof disabled")
		return nil
	}

	// 如果端口为0，则在主服务器上启用pprof
	if a.config.Pprof.Port == 0 {
		logger.Info("Pprof enabled on main server", "path", "/debug/pprof/")
		return nil
	}

	// 创建独立的pprof服务器
	pprofMux := http.NewServeMux()
	pprofMux.Handle("/debug/pprof/", http.DefaultServeMux)

	a.pprofServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.Pprof.Port),
		Handler: pprofMux,
	}

	logger.Info("Pprof server initialized successfully", "port", a.config.Pprof.Port)
	return nil
}

// Run 运行应用
func (a *App) Run() error {
	// 启动pprof服务器（如果配置了独立端口）
	if a.pprofServer != nil {
		go func() {
			logger.Info("Pprof server started", "address", fmt.Sprintf("http://localhost:%d/debug/pprof/", a.config.Pprof.Port))
			if err := a.pprofServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("Pprof server error", "error", err)
			}
		}()
	}

	// 启动HTTP服务器
	go func() {
		logger.Info("==================================================")
		logger.Info("HTTP服务器已启动",
			"address", fmt.Sprintf("http://localhost:%d", a.config.App.Port),
			"mode", a.config.App.Mode)
		if a.config.Pprof.Enabled && a.config.Pprof.Port == 0 {
			logger.Info("Pprof enabled", "address", fmt.Sprintf("http://localhost:%d/debug/pprof/", a.config.App.Port))
		}
		logger.Info("==================================================")

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭
	return a.Shutdown()
}

// Shutdown 优雅关闭应用
func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭pprof服务器
	if a.pprofServer != nil {
		if err := a.pprofServer.Shutdown(ctx); err != nil {
			logger.Error("Pprof server forced to shutdown", "error", err)
		} else {
			logger.Info("Pprof server stopped")
		}
	}

	// 关闭HTTP服务器
	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Error("Server forced to shutdown", "error", err)
		}
	}

	// 关闭数据库连接
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Error("Failed to close database connection", "error", err)
			} else {
				logger.Info("Database connection closed")
			}
		}
	}

	// 关闭Redis连接
	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			logger.Error("Failed to close Redis connection", "error", err)
		} else {
			logger.Info("Redis connection closed")
		}
	}

	logger.Info("Application stopped")
	return nil
}
