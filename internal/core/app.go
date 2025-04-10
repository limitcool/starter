package core

import (
	"sync"

	"github.com/limitcool/starter/configs"
)

// App 应用实例
type App struct {
	ComponentManager *ComponentManager
	AppConfig        *configs.Config
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
			AppConfig:        cfg,
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

// Config 获取配置
func (a *App) Config() *configs.Config {
	return a.AppConfig
}
