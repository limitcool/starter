package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// registerSeparateRoutes 注册分离模式的路由
func registerSeparateRoutes(r *gin.RouterGroup, params RouterParams) {
	ctx := context.Background()
	logger.InfoContext(ctx, "注册分离模式路由")

	// 用户路由
	userRoutes(r, params)

	// 管理员路由
	adminRoutes(r, params)
}

// userRoutes 用户路由
func userRoutes(r *gin.RouterGroup, params RouterParams) {
	// 用户路由
	users := r.Group("/users")
	{
		// 用户登录
		users.POST("/login", params.UserController.UserLogin)

		// 用户注册
		users.POST("/register", params.UserController.UserRegister)

		// 刷新令牌
		users.POST("/refresh-token", params.UserController.RefreshToken)

		// 需要认证的路由
		authenticated := users.Group("", middleware.JWTAuth(params.Config))
		{
			// 获取用户信息
			authenticated.GET("/info", params.UserController.UserInfo)

			// 修改密码
			authenticated.POST("/change-password", params.UserController.UserChangePassword)

			// 更新个人信息
			// 注意：以下方法需要实现
			// authenticated.PUT("/profile", params.UserController.UpdateUserProfile)
		}
	}
}

// adminRoutes 管理员路由
func adminRoutes(r *gin.RouterGroup, params RouterParams) {
	// 管理员路由
	admin := r.Group("/admin-api")
	{
		// 管理员登录
		admin.POST("/login", params.AdminUserController.AdminUserLogin)

		// 刷新令牌
		// 注意：以下方法需要实现
		// admin.POST("/refresh-token", params.AdminUserController.RefreshToken)

		// 需要认证的路由
		authenticated := admin.Group("", middleware.JWTAuth(params.Config))
		{
			// 获取管理员信息
			// 注意：以下方法需要实现
			// authenticated.GET("/info", params.AdminUserController.GetUserInfo)

			// 检查是否有角色控制器
			if params.RoleController != nil {
				// 角色管理
				roles := authenticated.Group("/roles")
				{
					// 获取角色列表
					roles.GET("", params.RoleController.GetRoles)

					// 获取角色详情
					roles.GET("/:id", params.RoleController.GetRole)

					// 创建角色
					roles.POST("", params.RoleController.CreateRole)

					// 更新角色
					roles.PUT("/:id", params.RoleController.UpdateRole)

					// 删除角色
					roles.DELETE("/:id", params.RoleController.DeleteRole)

					// 为角色分配菜单
					roles.POST("/menu", params.RoleController.AssignMenuToRole)

					// 为角色设置权限
					roles.POST("/permission", params.RoleController.SetRolePermission)

					// 删除角色权限
					roles.DELETE("/permission", params.RoleController.DeleteRolePermission)
				}
			}

			// 检查是否有菜单控制器
			if params.MenuController != nil {
				// 菜单管理
				menus := authenticated.Group("/menus")
				{
					// 获取菜单详情
					menus.GET("/:id", params.MenuController.GetMenu)

					// 获取菜单树
					menus.GET("/tree", params.MenuController.GetMenuTree)

					// 获取用户菜单
					menus.GET("/user", params.MenuController.GetUserMenus)

					// 获取用户菜单权限
					menus.GET("/user/perms", params.MenuController.GetUserMenuPerms)

					// 创建菜单
					menus.POST("", params.MenuController.CreateMenu)

					// 更新菜单
					menus.PUT("/:id", params.MenuController.UpdateMenu)

					// 删除菜单
					menus.DELETE("/:id", params.MenuController.DeleteMenu)
				}
			}

			// 检查是否有权限控制器
			if params.PermissionController != nil {
				// 权限管理
				permissions := authenticated.Group("/permissions")
				{
					// 获取权限列表
					permissions.GET("", params.PermissionController.GetPermissions)

					// 获取权限详情
					permissions.GET("/:id", params.PermissionController.GetPermission)

					// 获取用户权限
					permissions.GET("/user", params.PermissionController.GetUserPermissions)

					// 获取用户菜单
					permissions.GET("/user/menus", params.PermissionController.GetUserMenus)

					// 获取用户角色
					permissions.GET("/user/roles", params.PermissionController.GetUserRoles)

					// 为用户分配角色
					permissions.POST("/user/:id/roles", params.PermissionController.AssignRolesToUser)

					// 创建权限
					permissions.POST("", params.PermissionController.CreatePermission)

					// 更新权限
					permissions.PUT("/:id", params.PermissionController.UpdatePermission)

					// 删除权限
					permissions.DELETE("/:id", params.PermissionController.DeletePermission)

					// 为角色分配权限
					permissions.POST("/role/:id", params.PermissionController.AssignPermissionsToRole)

					// 获取角色权限
					permissions.GET("/role/:id", params.PermissionController.GetRolePermissions)

					// 更新权限系统设置
					permissions.PUT("/settings", params.PermissionController.UpdatePermissionSettings)
				}
			}

			// 检查是否有操作日志控制器
			if params.OperationLogController != nil {
				// 操作日志管理
				logs := authenticated.Group("/logs")
				{
					// 获取操作日志列表
					logs.GET("", params.OperationLogController.GetOperationLogs)

					// 删除操作日志
					logs.DELETE("/:id", params.OperationLogController.DeleteOperationLog)

					// 批量删除操作日志
					logs.POST("/batch-delete", params.OperationLogController.ClearOperationLogs)
				}
			}

			// 文件管理
			files := authenticated.Group("/files")
			{
				// 上传文件
				files.POST("/upload", params.FileController.UploadFile)

				// 注意：以下方法需要实现
				// files.GET("", params.FileController.GetFiles)
				// files.GET("/:id", params.FileController.GetFile)
				// files.DELETE("/:id", params.FileController.DeleteFile)
			}

			// 系统设置
			authenticated.GET("/settings", params.AdminController.GetSystemSettings)
			// 注意：以下方法需要实现
			// authenticated.PUT("/settings", params.AdminController.UpdateSystemSettings)
		}
	}
}
