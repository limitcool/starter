package casbin

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// Component Casbin组件实现
type Component struct {
	config  *configs.Config
	db      *gorm.DB
	service Service
	enabled bool
}

// NewComponent 创建Casbin组件
func NewComponent(cfg *configs.Config, db *gorm.DB) *Component {
	return &Component{
		config:  cfg,
		db:      db,
		enabled: cfg.Casbin.Enabled,
	}
}

// Name 返回组件名称
func (c *Component) Name() string {
	return "Casbin"
}

// Initialize 初始化Casbin组件
func (c *Component) Initialize() error {
	if !c.enabled {
		logger.Info("Casbin component disabled")
		return nil
	}

	logger.Info("Initializing Casbin component")
	
	// 创建Casbin服务
	c.service = NewService(c.db, c.config)
	
	// 初始化服务
	if err := c.service.Initialize(); err != nil {
		return err
	}
	
	logger.Info("Casbin component initialized successfully")
	return nil
}

// Cleanup 清理Casbin组件资源
func (c *Component) Cleanup() {
	// Casbin没有需要清理的资源
	logger.Debug("Cleaning up Casbin component")
}

// GetService 获取Casbin服务
func (c *Component) GetService() Service {
	return c.service
}

// IsEnabled 检查组件是否启用
func (c *Component) IsEnabled() bool {
	return c.enabled
}
