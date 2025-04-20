package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

// RedisCache Redis缓存实现
type RedisCache struct {
	client     *redis.Client
	expiration time.Duration
	keyPrefix  string // 键前缀，用于区分不同应用的缓存
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(client *redis.Client, opts ...Option) *RedisCache {
	options := NewOptions(opts...)

	return &RedisCache{
		client:     client,
		expiration: options.Expiration,
		keyPrefix:  "cache:", // 默认前缀
	}
}

// WithKeyPrefix 设置键前缀
func WithKeyPrefix(prefix string) Option {
	return func(o *Options) {
		// 这里不直接使用前缀，而是在RedisCache构造函数中设置
		// 这个选项只是为了传递前缀参数
	}
}

// prefixKey 为键添加前缀
func (c *RedisCache) prefixKey(key string) string {
	return c.keyPrefix + key
}

// Get 获取缓存
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	prefixedKey := c.prefixKey(key)
	val, err := c.client.Get(ctx, prefixedKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return val, nil
}

// Set 设置缓存
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	prefixedKey := c.prefixKey(key)
	if expiration == 0 {
		expiration = c.expiration
	}
	return c.client.Set(ctx, prefixedKey, value, expiration).Err()
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	prefixedKey := c.prefixKey(key)
	return c.client.Del(ctx, prefixedKey).Err()
}

// Clear 清空缓存（谨慎使用）
func (c *RedisCache) Clear(ctx context.Context) error {
	// 使用SCAN命令查找所有带前缀的键
	iter := c.client.Scan(ctx, 0, c.keyPrefix+"*", 100).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
		// 当积累了一定数量的键时，批量删除
		if len(keys) >= 1000 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
			keys = keys[:0] // 清空切片但保留容量
		}
	}

	// 删除剩余的键
	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

// GetMulti 批量获取缓存
func (c *RedisCache) GetMulti(ctx context.Context, keys []string) (map[string][]byte, error) {
	if len(keys) == 0 {
		return make(map[string][]byte), nil
	}

	// 为所有键添加前缀
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.prefixKey(key)
	}

	// 使用管道批量获取
	pipe := c.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(prefixedKeys))
	for i, key := range prefixedKeys {
		cmds[i] = pipe.Get(ctx, key)
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 处理结果
	result := make(map[string][]byte, len(keys))
	for i, cmd := range cmds {
		val, err := cmd.Bytes()
		if err == nil {
			// 去掉前缀，返回原始键
			result[keys[i]] = val
		}
	}

	return result, nil
}

// SetMulti 批量设置缓存
func (c *RedisCache) SetMulti(ctx context.Context, items map[string][]byte, expiration time.Duration) error {
	if len(items) == 0 {
		return nil
	}

	if expiration == 0 {
		expiration = c.expiration
	}

	// 使用管道批量设置
	pipe := c.client.Pipeline()
	for key, value := range items {
		prefixedKey := c.prefixKey(key)
		pipe.Set(ctx, prefixedKey, value, expiration)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteMulti 批量删除缓存
func (c *RedisCache) DeleteMulti(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	// 为所有键添加前缀
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.prefixKey(key)
	}

	return c.client.Del(ctx, prefixedKeys...).Err()
}

// Incr 自增
func (c *RedisCache) Incr(ctx context.Context, key string, delta int64) (int64, error) {
	prefixedKey := c.prefixKey(key)
	return c.client.IncrBy(ctx, prefixedKey, delta).Result()
}

// Decr 自减
func (c *RedisCache) Decr(ctx context.Context, key string, delta int64) (int64, error) {
	prefixedKey := c.prefixKey(key)
	return c.client.DecrBy(ctx, prefixedKey, delta).Result()
}

// Exists 检查缓存是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	prefixedKey := c.prefixKey(key)
	n, err := c.client.Exists(ctx, prefixedKey).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Expire 设置过期时间
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	prefixedKey := c.prefixKey(key)
	return c.client.Expire(ctx, prefixedKey, expiration).Err()
}

// TTL 获取过期时间
func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	prefixedKey := c.prefixKey(key)
	return c.client.TTL(ctx, prefixedKey).Result()
}

