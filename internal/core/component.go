package core

import "github.com/limitcool/starter/configs"

// Component 定义应用组件接口
type Component interface {
	// Initialize 初始化组件
	Initialize() error
	// Cleanup 清理组件资源
	Cleanup()
	// Name 返回组件名称（用于日志等）
	Name() string
}

// ComponentManager 管理应用组件
type ComponentManager struct {
	components []Component
	config     *configs.Config
}

// NewComponentManager 创建组件管理器
func NewComponentManager(cfg *configs.Config) *ComponentManager {
	return &ComponentManager{
		components: make([]Component, 0),
		config:     cfg,
	}
}

// AddComponent 添加组件到管理器
func (m *ComponentManager) AddComponent(c Component) {
	m.components = append(m.components, c)
}

// Initialize 初始化所有组件
func (m *ComponentManager) Initialize() error {
	for _, c := range m.components {
		if err := c.Initialize(); err != nil {
			return err
		}
	}
	return nil
}

// Cleanup 清理所有组件资源
func (m *ComponentManager) Cleanup() {
	// 反向顺序清理，确保依赖组件后清理
	for i := len(m.components) - 1; i >= 0; i-- {
		m.components[i].Cleanup()
	}
}

// GetConfig 获取配置
func (m *ComponentManager) GetConfig() *configs.Config {
	return m.config
}
