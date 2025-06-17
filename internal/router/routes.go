package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// registerRoutes 注册应用路由
func registerAppRoutes(r *gin.RouterGroup, params RouterParams) {
	ctx := context.Background()
	logger.InfoContext(ctx, "Registering application routes")

	// 公共路由
	public := r.Group("")
	{
		// 用户登录
		public.POST("/login", params.UserHandler.UserLogin)
		// 在lite版本中，管理员登录也使用普通用户登录接口
		public.POST("/admin/login", params.UserHandler.UserLogin)

		// 用户注册
		public.POST("/register", params.UserHandler.UserRegister)
	}

	// 需要认证的路由
	authenticated := r.Group("", middleware.JWTAuth(params.Config))

	// 管理员路由 - 使用简化的管理员检查中间件
	admin := authenticated.Group("/admin", middleware.AdminCheck())
	{
		// 文件管理
		files := admin.Group("/files")
		{
			// 上传文件
			files.POST("/upload", params.FileHandler.UploadFile)
		}

		// 系统设置
		admin.GET("/settings", params.AdminHandler.GetSystemSettings)
	}

	// 普通用户路由 - 使用JWT认证
	user := authenticated.Group("/user")
	{
		// 用户信息
		user.GET("/info", params.UserHandler.UserInfo)

		// 修改密码
		user.POST("/change-password", params.UserHandler.UserChangePassword)
	}
}