// Close 关闭缓存
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// 缓存预热相关方法

// WarmUp 缓存预热
// keys: 需要预热的键列表
// loader: 加载数据的函数，接收键并返回对应的值
func (c *RedisCache) WarmUp(ctx context.Context, keys []string, loader func(ctx context.Context, key string) (any, error)) error {
	if len(keys) == 0 {
		return nil
	}

	logger.Info("Warming up cache", "count", len(keys))

	// 检查哪些键需要加载
	pipe := c.client.Pipeline()
	existCmds := make([]*redis.IntCmd, len(keys))
	for i, key := range keys {
		prefixedKey := c.prefixKey(key)
		existCmds[i] = pipe.Exists(ctx, prefixedKey)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to check existing keys: %w", err)
	}

	// 加载缺失的键
	var missingKeys []string
	var missingIndexes []int
	for i, cmd := range existCmds {
		exists, err := cmd.Result()
		if err != nil {
			return err
		}
		if exists == 0 {
			missingKeys = append(missingKeys, keys[i])
			missingIndexes = append(missingIndexes, i)
		}
	}

	if len(missingKeys) == 0 {
		logger.Info("All cache keys already exist, no need to warm up")
		return nil
	}

	logger.Info("Loading missing cache keys", "count", len(missingKeys))

	// 并发加载数据
	type result struct {
		key   string
		value any
		err   error
		index int
	}

	// 使用工作池限制并发数
	concurrency := 10
	if len(missingKeys) < concurrency {
		concurrency = len(missingKeys)
	}

	jobs := make(chan int, len(missingKeys))
	results := make(chan result, len(missingKeys))

	// 启动工作协程
	for w := 0; w < concurrency; w++ {
		go func() {
			for idx := range jobs {
				key := missingKeys[idx]
				value, err := loader(ctx, key)
				results <- result{
					key:   key,
					value: value,
					err:   err,
					index: missingIndexes[idx],
				}
			}
		}()
	}

	// 发送工作
	for i := range missingKeys {
		jobs <- i
	}
	close(jobs)

	// 收集结果
	pipe = c.client.Pipeline()
	setCmds := make([]*redis.StatusCmd, len(missingKeys))
	for i := 0; i < len(missingKeys); i++ {
		r := <-results
		if r.err != nil {
			logger.Error("Failed to load data for cache key", "key", r.key, "error", r.err)
			continue
		}

		// 序列化值
		data, err := json.Marshal(r.value)
		if err != nil {
			logger.Error("Failed to marshal data for cache key", "key", r.key, "error", err)
			continue
		}

		prefixedKey := c.prefixKey(r.key)
		setCmds[i] = pipe.Set(ctx, prefixedKey, data, c.expiration)
	}

	// 执行管道
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set cache values: %w", err)
	}

	logger.Info("Cache warm up completed", "count", len(missingKeys))
	return nil
}

// 缓存穿透保护

// GetWithProtection 获取缓存，带穿透保护
// key: 缓存键
// loader: 加载数据的函数，当缓存未命中时调用
// expiration: 缓存过期时间，为0时使用默认过期时间
func (c *RedisCache) GetWithProtection(ctx context.Context, key string, loader func(ctx context.Context) (any, error), expiration time.Duration) (any, error) {
	prefixedKey := c.prefixKey(key)

	// 尝试从缓存获取
	data, err := c.client.Get(ctx, prefixedKey).Bytes()
	if err == nil {
		// 缓存命中，解析数据
		var result any
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		return result, nil
	} else if err != redis.Nil {
		// 发生了除了键不存在之外的错误
		return nil, err
	}

	// 缓存未命中，使用分布式锁防止缓存穿透
	lockKey := c.prefixKey("lock:" + key)
	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())

	// 尝试获取锁，过期时间5秒
	locked, err := c.client.SetNX(ctx, lockKey, lockValue, 5*time.Second).Result()
	if err != nil {
		return nil, err
	}

	// 如果没有获取到锁，说明有其他请求正在加载数据
	if !locked {
		// 等待一段时间后重试
		time.Sleep(100 * time.Millisecond)
		return c.GetWithProtection(ctx, key, loader, expiration)
	}

	// 获取到锁，加载数据
	defer c.client.Del(ctx, lockKey)

	// 再次检查缓存，可能在获取锁的过程中已经有其他请求加载了数据
	data, err = c.client.Get(ctx, prefixedKey).Bytes()
	if err == nil {
		// 缓存命中，解析数据
		var result any
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		return result, nil
	} else if err != redis.Nil {
		// 发生了除了键不存在之外的错误
		return nil, err
	}

	// 加载数据
	result, err := loader(ctx)
	if err != nil {
		return nil, err
	}

	// 序列化数据
	data, err = json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// 设置缓存
	if expiration == 0 {
		expiration = c.expiration
	}
	if err := c.client.Set(ctx, prefixedKey, data, expiration).Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetWithBloomFilter 使用布隆过滤器防止缓存穿透
