package core

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/storage/database"
)

// App 应用实例
type App struct {
	ComponentManager *ComponentManager
	AppConfig        *configs.Config
	ServiceFactory   any // 服务工厂
}

// Setup 创建应用实例
func Setup(cfg *configs.Config) *App {
	return &App{
		ComponentManager: NewComponentManager(cfg),
		AppConfig:        cfg,
	}
}

// SetServiceFactory 设置服务工厂
func (a *App) SetServiceFactory(factory any) {
	a.ServiceFactory = factory
}

// GetServiceFactory 获取服务工厂
func (a *App) GetServiceFactory() any {
	return a.ServiceFactory
}

// GetDatabase 获取数据库实例
func (a *App) GetDatabase() database.Database {
	// 从组件管理器中获取数据库组件
	for _, component := range a.ComponentManager.components {
		if db, ok := component.(database.Database); ok {
			return db
		}
	}
	return nil
}

// Initialize 初始化应用
func (a *App) Initialize() error {
	return a.ComponentManager.Initialize()
}

// Cleanup 清理资源
func (a *App) Cleanup() {
	a.ComponentManager.Cleanup()
}

// Config 获取配置
func (a *App) Config() *configs.Config {
	return a.AppConfig
}
