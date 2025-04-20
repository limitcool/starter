package redisdb

import (
	"context"
	"time"

	"github.com/limitcool/starter/internal/pkg/cache"
	"github.com/redis/go-redis/v9"
)

// RedisCache Redis缓存适配器
// 实现cache.Cache接口，但提供额外的Redis特定功能
type RedisCache struct {
	cache  cache.Cache
	client *redis.Client
}

// NewRedisCache 创建Redis缓存适配器
func NewRedisCache(client *redis.Client, expiration time.Duration) *RedisCache {
	// 创建基础缓存
	baseCache := cache.NewRedisCache(
		client,
		cache.WithExpiration(expiration),
	)

	return &RedisCache{
		cache:  baseCache,
		client: client,
	}
}

// Get 获取缓存
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.cache.Get(ctx, key)
}

// Set 设置缓存
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return c.cache.Set(ctx, key, value, expiration)
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, key)
}

// Clear 清空缓存
func (c *RedisCache) Clear(ctx context.Context) error {
	return c.cache.Clear(ctx)
}

// GetMulti 批量获取缓存
func (c *RedisCache) GetMulti(ctx context.Context, keys []string) (map[string][]byte, error) {
	return c.cache.GetMulti(ctx, keys)
}

// SetMulti 批量设置缓存
func (c *RedisCache) SetMulti(ctx context.Context, items map[string][]byte, expiration time.Duration) error {
	return c.cache.SetMulti(ctx, items, expiration)
}

// DeleteMulti 批量删除缓存
func (c *RedisCache) DeleteMulti(ctx context.Context, keys []string) error {
	return c.cache.DeleteMulti(ctx, keys)
}

// Incr 自增
func (c *RedisCache) Incr(ctx context.Context, key string, delta int64) (int64, error) {
	return c.cache.Incr(ctx, key, delta)
}

// Decr 自减
func (c *RedisCache) Decr(ctx context.Context, key string, delta int64) (int64, error) {
	return c.cache.Decr(ctx, key, delta)
}

// Exists 检查缓存是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	return c.cache.Exists(ctx, key)
}

// Expire 设置过期时间
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.cache.Expire(ctx, key, expiration)
}

// TTL 获取过期时间
func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.cache.TTL(ctx, key)
}

// Close 关闭缓存
func (c *RedisCache) Close() error {
	return c.cache.Close()
}

// GetClient 获取Redis客户端
// 这是RedisCache特有的方法，不是cache.Cache接口的一部分
func (c *RedisCache) GetClient() *redis.Client {
	return c.client
}

// WarmUp 缓存预热
// 这是RedisCache特有的方法，不是cache.Cache接口的一部分
func (c *RedisCache) WarmUp(ctx context.Context, keys []string, loader func(ctx context.Context, key string) (any, error)) error {
	// 检查是否为Redis缓存
	redisCache, ok := c.cache.(*cache.RedisCache)
	if !ok {
		return nil
	}
	return redisCache.WarmUp(ctx, keys, loader)
}

// GetWithProtection 获取缓存，带穿透保护
// 这是RedisCache特有的方法，不是cache.Cache接口的一部分
func (c *RedisCache) GetWithProtection(ctx context.Context, key string, loader func(ctx context.Context) (any, error), expiration time.Duration) (any, error) {
	// 检查是否为Redis缓存
	redisCache, ok := c.cache.(*cache.RedisCache)
	if !ok {
		return nil, nil
	}
	return redisCache.GetWithProtection(ctx, key, loader, expiration)
}
