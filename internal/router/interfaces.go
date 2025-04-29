package router

import (
	"github.com/gin-gonic/gin"
)

// RouteRegistrarInterface 路由注册接口
type RouteRegistrarInterface interface {
	// RegisterRoutes 注册路由
	RegisterRoutes(r *gin.RouterGroup, params RouterParams)
}

// RouteRegistrar 路由注册器
type RouteRegistrar struct{}

// RegisterRoutes 注册应用路由
func (rr *RouteRegistrar) RegisterRoutes(r *gin.RouterGroup, params RouterParams) {
	registerAppRoutes(r, params)
}
