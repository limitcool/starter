package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// registerSimpleRoutes 注册简单模式的路由
func registerSimpleRoutes(r *gin.RouterGroup, params RouterParams) {
	ctx := context.Background()
	logger.InfoContext(ctx, "注册简单模式路由")

	// 公共路由
	public := r.Group("")
	{
		// 用户登录
		public.POST("/login", params.UserController.UserLogin)
		public.POST("/admin/login", params.AdminUserController.AdminUserLogin)

		// 刷新令牌
		public.POST("/refresh-token", params.UserController.RefreshToken)
	}

	// 需要认证的路由
	authenticated := r.Group("", middleware.JWTAuth(params.Config))

	// 管理员路由 - 使用简化的管理员检查中间件
	// 检查UserRepo是否为空
	var admin *gin.RouterGroup
	if params.UserRepo != nil {
		admin = authenticated.Group("/admin", middleware.SimpleAdminCheck(params.UserRepo))
	} else {
		logger.WarnContext(ctx, "用户仓库为空，无法添加管理员检查中间件，使用默认路由")
		admin = authenticated.Group("/admin")
	}
	{
		// 用户管理
		// 注意：这里的用户管理功能需要实现
		// users := admin.Group("/users")

		// 文件管理
		files := admin.Group("/files")
		{
			// 上传文件
			files.POST("/upload", params.FileController.UploadFile)

			// 注意：以下方法需要实现
			// files.GET("", params.FileController.GetFiles)
			// files.GET("/:id", params.FileController.GetFile)
			// files.DELETE("/:id", params.FileController.DeleteFile)
		}

		// 系统设置
		admin.GET("/settings", params.AdminController.GetSystemSettings)
		// 注意：以下方法需要实现
		// admin.PUT("/settings", params.AdminController.UpdateSystemSettings)
	}

	// 普通用户路由 - 使用JWT认证
	user := authenticated.Group("/user")
	{
		// 用户信息
		user.GET("/info", params.UserController.UserInfo)

		// 修改密码
		user.POST("/change-password", params.UserController.UserChangePassword)

		// 更新个人信息
		// 注意：以下方法需要实现
		// user.PUT("/profile", params.UserController.UpdateUserProfile)
	}
}
