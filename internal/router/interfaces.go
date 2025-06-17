package router

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/handler"
)

// RouteRegistrarInterface 路由注册接口
type RouteRegistrarInterface interface {
	// RegisterRoutes 注册路由
	RegisterRoutes(r *gin.RouterGroup, params RouterParams)
	// RegisterRoutesSimple 简单注册路由（不依赖fx）
	RegisterRoutesSimple(r *gin.Engine, config *configs.Config, userHandler *handler.UserHandler, fileHandler *handler.FileHandler, adminHandler *handler.AdminHandler)
}

// RouteRegistrar 路由注册器
type RouteRegistrar struct{}

// RegisterRoutes 注册应用路由
func (rr *RouteRegistrar) RegisterRoutes(r *gin.RouterGroup, params RouterParams) {
	// 注册应用路由
	registerAppRoutes(r, params)
}

// RegisterRoutesSimple 简单注册路由（不依赖fx）
func (rr *RouteRegistrar) RegisterRoutesSimple(r *gin.Engine, config *configs.Config, userHandler *handler.UserHandler, fileHandler *handler.FileHandler, adminHandler *handler.AdminHandler) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 创建API路由组
	api := r.Group("/api/v1")

	// 创建参数结构
	params := RouterParams{
		Config:       config,
		UserHandler:  userHandler,
		FileHandler:  fileHandler,
		AdminHandler: adminHandler,
	}

	// 注册应用路由
	registerAppRoutes(api, params)
}
