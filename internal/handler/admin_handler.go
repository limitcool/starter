package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	db     *gorm.DB
	config *configs.Config
	logger *logger.Logger
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(db *gorm.DB, config *configs.Config, lc fx.Lifecycle) *AdminHandler {
	handler := &AdminHandler{
		db:     db,
		config: config,
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "AdminHandler initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "AdminHandler stopped")
			return nil
		},
	})

	return handler
}

// GetSystemSettings 获取系统设置
func (h *AdminHandler) GetSystemSettings(ctx *gin.Context) {
	// 获取请求上下文
	reqCtx := ctx.Request.Context()

	// 记录请求
	logger.InfoContext(reqCtx, "GetSystemSettings 获取系统设置")

	// 返回系统设置
	settings := map[string]any{
		"app_name":    h.config.App.Name,
		"app_version": "1.0.0",
		"app_mode":    h.config.App.Mode,
	}

	response.Success(ctx, settings)
}
