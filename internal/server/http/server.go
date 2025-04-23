package http

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

// ServerParams HTTP服务器参数
type ServerParams struct {
	fx.In

	LC     fx.Lifecycle
	Config *configs.Config
	Router *gin.Engine
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(params ServerParams) *http.Server {
	// 创建HTTP服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", params.Config.App.Port),
		Handler:        params.Router,
		ReadTimeout:    time.Duration(params.Config.App.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(params.Config.App.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "==================================================")
			logger.InfoContext(ctx, "HTTP服务器已启动",
				"address", fmt.Sprintf("http://localhost:%d", params.Config.App.Port),
				"mode", params.Config.App.Mode)
			logger.InfoContext(ctx, "==================================================")

			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.ErrorContext(context.Background(), "HTTP server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "Stopping HTTP server")

			// 创建一个5秒超时的上下文
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			// 优雅关闭服务器
			if err := srv.Shutdown(ctx); err != nil {
				logger.ErrorContext(ctx, "HTTP server shutdown error", "error", err)
				return err
			}

			logger.InfoContext(ctx, "HTTP server stopped")
			return nil
		},
	})

	return srv
}
