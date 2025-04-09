package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/controller"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/pkg/storage"
)

// SetupFileRoutes 文件路由设置
func SetupFileRoutes(router *gin.RouterGroup, storage *storage.Storage) {
	fileController := controller.NewFileController(storage)

	fileRoutes := router.Group("/files")
	{
		// 上传文件需要登录
		fileRoutes.POST("/upload", middleware.JWTAuth(), fileController.UploadFile)

		// 获取文件信息
		fileRoutes.GET("/:id", fileController.GetFile)

		// 下载文件
		fileRoutes.GET("/:id/download", fileController.DownloadFile)

		// 删除文件需要登录
		fileRoutes.DELETE("/:id", middleware.JWTAuth(), fileController.DeleteFile)
	}

	// 用户头像相关路由
	userRoutes := router.Group("/users")
	{
		// 更新头像需要用户登录
		userRoutes.PUT("/avatar", middleware.JWTAuth(), fileController.UpdateUserAvatar)
	}

	// 系统用户头像相关路由
	sysUserRoutes := router.Group("/sysusers")
	{
		// 更新头像需要系统用户登录
		sysUserRoutes.PUT("/avatar", middleware.JWTAuth(), fileController.UpdateSysUserAvatar)
	}
}
