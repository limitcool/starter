package redisdb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/limitcool/starter/internal/pkg/cache"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

// CacheStore Redis缓存存储
type CacheStore struct {
	manager *Manager
	caches  map[string]cache.Cache
	mu      sync.RWMutex
}

// NewCacheStore 创建缓存存储
func NewCacheStore(manager *Manager) *CacheStore {
	return &CacheStore{
		manager: manager,
		caches:  make(map[string]cache.Cache),
	}
}

// RegisterCache 注册缓存
func (s *CacheStore) RegisterCache(name string, client *redis.Client, expiration time.Duration) (cache.Cache, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查缓存是否已存在
	if _, ok := s.caches[name]; ok {
		return nil, fmt.Errorf("cache %s already registered", name)
	}

	// 获取缓存配置
	cacheConfig := s.manager.config.Redis.Cache

	// 如果未指定过期时间，使用配置中的默认值
	if expiration == 0 {
		expiration = cacheConfig.DefaultTTL
	}

	// 创建Redis缓存
	redisCache := NewRedisCache(client, expiration)

	// 保存缓存
	s.caches[name] = redisCache
	logger.Debug("Registered Redis cache", "name", name, "expiration", expiration)
	return redisCache, nil
}

// GetCache 获取缓存
func (s *CacheStore) GetCache(name string) (cache.Cache, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查缓存是否存在
	c, ok := s.caches[name]
	if !ok {
		return nil, fmt.Errorf("cache %s not found", name)
	}

	return c, nil
}

// Close 关闭所有缓存
func (s *CacheStore) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, c := range s.caches {
		logger.Debug("Closing cache", "name", name)
		_ = c.Close()
	}

	// 清空缓存映射
	s.caches = make(map[string]cache.Cache)
}

// WarmUpCache 预热缓存
func (s *CacheStore) WarmUpCache(name string, keys []string, loader func(key string) (any, error)) error {
	// 获取缓存
	c, err := s.GetCache(name)
	if err != nil {
		return err
	}

	// 检查是否为Redis缓存
	redisCache, ok := c.(*RedisCache)
	if !ok {
		return fmt.Errorf("cache %s is not a Redis cache", name)
	}

	// 预热缓存
	return redisCache.WarmUp(context.Background(), keys, func(ctx context.Context, key string) (any, error) {
		return loader(key)
	})
}
