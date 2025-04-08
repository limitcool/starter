package redisdb

import (
	"errors"
	"sync"

	"github.com/redis/go-redis/v9"
)

// 错误定义
var (
	// ErrNotInitialized Redis组件未初始化错误
	ErrNotInitialized = errors.New("redis component not initialized")
)

var (
	instance *Component
	once     sync.Once
)

// setupInstance 设置全局Redis组件实例（仅内部使用）
func setupInstance(component *Component) {
	once.Do(func() {
		instance = component
	})
}

// Instance 获取全局Redis组件实例
func Instance() *Component {
	return instance
}

// IsEnabled 检查Redis是否启用
func IsEnabled() bool {
	return instance != nil && instance.IsEnabled()
}

// GetClient 获取默认Redis客户端
func GetClient() (*redis.Client, error) {
	if instance == nil {
		return nil, ErrNotInitialized
	}
	return instance.GetClient()
}

// GetClientByName 获取指定名称的Redis客户端
func GetClientByName(name string) (*redis.Client, error) {
	if instance == nil {
		return nil, ErrNotInitialized
	}
	return instance.GetClientByName(name)
}
