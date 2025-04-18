package database

import (
	"gorm.io/gorm"
)

// DB 数据库接口
// 简化的接口，只包含必要的方法
type DB interface {
	// GetDB 获取数据库连接
	GetDB() *gorm.DB

	// Close 关闭数据库连接
	Close() error
}

// GormDB 基于GORM的数据库实现
type GormDB struct {
	DB *gorm.DB
}

// NewGormDB 创建GORM数据库实例
func NewGormDB(db *gorm.DB) *GormDB {
	return &GormDB{DB: db}
}

// GetDB 获取数据库连接
func (g *GormDB) GetDB() *gorm.DB {
	return g.DB
}

// Close 关闭数据库连接
func (g *GormDB) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
