package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// registerPublicRoutes 注册公开路由（无需认证）
func registerPublicRoutes(r *gin.Engine, params RouterParams) {
	// 公开文件访问
	publicFiles := r.Group("/public")
	{
		publicFiles.GET("/files/:id", params.FileHandler.GetFileInfo)
	}
}

// registerAppRoutes 注册应用路由
func registerAppRoutes(r *gin.RouterGroup, params RouterParams) {
	ctx := context.Background()
	logger.InfoContext(ctx, "Registering application routes")

	// 公共路由
	public := r.Group("")
	{
		// 用户登录（管理员和普通用户使用同一接口）
		public.POST("/login", params.UserHandler.UserLogin)

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
			files.POST("/upload-url", params.FileHandler.GetUploadURL)
			files.POST("/confirm", params.FileHandler.ConfirmUpload)
			files.GET("/:id/download", params.FileHandler.GetDownloadURL)
			files.DELETE("/:id", params.FileHandler.DeleteFile)
		}

		// 系统设置
		admin.GET("/settings", params.AdminHandler.GetSystemSettings)
	}

	// 文件上传接口（统一支持本地和MinIO存储，需要认证但不需要管理员权限）
	upload := authenticated.Group("/upload")
	{
		upload.POST("/file", params.FileHandler.UploadFile) // 统一上传接口
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