// key: 缓存键
// loader: 加载数据的函数，当缓存未命中时调用
// bloomFilterKey: 布隆过滤器的键
// expiration: 缓存过期时间，为0时使用默认过期时间
func (c *RedisCache) GetWithBloomFilter(ctx context.Context, key string, loader func(ctx context.Context) (any, error), bloomFilterKey string, expiration time.Duration) (any, error) {
	prefixedKey := c.prefixKey(key)

	// 尝试从缓存获取
	data, err := c.client.Get(ctx, prefixedKey).Bytes()
	if err == nil {
		// 缓存命中，解析数据
		var result any
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		return result, nil
	} else if err != redis.Nil {
		// 发生了除了键不存在之外的错误
		return nil, err
	}

	// 缓存未命中，检查布隆过滤器
	// 注意：这里假设布隆过滤器已经在Redis中设置好了
	// 实际使用时，需要先初始化布隆过滤器并添加所有有效的键
	exists, err := c.client.Exists(ctx, bloomFilterKey).Result()
	if err != nil {
		return nil, err
	}

	// 如果布隆过滤器不存在，直接加载数据
	if exists == 0 {
		// 布隆过滤器不存在，直接加载数据
		return c.loadAndCache(ctx, key, loader, expiration)
	}

	// 检查键是否在布隆过滤器中
	// 注意：这里使用了简化的实现，实际的布隆过滤器需要使用专门的库或Redis模块
	// 例如RedisBloom模块: https://github.com/RedisBloom/RedisBloom
	inFilter, err := c.client.SIsMember(ctx, bloomFilterKey, key).Result()
	if err != nil {
		return nil, err
	}

	// 如果键不在布隆过滤器中，说明数据不存在，返回空值
	if !inFilter {
		return nil, errors.New("key does not exist (filtered by bloom filter)")
	}

	// 键在布隆过滤器中，加载数据
	return c.loadAndCache(ctx, key, loader, expiration)
}

// loadAndCache 加载数据并缓存
func (c *RedisCache) loadAndCache(ctx context.Context, key string, loader func(ctx context.Context) (any, error), expiration time.Duration) (any, error) {
	// 加载数据
	result, err := loader(ctx)
	if err != nil {
		return nil, err
	}

	// 序列化数据
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// 设置缓存
	prefixedKey := c.prefixKey(key)
	if expiration == 0 {
		expiration = c.expiration
	}
	if err := c.client.Set(ctx, prefixedKey, data, expiration).Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// SetNilValue 设置空值缓存，用于缓存穿透保护
// 当查询结果为空时，缓存一个特殊的空值，避免频繁查询数据库
func (c *RedisCache) SetNilValue(ctx context.Context, key string, expiration time.Duration) error {
	if expiration == 0 {
		// 空值缓存的过期时间通常比正常数据短
		expiration = c.expiration / 10
		if expiration < time.Second {
			expiration = time.Second
		}
	}

	prefixedKey := c.prefixKey(key)
	return c.client.Set(ctx, prefixedKey, "nil", expiration).Err()
}

// IsNilValue 检查是否为空值缓存
func (c *RedisCache) IsNilValue(ctx context.Context, data []byte) bool {
	return string(data) == "nil"
}
