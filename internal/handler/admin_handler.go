package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	*BaseHandler
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(db *gorm.DB, config *configs.Config) *AdminHandler {
	handler := &AdminHandler{
		BaseHandler: NewBaseHandler(db, config),
	}

	handler.LogInit("AdminHandler")
	return handler
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
