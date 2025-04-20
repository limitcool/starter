package router

import (
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// printRoutes 打印所有路由信息
func printRoutes(engine *gin.Engine) {
	routes := getRoutes(engine)
	if len(routes) == 0 {
		return
	}

	// 按路径排序
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Path < routes[j].Path
	})

	// 打印路由信息
	logger.Info("==================================================")
	logger.Info("路由信息:")
	for _, route := range routes {
		logger.Info(fmt.Sprintf("%-7s %s", route.Method, route.Path))
	}
	logger.Info("==================================================")
}

// routeInfo 路由信息
type routeInfo struct {
	Method string
	Path   string
}

// getRoutes 获取所有路由信息
func getRoutes(engine *gin.Engine) []routeInfo {
	routes := make([]routeInfo, 0)

	// 获取所有路由
	ginRoutes := engine.Routes()

	// 遍历路由
	for _, route := range ginRoutes {
		routes = append(routes, routeInfo{
			Method: route.Method,
			Path:   route.Path,
		})
	}

	return routes
}
