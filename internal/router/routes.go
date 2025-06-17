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
			// 获取上传预签名URL（推荐方式）
			files.POST("/upload-url", params.FileHandler.GetUploadURL)
			// 确认上传完成
			files.POST("/confirm-upload", params.FileHandler.ConfirmUpload)
			// 上传文件（废弃，保留兼容性）
			files.POST("/upload", params.FileHandler.UploadFile)
			// 获取文件列表
			files.GET("/list", params.FileHandler.ListFiles)
			// 获取文件信息
			files.GET("/:id", params.FileHandler.GetFileInfo)
			// 获取文件访问URL
			files.GET("/:id/url", params.FileHandler.GetFileURL)
		}

		// 系统设置
		admin.GET("/settings", params.AdminHandler.GetSystemSettings)
	}

	// 公开文件访问（无需认证）
	publicFiles := r.Group("/public")
	{
		// 公开文件访问（长期URL的实现）
		publicFiles.GET("/files/:id", params.FileHandler.ServePublicFile)
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
