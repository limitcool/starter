package main

import (
	"context"
	"fmt"
	"io"

	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/database"
	"github.com/limitcool/starter/internal/database/mongodb"
	"github.com/limitcool/starter/internal/storage/redisdb"
	"github.com/limitcool/starter/pkg/env"
	"github.com/limitcool/starter/pkg/logger"
	"github.com/limitcool/starter/routers"
	"github.com/spf13/viper"
)

func loadConfig() *configs.Config {
	env := env.Get()

	// 直接读取环境对应的配置文件
	configName := env.String() // 使用环境名称作为配置文件名: dev.yaml, test.yaml, prod.yaml

	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")

	// 读取环境配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config file", "env", env, "error", err)
	}

	// 解析配置到结构体
	cfg := &configs.Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal("Config unmarshal failed", "error", err)
	}

	// 配置日志系统
	logger.Setup(cfg.Log)

	// 记录环境信息
	log.Info("Environment configured", "env", env)

	// 设置全局配置
	global.Config = cfg

	return cfg
}

func main() {
	// 设置基本日志
	cfg := loadConfig()
	log.SetPrefix("🌏 starter ")

	// 获取环境
	currentEnv := env.Get()

	// 根据环境设置Gin模式
	if currentEnv == env.Dev {
		// 在开发环境中，我们可以保留Gin的调试输出
		gin.SetMode(gin.DebugMode)

		// 但仍然将它重定向到我们的日志系统
		logger.SetupGinLogger()
	} else {
		// 在非开发环境中，完全禁用Gin的调试输出
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
	}

	// 即使在开发环境中，也可以选择禁用Gin的调试日志
	gin.DisableConsoleColor()

	// 使用配置更新日志设置
	logger.Setup(cfg.Log)

	// 日志系统配置完成后的第一条日志
	log.Info("Application starting", "name", cfg.App.Name)

	// 根据环境设置Gin模式
	if env.IsProd() {
		log.Info("Running in production mode")
	} else if env.IsTest() {
		log.Info("Running in test mode")
	} else {
		log.Info("Running in debug mode")
	}

	// 初始化应用核心
	app := core.Setup(cfg)

	// 初始化数据库
	initDatabase(cfg)

	// 添加Redis组件
	redisComponent := redisdb.NewComponent(cfg)
	app.ComponentManager.AddComponent(redisComponent)

	// 初始化所有组件
	if err := app.Initialize(); err != nil {
		log.Fatal("Failed to initialize application", "error", err)
	}

	// 确保资源清理
	defer app.Cleanup()

	// 初始化路由
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           fmt.Sprint("0.0.0.0:", cfg.App.Port),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}
	log.Info("Server started", "url", fmt.Sprintf("http://127.0.0.1:%d", cfg.App.Port))
	go func() {
		// 服务连接 监听
		if err := s.ListenAndServe(); err != nil {
			log.Fatal("Server listening failed", "error", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器,这里需要缓冲
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	//(设置5秒超时时间)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		// 处理错误，例如记录日志、返回错误等
		log.Info("Error during server shutdown", "error", err)
	}
}

// initDatabase 初始化数据库
func initDatabase(cfg *configs.Config) {
	switch cfg.Driver {
	case configs.DriverMongo:
		log.Info("Using database driver", "driver", "mongo")
		_, err := mongodb.NewMongoDBConn(context.Background(), &cfg.Mongo)
		if err != nil {
			log.Fatal("MongoDB connection failed", "error", err)
		}
	case configs.DriverMysql, configs.DriverPostgres, configs.DriverSqlite, configs.DriverMssql, configs.DriverOracle:
		log.Info("Using database driver", "driver", cfg.Driver)
		db := database.NewDB(*cfg)
		db.AutoMigrate()
	default:
		log.Fatal("No database driver", "driver", "none")
	}
}
