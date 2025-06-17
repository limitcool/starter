package handler

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// BaseHandler 基础处理器，包含所有Handler的公共字段和方法
type BaseHandler struct {
	DB       *gorm.DB
	Config   *configs.Config
	Helper   *HandlerHelper
	FileUtil *FileUtil
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler(db *gorm.DB, config *configs.Config) *BaseHandler {
	return &BaseHandler{
		DB:       db,
		Config:   config,
		Helper:   NewHandlerHelper(),
		FileUtil: NewFileUtil("/uploads"), // 默认基础URL
	}
}

// LogInit 记录Handler初始化日志
func (h *BaseHandler) LogInit(handlerName string) {
	logger.Info(handlerName + " initialized")
}
