package router

import (
	"github.com/gin-gonic/gin"
)

// RouteRegistrarInterface 路由注册接口
type RouteRegistrarInterface interface {
	// RegisterRoutes 注册路由
	RegisterRoutes(r *gin.RouterGroup, params RouterParams)
}

// SeparateRouteRegistrar 分离模式路由注册器
type SeparateRouteRegistrar struct{}

// RegisterRoutes 注册分离模式路由
func (srr *SeparateRouteRegistrar) RegisterRoutes(r *gin.RouterGroup, params RouterParams) {
	registerSeparateRoutes(r, params)
}
