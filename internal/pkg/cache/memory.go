package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	// ErrNotFound 缓存未找到错误
	ErrNotFound = errors.New("cache: key not found")

	// ErrKeyExists 缓存键已存在错误
	ErrKeyExists = errors.New("cache: key already exists")
)

// MemoryCache 内存缓存
type MemoryCache struct {
	cache  *cache.Cache
	mu     sync.RWMutex
	closed bool
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(opts ...Option) *MemoryCache {
	options := NewOptions(opts...)

	c := cache.New(options.Expiration, options.Expiration/2)

	if options.OnEvicted != nil {
		c.OnEvicted(func(key string, value interface{}) {
			if value != nil {
				options.OnEvicted(key, value.([]byte))
			}
		})
	}

	return &MemoryCache{
		cache: c,
	}
}

// Get 获取缓存
func (c *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, errors.New("cache: cache is closed")
	}

	value, found := c.cache.Get(key)
	if !found {
		return nil, ErrNotFound
	}

	return value.([]byte), nil
}

// Set 设置缓存
func (c *MemoryCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("cache: cache is closed")
	}

	c.cache.Set(key, value, expiration)
	return nil
}

// Delete 删除缓存
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("cache: cache is closed")
	}

	c.cache.Delete(key)
	return nil
}

// Clear 清空缓存
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("cache: cache is closed")
	}

	c.cache.Flush()
	return nil
}

// GetMulti 批量获取缓存
func (c *MemoryCache) GetMulti(ctx context.Context, keys []string) (map[string][]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, errors.New("cache: cache is closed")
	}

	result := make(map[string][]byte, len(keys))
	for _, key := range keys {
		value, found := c.cache.Get(key)
		if found {
			result[key] = value.([]byte)
		}
	}

	return result, nil
}

// SetMulti 批量设置缓存
func (c *MemoryCache) SetMulti(ctx context.Context, items map[string][]byte, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("cache: cache is closed")
	}

	for key, value := range items {
		c.cache.Set(key, value, expiration)
	}

	return nil
}

// DeleteMulti 批量删除缓存
func (c *MemoryCache) DeleteMulti(ctx context.Context, keys []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("cache: cache is closed")
	}

	for _, key := range keys {
		c.cache.Delete(key)
	}

	return nil
}

// Incr 自增
func (c *MemoryCache) Incr(ctx context.Context, key string, delta int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return 0, errors.New("cache: cache is closed")
	}

	value, err := c.cache.IncrementInt64(key, delta)
	if err != nil {
		return 0, errors.New("cache: increment failed: " + err.Error())
	}

	return value, nil
}

// Decr 自减
func (c *MemoryCache) Decr(ctx context.Context, key string, delta int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return 0, errors.New("cache: cache is closed")
	}

	value, err := c.cache.DecrementInt64(key, delta)
	if err != nil {
		return 0, errors.New("cache: decrement failed: " + err.Error())
	}

	return value, nil
}

// Exists 检查缓存是否存在
func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return false, errors.New("cache: cache is closed")
	}

	_, found := c.cache.Get(key)
	return found, nil
}

// Expire 设置过期时间
func (c *MemoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("cache: cache is closed")
	}

	item, found := c.cache.Items()[key]
	if !found {
		return ErrNotFound
	}

	item.Expiration = time.Now().Add(expiration).UnixNano()
	return nil
}

// TTL 获取过期时间
func (c *MemoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, errors.New("cache: cache is closed")
	}

	item, found := c.cache.Items()[key]
	if !found {
		return 0, ErrNotFound
	}

	if item.Expiration == 0 {
		return 0, nil
	}

	expiration := time.Unix(0, item.Expiration)
	return expiration.Sub(time.Now()), nil
}

// Close 关闭缓存
func (c *MemoryCache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("cache: cache is closed")
	}

	c.closed = true
	c.cache.Flush()
	return nil
}
