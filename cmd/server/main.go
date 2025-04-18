package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/router"
	"github.com/limitcool/starter/internal/storage/casbin"
	"github.com/limitcool/starter/internal/storage/database"
	"github.com/limitcool/starter/internal/storage/mongodb"
	"github.com/limitcool/starter/internal/storage/redisdb"
	"github.com/limitcool/starter/internal/storage/sqldb"
)

func main() {
	// 加载配置
	cfg := configs.LoadConfig()

	// 设置日志
	setupLogger(cfg)

	// 日志系统配置完成后的第一条日志
	log.Info("Application starting", "name", cfg.App.Name)

	// 初始化应用核心
	app := core.Setup(cfg)

	// 添加SQL数据库组件（如果配置了启用）
	if cfg.Database.Enabled {
		log.Info("Adding SQL database component", "driver", cfg.Driver)
		dbComponent := sqldb.NewComponent(cfg)
		app.ComponentManager.AddComponent(dbComponent)
	}

	// 添加MongoDB组件（如果配置了启用）
	if cfg.Mongo.Enabled {
		log.Info("Adding MongoDB component")
		mongoComponent := mongodb.NewComponent(cfg)
		app.ComponentManager.AddComponent(mongoComponent)
	}

	// 添加Redis组件（如果配置了启用）
	redisComponent := redisdb.NewComponent(cfg)
	app.ComponentManager.AddComponent(redisComponent)

	// 添加Casbin组件（如果配置了启用）
	if cfg.Casbin.Enabled {
		log.Info("Adding Casbin component")
		// 获取数据库组件
		var dbComponent *sqldb.Component
		for _, component := range app.ComponentManager.GetComponents() {
			if c, ok := component.(*sqldb.Component); ok {
				dbComponent = c
				break
			}
		}

		if dbComponent != nil {
			casbinComponent := casbin.NewComponent(cfg, dbComponent.DB())
			app.ComponentManager.AddComponent(casbinComponent)
		} else {
			log.Warn("Cannot add Casbin component: database component not found")
		}
	}

	// 初始化所有组件
	if err := app.Initialize(); err != nil {
		log.Fatal("Failed to initialize application", "error", err)
	}

	// 获取数据库组件
	var dbComponent *sqldb.Component
	for _, component := range app.ComponentManager.GetComponents() {
		if c, ok := component.(*sqldb.Component); ok {
			dbComponent = c
			break
		}
	}

	if dbComponent == nil {
		log.Fatal("Failed to get database component")
	}

	// 创建数据库适配器
	db := database.NewGormDB(dbComponent.DB())

	// 确保资源清理
	defer app.Cleanup()
	
	// 初始化路由
	r := router.SetupRouter(db)
	
	s := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%d", cfg.App.Port),
		Handler:        r,
		MaxHeaderBytes: 1 << 20,
	}
	
	log.Info("Server started", "url", fmt.Sprintf("http://127.0.0.1:%d", cfg.App.Port))
	
	go func() {
		// 服务连接监听
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server listening failed", "error", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器,这里需要缓冲
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	log.Info("Shutting down server...")
	if err := s.Shutdown(ctx); err != nil {
		// 处理错误，例如记录日志、返回错误等
		log.Info("Error during server shutdown", "error", err)
	}

	log.Info("Server stopped gracefully")
}

// setupLogger 设置日志系统
func setupLogger(cfg *configs.Config) {
	// 设置日志级别
	switch cfg.Log.Level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// 设置日志格式
	if cfg.Log.Format == "json" {
		log.SetFormatter(log.JSONFormatter)
	}
}
