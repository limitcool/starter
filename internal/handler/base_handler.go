package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/filestore"
	"github.com/limitcool/starter/internal/pkg/cache"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

type RouterInitializer interface {
	InitRouters(g *gin.RouterGroup, root *gin.Engine)
}

type AppContext interface {
	GetConfig() *configs.Config
	GetDB() *gorm.DB
	GetCache() cache.Cache
	GetStorage() filestore.FileStorage
}

// BaseHandler 基础处理器，包含所有Handler的公共字段和方法
type BaseHandler struct {
	DB     *gorm.DB
	Config *configs.Config
	Helper *HandlerHelper
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler(db *gorm.DB, config *configs.Config) *BaseHandler {
	return &BaseHandler{
		DB:     db,
		Config: config,
		Helper: NewHandlerHelper(),
	}
}

// LogInit 记录Handler初始化日志
func (h *BaseHandler) LogInit(handlerName string) {
	logger.Info(handlerName + " initialized")
}
