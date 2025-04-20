package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
)

// Module 应用程序模块
// 包含所有依赖项和生命周期钩子
var Module = fx.Options(
	// 提供配置
	fx.Provide(configs.LoadConfig),

	// 提供日志
	fx.Provide(NewLogger),

	// 提供HTTP服务器
	fx.Provide(NewHTTPServer),

	// 注册生命周期钩子
	fx.Invoke(RegisterHooks),
)

// NewLogger 创建日志记录器
func NewLogger(cfg *configs.Config) (logger.Logger, error) {
	// 初始化日志
	logger.Setup(cfg.Log)
	return logger.Default(), nil
}

// HTTPServer HTTP服务器
type HTTPServer struct {
	Server *http.Server
	Router *gin.Engine
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(lc fx.Lifecycle, cfg *configs.Config, router *gin.Engine) *HTTPServer {
	// 创建HTTP服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.App.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second, // 默认值
		WriteTimeout:   10 * time.Second, // 默认值
		MaxHeaderBytes: 1 << 20,          // 1MB
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server", "port", cfg.App.Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server failed", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			// 创建一个带超时的上下文
			stopCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			return srv.Shutdown(stopCtx)
		},
	})

	return &HTTPServer{
		Server: srv,
		Router: router,
	}
}

// RegisterHooks 注册其他生命周期钩子
func RegisterHooks(lc fx.Lifecycle, log *logger.Logger) {
	// 注册日志生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Application starting")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Application stopping")
			return nil
		},
	})
}
