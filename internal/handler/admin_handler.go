package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	*BaseHandler
	app IAPP
}

var _ IHandler = (*AdminHandler)(nil) // 用于接口断言，_ 变量编译后会被移除

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(app IAPP) *AdminHandler {
	handler := &AdminHandler{
		BaseHandler: NewBaseHandler(app.GetDB(), app.GetConfig()),
		app:         app,
	}

	handler.LogInit("AdminHandler")
	return handler
}

func (h *AdminHandler) InitRouters(g *gin.RouterGroup, root *gin.Engine) {

	// 需要认证的路由
	authenticated := g.Group("", middleware.JWTAuth(h.app.GetConfig()))

	// 管理员路由 - 使用简化的管理员检查中间件
	admin := authenticated.Group("/admin", middleware.AdminCheck())
	{

		// 系统设置
		admin.GET("/settings", h.GetSystemSettings)
	}
}

// GetSystemSettings 获取系统设置
func (h *AdminHandler) GetSystemSettings(ctx *gin.Context) {
	// 获取请求上下文
	reqCtx := ctx.Request.Context()

	// 记录请求
	logger.InfoContext(reqCtx, "GetSystemSettings 获取系统设置")

	// 返回系统设置
	response.Success(ctx, &dto.SystemSettingsResponse{
		AppName:    h.Config.App.Name,
		AppVersion: "1.0.0",
		AppMode:    h.Config.App.Mode,
	})
}
