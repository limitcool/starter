package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/handlers"
	"github.com/limitcool/starter/internal/middleware"
)

// Load loads the middlewares, routes, handlers.
func NewRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	g := gin.New()
	g.Use(gin.Recovery(), middleware.LoggerWithCharmbracelet(), middleware.Cors())
	// v1 router
	apiV1 := g.Group("/api/v1")
	// apiV1.Use()
	// {
	// 	apiV1.GET("/ping", handler.Ping)
	// }
	auth := apiV1.Use(middleware.AuthMiddleware())
	{
		auth.GET("/ping", handlers.Ping)
	}
	return g
}
