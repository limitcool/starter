package router

import (
	"github.com/gin-gonic/gin"
)

// RouteRegistrarInterface 路由注册接口
type RouteRegistrarInterface interface {
	// RegisterRoutes 注册路由
	RegisterRoutes(r *gin.RouterGroup, params RouterParams)
}
