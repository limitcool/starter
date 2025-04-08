package core

import (
	"sync"

	"github.com/limitcool/starter/configs"
)

// App 应用实例
type App struct {
	ComponentManager *ComponentManager
	Config           *configs.Config
}

var (
	instance *App
	once     sync.Once
)

// Instance 获取应用实例
func Instance() *App {
	return instance
}

// Setup 设置应用实例
func Setup(cfg *configs.Config) *App {
	once.Do(func() {
		instance = &App{
			ComponentManager: NewComponentManager(cfg),
			Config:           cfg,
		}
	})
	return instance
}

// Initialize 初始化应用
func (a *App) Initialize() error {
	return a.ComponentManager.Initialize()
}

// Cleanup 清理资源
func (a *App) Cleanup() {
	a.ComponentManager.Cleanup()
}

// GetConfig 获取配置
func (a *App) GetConfig() *configs.Config {
	return a.Config
}
