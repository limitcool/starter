package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/handler"
)

// RouteRegistrar 路由注册器
type RouteRegistrar struct{}

// RegisterRoutes 注册路由
func (rr *RouteRegistrar) RegisterRoutes(r *gin.Engine, config *configs.Config, userHandler *handler.UserHandler, fileHandler *handler.FileHandler, adminHandler *handler.AdminHandler) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, &dto.HealthResponse{
			Status: "ok",
		})
	})

	// 如果启用了pprof且使用主服务器端口，则添加pprof路由
	if config.Pprof.Enabled && config.Pprof.Port == 0 {
		rr.registerPprofRoutes(r)
	}

	// 创建参数结构
	params := RouterParams{
		Config:       config,
		UserHandler:  userHandler,
		FileHandler:  fileHandler,
		AdminHandler: adminHandler,
	}

	// 注册公开路由（无需认证）
	registerPublicRoutes(r, params)

	// 创建API路由组
	api := r.Group("/api/v1")

	// 注册应用路由
	registerAppRoutes(api, params)
}

// registerPprofRoutes 注册pprof路由
func (rr *RouteRegistrar) registerPprofRoutes(r *gin.Engine) {
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
