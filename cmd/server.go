package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/router"
	"github.com/limitcool/starter/internal/storage/database"
	"github.com/limitcool/starter/internal/storage/mongodb"
	"github.com/limitcool/starter/internal/storage/redisdb"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"github.com/spf13/cobra"
)

// serverCmd 表示server子命令
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server",
	Long: `Start HTTP server and provide Web API services.

The server will load the configuration file and initialize all necessary components, including database connections, logging systems, etc.
The server gracefully handles shutdown signals, ensuring all requests are processed and resources are safely closed.`,
	Run: runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// 添加服务器特定的标志
	serverCmd.Flags().IntP("port", "p", 0, "HTTP server port number, overrides the setting in the configuration file")
}

// runServer 运行HTTP服务器
func runServer(cmd *cobra.Command, args []string) {
	// 加载配置
	cfg := InitConfig(cmd, args)

	// 设置日志
	InitLogger(cfg)

	// 检查是否从命令行指定了端口
	port, _ := cmd.Flags().GetInt("port")
	if port > 0 {
		cfg.App.Port = port
	}

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

	// 初始化MongoDB集合
	if cfg.Mongo.Enabled {
		// 获取MongoDB组件
		var mongoComponent *mongodb.Component
		for _, component := range app.ComponentManager.GetComponents() {
			if c, ok := component.(*mongodb.Component); ok {
				mongoComponent = c
				break
			}
		}

		if mongoComponent != nil {
			// 设置操作日志集合
			model.SetOperationLogCollection(mongoComponent.GetCollection("operation_log"))
		}
	}

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
		// 服务连接 监听
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
