package sqldb

import (
	"fmt"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// Component SQL数据库组件
type Component struct {
	Config  *configs.Config // 导出字段
	db      *gorm.DB
	enabled bool
}

// Component 实现 database.Database 接口

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
		logger.Info("SQL database component disabled")
		return nil
	}

	logger.Info("Initializing SQL database component")

	db := newDbConn(c.Config)
	if db == nil {
		return fmt.Errorf("failed to create database connection")
	}

	c.db = db

	// 不再设置全局实例
	// setupInstance(c) // 已移除，使用依赖注入代替

	// 不再自动执行迁移
	// if err := db.AutoMigrate(); err != nil {
	// 	return fmt.Errorf("database migration failed: %w", err)
	// }

	logger.Info("SQL database component initialized successfully")
	return nil
}

// Migrate 执行数据库迁移
func (c *Component) Migrate() error {
	if c.db == nil {
		return fmt.Errorf("database not initialized")
	}

	logger.Info("Running database migrations")
	if err := c.db.AutoMigrate(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// Cleanup 清理数据库资源
func (c *Component) Cleanup() {
	_ = c.Close()
}

// Close 关闭数据库连接
func (c *Component) Close() error {
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err == nil {
			logger.Info("Closing database connection")
			return sqlDB.Close()
		}
		return err
	}
	return nil
}

// IsEnabled 检查组件是否启用
func (c *Component) IsEnabled() bool {
	return c.enabled
}

// DB 获取数据库连接
// 这是Component实例的方法，返回组件管理的db实例。
// 当您已经持有Component实例时，应优先使用此方法而非包级函数GetDB()，
// 这样可以获得更好的代码组织和依赖管理。
func (c *Component) DB() *gorm.DB {
	return c.db
}

// GetDB 获取数据库连接
// 实现database.DB接口
func (c *Component) GetDB() *gorm.DB {
	return c.db
}
