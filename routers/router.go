package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/middleware"
)

// Load loads the middlewares, routes, handlers.
func NewRouter() *gin.Engine {
	// 创建不带默认中间件的路由
	r := gin.New()

	// 使用我们的结构化日志中间件
	r.Use(middleware.LoggerWithCharmbracelet())

	// 使用恢复中间件
	r.Use(gin.Recovery())
	// 使用CORS中间件
	r.Use(middleware.Cors())
	// v1 router
	apiV1 := r.Group("/api/v1")
	// apiV1.Use()
	// {
	// 	apiV1.GET("/ping", handler.Ping)
	// }
	auth := apiV1.Use(middleware.AuthMiddleware())
	{
		auth.GET("/ping", controller.Ping)
	}
	return r
}
