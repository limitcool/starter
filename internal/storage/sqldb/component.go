package sqldb

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"gorm.io/gorm"
)

// Component SQL数据库组件
type Component struct {
	Config  *configs.Config // 导出字段
	db      *gorm.DB
	enabled bool
}

// NewComponent 创建数据库组件
func NewComponent(cfg *configs.Config) *Component {
	return &Component{
		Config:  cfg,
		enabled: cfg.Database.Enabled,
	}
}

// Name 返回组件名称
func (c *Component) Name() string {
	return "SQLDatabase"
}

// Initialize 初始化数据库连接
func (c *Component) Initialize() error {
	if !c.enabled {
		log.Info("SQL database component disabled")
		return nil
	}

	log.Info("Initializing SQL database component")

	db := newDbConn(c.Config)
	if db == nil {
		return fmt.Errorf("failed to create database connection")
	}

	c.db = db
	DB = db // 设置全局访问点

	// 执行设置全局实例
	setupInstance(c)

	// 执行数据库迁移
	if err := db.AutoMigrate(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	log.Info("SQL database component initialized successfully")
	return nil
}

// Cleanup 清理数据库资源
func (c *Component) Cleanup() {
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err == nil {
			log.Info("Closing database connection")
			_ = sqlDB.Close()
		}
	}
}

// IsEnabled 检查组件是否启用
func (c *Component) IsEnabled() bool {
	return c.enabled
}

// GetDB 获取数据库连接
func (c *Component) GetDB() *gorm.DB {
	return c.db
}
