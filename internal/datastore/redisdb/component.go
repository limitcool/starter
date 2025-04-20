package redisdb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/cache"
	"github.com/limitcool/starter/internal/pkg/logger"
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
	config     *configs.Config
	manager    *Manager
	cacheStore *CacheStore
	cleanup    func()
	enabled    bool
}

// NewComponent 创建Redis组件
func NewComponent(cfg *configs.Config) *Component {
	// 检查是否存在默认配置并且显式启用
	enabled := false
	if defaultCfg, exists := cfg.Redis.Instances[DefaultRedisName]; exists && defaultCfg.Enabled {
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
		logger.Info("Redis component disabled (no configuration or not enabled)")
		return nil
	}

	logger.Info("Initializing Redis component")
	r.manager = NewManager(r.config)

	// 尝试初始化默认Redis客户端
	client, err := r.manager.GetClient(DefaultRedisName)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}

	// 初始化缓存存储
	r.cacheStore = NewCacheStore(r.manager)

	// 注册默认缓存
	// 使用配置中的默认过期时间
	_, err = r.cacheStore.RegisterCache("default", client, r.config.Redis.Cache.DefaultTTL)
	if err != nil {
		logger.Warn("Failed to register default cache", "error", err)
	}

	// 注册其他Redis实例的缓存
	for name, instance := range r.config.Redis.Instances {
		// 跳过默认实例（已经注册）和未启用的实例
		if name == DefaultRedisName || !instance.Enabled {
			continue
		}

		// 获取Redis客户端
		instanceClient, err := r.manager.GetClient(name)
		if err != nil {
			logger.Warn("Failed to get Redis client", "name", name, "error", err)
			continue
		}

		// 注册缓存
		_, err = r.cacheStore.RegisterCache(name, instanceClient, r.config.Redis.Cache.DefaultTTL)
		if err != nil {
			logger.Warn("Failed to register cache", "name", name, "error", err)
		}
	}

	r.cleanup = func() {
		// 关闭所有Redis连接
		for name, client := range r.manager.clients {
			logger.Debug("Closing Redis connection", "name", name)
			_ = client.Close()
		}

		// 关闭所有缓存
		if r.cacheStore != nil {
			r.cacheStore.Close()
		}
	}

	// 不再设置全局访问实例
	// setupInstance(r) // 已移除，使用依赖注入代替

	logger.Info("Redis component initialized successfully",
		"instances", len(r.manager.clients),
		"caches", len(r.cacheStore.caches))
	return nil
}

// Cleanup 清理Redis资源
func (r *Component) Cleanup() {
	if r.cleanup != nil {
		logger.Info("Cleaning up Redis resources")
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

// GetCache 获取默认缓存
func (r *Component) GetCache() (cache.Cache, error) {
	if !r.enabled || r.cacheStore == nil {
		return nil, fmt.Errorf("redis component not initialized or disabled")
	}
	return r.cacheStore.GetCache("default")
}

// GetCacheByName 获取指定名称的缓存
func (r *Component) GetCacheByName(name string) (cache.Cache, error) {
	if !r.enabled || r.cacheStore == nil {
		return nil, fmt.Errorf("redis component not initialized or disabled")
	}
	return r.cacheStore.GetCache(name)
}

// RegisterCache 注册新的缓存
func (r *Component) RegisterCache(name string, expiration time.Duration) (cache.Cache, error) {
	if !r.enabled || r.cacheStore == nil {
		return nil, fmt.Errorf("redis component not initialized or disabled")
	}

	// 获取默认Redis客户端
	client, err := r.GetClient()
	if err != nil {
		return nil, err
	}

	return r.cacheStore.RegisterCache(name, client, expiration)
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
	if _, exists := m.config.Redis.Instances[name]; !exists {
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

	redisConfig := m.config.Redis.Instances[name]
	logger.Debug("Creating Redis client", "name", name, "addr", redisConfig.Addr)
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
