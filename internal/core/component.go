package core

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/i18n"
)

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
	// 记录初始化开始时间
	startTime := time.Now()

	// 初始化国际化组件
	if m.config.I18n.Enabled {
		log.Info("Initializing internationalization component...")
		err := initI18n(m.config.I18n)
		if err != nil {
			log.Error("Failed to initialize i18n component", "err", err)
			return err
		}
		log.Info("Internationalization component initialized successfully")
	}

	// 初始化其他组件
	for _, c := range m.components {
		log.Info("Initializing component", "component", c.Name())
		if err := c.Initialize(); err != nil {
			log.Error("Failed to initialize component", "component", c.Name(), "err", err)
			return err
		}
		log.Info("Component initialized successfully", "component", c.Name())
	}

	// 记录初始化耗时
	elapsedTime := time.Since(startTime)
	log.Info("All components initialized", "duration", elapsedTime)

	return nil
}

// 初始化国际化组件
func initI18n(config configs.I18n) error {
	// 打开本地资源目录
	fsys := os.DirFS(config.ResourcesPath)

	// 配置i18n
	i18nConfig := i18n.Config{
		DefaultLanguage:  config.DefaultLanguage,
		SupportLanguages: config.SupportLanguages,
	}

	// 初始化i18n服务
	return i18n.Setup(i18nConfig, fsys)
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
