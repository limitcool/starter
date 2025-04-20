package database

import (
	"gorm.io/gorm"
)

// Database 数据库接口
type Database interface {
	// DB 获取数据库连接
	DB() *gorm.DB
	
	// Close 关闭数据库连接
	Close() error
	
	// IsEnabled 检查数据库是否启用
	IsEnabled() bool
}

// Repository 数据库仓库接口
// 所有需要数据库操作的服务都应该实现这个接口
type Repository interface {
	// SetDB 设置数据库连接
	SetDB(db Database)
	
	// GetDB 获取数据库连接
	GetDB() *gorm.DB
}

// BaseRepository 基础仓库实现
// 可以被其他服务嵌入使用
type BaseRepository struct {
	db Database
}

// SetDB 设置数据库连接
func (r *BaseRepository) SetDB(db Database) {
	r.db = db
}

// GetDB 获取数据库连接
func (r *BaseRepository) GetDB() *gorm.DB {
	if r.db != nil {
		return r.db.DB()
	}
	return nil
}
