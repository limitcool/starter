package cache

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) ([]byte, error)
	
	// Set 设置缓存
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	
	// Delete 删除缓存
	Delete(ctx context.Context, key string) error
	
	// Clear 清空缓存
	Clear(ctx context.Context) error
	
	// GetMulti 批量获取缓存
	GetMulti(ctx context.Context, keys []string) (map[string][]byte, error)
	
	// SetMulti 批量设置缓存
	SetMulti(ctx context.Context, items map[string][]byte, expiration time.Duration) error
	
	// DeleteMulti 批量删除缓存
	DeleteMulti(ctx context.Context, keys []string) error
	
	// Incr 自增
	Incr(ctx context.Context, key string, delta int64) (int64, error)
	
	// Decr 自减
	Decr(ctx context.Context, key string, delta int64) (int64, error)
	
	// Exists 检查缓存是否存在
	Exists(ctx context.Context, key string) (bool, error)
	
	// Expire 设置过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) error
	
	// TTL 获取过期时间
	TTL(ctx context.Context, key string) (time.Duration, error)
	
	// Close 关闭缓存
	Close() error
}

// Options 缓存选项
type Options struct {
	// Expiration 默认过期时间
	Expiration time.Duration
	
	// MaxEntries 最大缓存条目数
	MaxEntries int
	
	// OnEvicted 缓存淘汰回调函数
	OnEvicted func(key string, value []byte)
}

// DefaultOptions 默认缓存选项
var DefaultOptions = Options{
	Expiration: 5 * time.Minute,
	MaxEntries: 10000,
}

// NewOptions 创建缓存选项
func NewOptions(opts ...Option) Options {
	options := DefaultOptions
	
	for _, opt := range opts {
		opt(&options)
	}
	
	return options
}

// Option 缓存选项函数
type Option func(*Options)

// WithExpiration 设置默认过期时间
func WithExpiration(expiration time.Duration) Option {
	return func(o *Options) {
		o.Expiration = expiration
	}
}

// WithMaxEntries 设置最大缓存条目数
func WithMaxEntries(maxEntries int) Option {
	return func(o *Options) {
		o.MaxEntries = maxEntries
	}
}

// WithOnEvicted 设置缓存淘汰回调函数
func WithOnEvicted(onEvicted func(key string, value []byte)) Option {
	return func(o *Options) {
		o.OnEvicted = onEvicted
	}
}
