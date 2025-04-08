package redisdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// Constants
const (
	// DefaultRedisName default redis name
	DefaultRedisName = "default"
)

// Component Redis组件实现
type Component struct {
	config  *configs.Config
	manager *Manager
	cleanup func()
	enabled bool
}

// NewComponent 创建Redis组件
func NewComponent(cfg *configs.Config) *Component {
	// 检查是否存在默认配置并且显式启用
	enabled := false
	if defaultCfg, exists := cfg.Redis[DefaultRedisName]; exists && defaultCfg.Enabled {
		enabled = true
	}

	return &Component{
		config:  cfg,
		enabled: enabled,
	}
}

// Name 返回组件名称
func (r *Component) Name() string {
	return "Redis"
}

// Initialize 初始化Redis组件
func (r *Component) Initialize() error {
	if !r.enabled {
		log.Info("Redis component disabled (no configuration or not enabled)")
		return nil
	}

	log.Info("Initializing Redis component")
	r.manager = NewManager(r.config)

	// 尝试初始化默认Redis客户端
	_, err := r.manager.GetClient(DefaultRedisName)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}

	r.cleanup = func() {
		for name, client := range r.manager.clients {
			log.Debug("Closing Redis connection", "name", name)
			_ = client.Close()
		}
	}

	// 设置全局访问实例
	setupInstance(r)

	log.Info("Redis component initialized successfully")
	return nil
}

// Cleanup 清理Redis资源
func (r *Component) Cleanup() {
	if r.cleanup != nil {
		log.Info("Cleaning up Redis resources")
		r.cleanup()
	}
}

// IsEnabled 检查组件是否启用
func (r *Component) IsEnabled() bool {
	return r.enabled
}

// GetClient 获取默认Redis客户端
func (r *Component) GetClient() (*redis.Client, error) {
	if !r.enabled || r.manager == nil {
		return nil, fmt.Errorf("redis component not initialized or disabled")
	}
	return r.manager.GetClient(DefaultRedisName)
}

// GetClientByName 获取指定名称的Redis客户端
func (r *Component) GetClientByName(name string) (*redis.Client, error) {
	if !r.enabled || r.manager == nil {
		return nil, fmt.Errorf("redis component not initialized or disabled")
	}
	return r.manager.GetClient(name)
}

// Manager Redis连接管理器
type Manager struct {
	clients map[string]*redis.Client
	config  *configs.Config
	sync.RWMutex
}

// NewManager 创建Redis管理器
func NewManager(cfg *configs.Config) *Manager {
	return &Manager{
		clients: make(map[string]*redis.Client),
		config:  cfg,
	}
}

// GetClient 获取Redis客户端
func (m *Manager) GetClient(name string) (*redis.Client, error) {
	// 检查配置是否存在
	if _, exists := m.config.Redis[name]; !exists {
		return nil, fmt.Errorf("redis configuration for '%s' not found", name)
	}

	// 从缓存获取客户端
	m.RLock()
	if client, ok := m.clients[name]; ok {
		m.RUnlock()
		return client, nil
	}
	m.RUnlock()

	// 创建新客户端
	m.Lock()
	defer m.Unlock()

	// 双重检查
	if client, ok := m.clients[name]; ok {
		return client, nil
	}

	redisConfig := m.config.Redis[name]
	log.Debug("Creating Redis client", "name", name, "addr", redisConfig.Addr)
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisConfig.Addr,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		MinIdleConns: redisConfig.MinIdleConn,
		DialTimeout:  redisConfig.DialTimeout,
		ReadTimeout:  redisConfig.ReadTimeout,
		WriteTimeout: redisConfig.WriteTimeout,
		PoolSize:     redisConfig.PoolSize,
		PoolTimeout:  redisConfig.PoolTimeout,
	})

	// 检查连接
	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	// 链路追踪
	if redisConfig.EnableTrace {
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			return nil, err
		}
	}

	m.clients[name] = rdb
	return rdb, nil
}
