package database

import (
	"gorm.io/gorm"
)

// 注意：我们使用 database.go 中定义的 Database 接口
// 这里不再定义重复的 DB 接口

// GormDB 基于GORM的数据库实现
type GormDB struct {
	db      *gorm.DB
	enabled bool
}

// NewGormDB 创建GORM数据库实例
func NewGormDB(db *gorm.DB) *GormDB {
	return &GormDB{db: db, enabled: true}
}

// DB 获取数据库连接
// 实现 Database 接口
func (g *GormDB) DB() *gorm.DB {
	return g.db
}

// Close 关闭数据库连接
func (g *GormDB) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// IsEnabled 检查数据库是否启用
func (g *GormDB) IsEnabled() bool {
	return g.enabled
}
