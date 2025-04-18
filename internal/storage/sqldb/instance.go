package sqldb

import (
	"errors"
	"sync"

	"github.com/charmbracelet/log"
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

// setupInstance 设置全局实例
// 已弃用: 请使用依赖注入而不是全局实例
func setupInstance(component *Component) {
	once.Do(func() {
		instance = component
		log.Warn("使用了已弃用的全局实例设置方法，请使用依赖注入")
	})
}

// Instance 获取数据库组件实例
// 已弃用: 请使用依赖注入而不是全局实例
func Instance() *Component {
	if instance == nil {
		log.Warn("使用了未初始化的全局数据库实例，请使用依赖注入")
	}
	return instance
}

// GetDBWithError 获取数据库连接（带错误返回）
// 已弃用: 请使用依赖注入而不是全局实例
func GetDBWithError() (*gorm.DB, error) {
	if instance != nil && instance.DB() != nil {
		log.Warn("使用了已弃用的全局数据库实例获取方法，请使用依赖注入")
		return instance.DB(), nil
	}

	return nil, ErrDBNotInitialized
}
