package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/epkgs/i18n"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/handler"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// newRouter 创建路由器（不依赖fx）
func newRouter(config *configs.Config, handlers ...handler.RouterInitializer) (*gin.Engine, error) {
	// 设置Gin模式
	gin.SetMode(config.App.Mode)

	// 创建路由器
	r := gin.New()

	// 添加中间件
	r.Use(middleware.RequestLoggerMiddleware())
	r.Use(middleware.Cors())

	// 添加国际化中间件
	if config.I18n.Enabled {
		r.Use(i18n.GinMiddleware(config.I18n.DefaultLanguage))
	}

	// 添加错误处理中间件（替换gin.Recovery()）
	r.Use(middleware.PanicRecovery())
	r.Use(middleware.GlobalErrorHandler())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, &dto.HealthResponse{
			Status: "ok",
		})
	})

	// 如果启用了pprof且使用主服务器端口，则添加pprof路由
	if config.Pprof.Enabled && config.Pprof.Port == 0 {
		registerPprofRoutes(r)
	}

	ctx := context.Background()
	logger.InfoContext(ctx, "Registering application routes")

	// 创建API路由组
	api := r.Group("/api/v1")

	// 注册应用路由
	for _, h := range handlers {
		h.InitRouters(api, r)
	}

	// 打印路由信息
	logger.Info("==================================================")
	logger.Info("Route information:")

	// 获取所有路由
	routes := r.Routes()
	for _, route := range routes {
		logger.Info(fmt.Sprintf("%-7s %s", route.Method, route.Path))
	}

	logger.Info("==================================================")

	return r, nil
}

// registerPprofRoutes 注册pprof路由
func registerPprofRoutes(r *gin.Engine) {
	// 创建pprof路由组
	pprofGroup := r.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
			pprof.Index(w, r)
		}))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
		pprofGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
		pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
		pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
		pprofGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
		pprofGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	}
}
