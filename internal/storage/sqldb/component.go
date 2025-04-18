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

// Database 数据库接口
type Database interface {
	// DB 获取数据库连接
	DB() *gorm.DB

	// Close 关闭数据库连接
	Close() error

	// IsEnabled 检查数据库是否启用
	IsEnabled() bool
}

// DB 数据库接口
type DB interface {
	// GetDB 获取数据库连接
	GetDB() *gorm.DB

	// Close 关闭数据库连接
	Close() error
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

	// 执行设置全局实例
	setupInstance(c)

	// 不再自动执行迁移
	// if err := db.AutoMigrate(); err != nil {
	// 	return fmt.Errorf("database migration failed: %w", err)
	// }

	log.Info("SQL database component initialized successfully")
	return nil
}

// Migrate 执行数据库迁移
func (c *Component) Migrate() error {
	if c.db == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Info("Running database migrations")
	if err := c.db.AutoMigrate(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	log.Info("Database migrations completed successfully")
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
			log.Info("Closing database connection")
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
