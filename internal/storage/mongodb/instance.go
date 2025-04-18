package mongodb

import (
	"errors"
	"sync"

	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// 错误定义
var (
	// ErrNotInitialized MongoDB组件未初始化错误
	ErrNotInitialized = errors.New("mongodb component not initialized")
)

// 已弃用: 请使用 Instance().GetClient() 获取MongoDB客户端
var Mongo *mongo.Client

var (
	instance *Component
	once     sync.Once
)

// setupInstance 设置全局MongoDB组件实例（仅内部使用）
func setupInstance(component *Component) {
	once.Do(func() {
		instance = component
	})
}

// Instance 获取全局MongoDB组件实例
func Instance() *Component {
	return instance
}

// GetDatabaseName 获取MongoDB数据库名称
func GetDatabaseName() string {
	if instance != nil && instance.Config != nil {
		return instance.Config.Mongo.DB
	}
	return ""
}

// IsEnabled 检查MongoDB是否启用
func IsEnabled() bool {
	return instance != nil && instance.IsEnabled()
}

// GetClient 获取MongoDB客户端
func GetClient() (*mongo.Client, error) {
	if instance == nil || !instance.IsEnabled() {
		return nil, ErrNotInitialized
	}
	return instance.GetClient(), nil
}

// GetDB 获取默认数据库
func GetDB() (*mongo.Database, error) {
	if instance == nil || !instance.IsEnabled() {
		return nil, ErrNotInitialized
	}
	return instance.GetDB(), nil
}

// GetCollection 获取集合
func GetCollection(name string) (*mongo.Collection, error) {
	if instance == nil || !instance.IsEnabled() {
		return nil, ErrNotInitialized
	}
	return instance.GetCollection(name), nil
}

// 为兼容性提供的方法
// 已弃用: 请使用 Instance().GetCollection(name) 获取集合
func Collection(name string) *mongo.Collection {
	if instance != nil && instance.IsEnabled() {
		return instance.GetCollection(name)
	}
	log.Warn("使用了已弃用的MongoDB访问方式，请使用 Instance().GetCollection()", "caller", "mongodb.Collection()")
	return nil
}
