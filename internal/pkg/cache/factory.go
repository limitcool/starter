package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// CacheType 缓存类型
type CacheType string

const (
	// Memory 内存缓存
	Memory CacheType = "memory"
	// Redis Redis缓存
	Redis CacheType = "redis"
)

var (
	// 缓存工厂实例
	factory *Factory
	// 缓存工厂单例锁
	factoryOnce sync.Once
)

// Factory 缓存工厂
type Factory struct {
	caches map[string]Cache
	mu     sync.RWMutex
}

// GetFactory 获取缓存工厂实例
func GetFactory() *Factory {
	factoryOnce.Do(func() {
		factory = &Factory{
			caches: make(map[string]Cache),
		}
	})
	return factory
}

// RedisOptions Redis缓存选项
type RedisOptions struct {
	Client    *redis.Client // Redis客户端
	KeyPrefix string        // 键前缀
}

// WithRedisClient 设置Redis客户端
func WithRedisClient(client *redis.Client) func(*RedisOptions) {
	return func(o *RedisOptions) {
		o.Client = client
	}
}

// WithRedisKeyPrefix 设置Redis键前缀
func WithRedisKeyPrefix(prefix string) func(*RedisOptions) {
	return func(o *RedisOptions) {
		o.KeyPrefix = prefix
	}
}

// Create 创建缓存
func (f *Factory) Create(name string, cacheType CacheType, opts ...Option) (Cache, error) {
	return f.CreateWithContext(context.Background(), name, cacheType, opts...)
}

// CreateWithContext 使用上下文创建缓存
func (f *Factory) CreateWithContext(ctx context.Context, name string, cacheType CacheType, opts ...Option) (Cache, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查缓存是否已存在
	if _, ok := f.caches[name]; ok {
		return nil, fmt.Errorf("cache: cache %s already exists", name)
	}

	// 创建缓存
	var cache Cache
	switch cacheType {
	case Memory:
		cache = NewMemoryCache(opts...)
		logger.DebugContext(ctx, "Created memory cache", "name", name)
	case Redis:
		return nil, fmt.Errorf("cache: use CreateRedis method to create Redis cache")
	default:
		return nil, fmt.Errorf("cache: unknown cache type %s", cacheType)
	}

	// 保存缓存
	f.caches[name] = cache
	return cache, nil
}

// CreateRedis 创建Redis缓存
func (f *Factory) CreateRedis(name string, client *redis.Client, opts ...Option) (Cache, error) {
	return f.CreateRedisWithContext(context.Background(), name, client, opts...)
}

// CreateRedisWithContext 使用上下文创建Redis缓存
func (f *Factory) CreateRedisWithContext(ctx context.Context, name string, client *redis.Client, opts ...Option) (Cache, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查缓存是否已存在
	if _, ok := f.caches[name]; ok {
		return nil, fmt.Errorf("cache: cache %s already exists", name)
	}

	// 创建Redis缓存
	cache := NewRedisCache(client, opts...)
	logger.DebugContext(ctx, "Created Redis cache", "name", name)

	// 保存缓存
	f.caches[name] = cache
	return cache, nil
}

// Get 获取缓存
func (f *Factory) Get(name string) (Cache, error) {
	return f.GetWithContext(context.Background(), name)
}

// GetWithContext 使用上下文获取缓存
func (f *Factory) GetWithContext(ctx context.Context, name string) (Cache, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	cache, ok := f.caches[name]
	if !ok {
		return nil, fmt.Errorf("cache: cache %s not found", name)
	}

	return cache, nil
}

// Delete 删除缓存
func (f *Factory) Delete(name string) error {
	return f.DeleteWithContext(context.Background(), name)
}

// DeleteWithContext 使用上下文删除缓存
func (f *Factory) DeleteWithContext(ctx context.Context, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	cache, ok := f.caches[name]
	if !ok {
		return fmt.Errorf("cache: cache %s not found", name)
	}

	// 关闭缓存
	if err := cache.Close(); err != nil {
		return err
	}

	// 删除缓存
	delete(f.caches, name)
	return nil
}

// Close 关闭所有缓存
func (f *Factory) Close() error {
	return f.CloseWithContext(context.Background())
}

// CloseWithContext 使用上下文关闭所有缓存
func (f *Factory) CloseWithContext(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for name, cache := range f.caches {
		if err := cache.Close(); err != nil {
			return err
		}
		delete(f.caches, name)
	}

	return nil
}
