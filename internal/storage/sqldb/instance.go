package sqldb

import (
	"errors"
	"sync"

	"gorm.io/gorm"
)

// 错误定义
var (
	ErrDBNotInitialized = errors.New("database not initialized")
)

var (
	instance *Component
	once     sync.Once
)

// db.go中已经定义了DB和GetDB，这里不需要重复定义
// var DB *gorm.DB

// setupInstance 设置全局实例
func setupInstance(component *Component) {
	once.Do(func() {
		instance = component
	})
}

// Instance 获取数据库组件实例
func Instance() *Component {
	return instance
}

// GetDBWithError 获取数据库连接（带错误返回）
func GetDBWithError() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	if instance != nil && instance.DB() != nil {
		return instance.DB(), nil
	}

	return nil, ErrDBNotInitialized
}
