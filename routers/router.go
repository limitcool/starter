package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/storage/casbin"
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

	// 初始化Casbin组件
	casbinComponent := casbin.NewComponent(global.Config)
	if err := casbinComponent.Initialize(); err != nil {
		panic("Casbin组件初始化失败: " + err.Error())
	}

	// v1 router
	apiV1 := r.Group("/api/v1")

	// 公共路由
	{
		apiV1.GET("/ping", controller.Ping)
	}

	// 需要认证的路由
	auth := apiV1.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		// 获取当前用户菜单
		auth.GET("/user/menus", controller.GetUserMenus)
		// 获取当前用户权限标识
		auth.GET("/user/perms", controller.GetUserMenuPerms)
	}

	// 需要权限控制的路由
	admin := auth.Group("/admin")
	admin.Use(middleware.CasbinComponentMiddleware())
	{
		// 菜单管理
		menu := admin.Group("/menu")
		{
			menu.POST("", controller.CreateMenu)
			menu.PUT("/:id", controller.UpdateMenu)
			menu.DELETE("/:id", controller.DeleteMenu)
			menu.GET("/:id", controller.GetMenu)
			menu.GET("", controller.GetMenuTree)
		}

		// 角色管理
		role := admin.Group("/role")
		{
			role.POST("", controller.CreateRole)
			role.PUT("/:id", controller.UpdateRole)
			role.DELETE("/:id", controller.DeleteRole)
			role.GET("/:id", controller.GetRole)
			role.GET("", controller.GetRoles)
			// 为角色分配菜单
			role.POST("/menu", controller.AssignMenuToRole)
			// 为角色设置权限
			role.POST("/permission", controller.SetRolePermission)
			// 删除角色权限
			role.DELETE("/permission", controller.DeleteRolePermission)
		}

		// 权限管理
		permission := admin.Group("/permission")
		{
			permission.GET("", controller.GetPermissions)
			permission.GET("/:id", controller.GetPermission)
			permission.POST("", controller.CreatePermission)
			permission.PUT("/:id", controller.UpdatePermission)
			permission.DELETE("/:id", controller.DeletePermission)
		}
	}

	return r
}
