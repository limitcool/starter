package redisdb

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/cache"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
)

// DefaultRedisName 默认Redis实例名称
const DefaultRedisName = "default"

// Module Redis模块
var Module = fx.Options(
	// 提供Redis客户端
	fx.Provide(NewRedisClient),

	// 提供缓存
	fx.Provide(NewRedisCache),
)

// NewRedisClient 创建Redis客户端
func NewRedisClient(lc fx.Lifecycle, cfg *configs.Config) (*redis.Client, error) {
	// 检查是否存在默认配置并且显式启用
	var redisConfig configs.RedisInstance
	var exists bool

	if redisConfig, exists = cfg.Redis.Instances[DefaultRedisName]; !exists || !redisConfig.Enabled {
		logger.Info("Redis disabled")
		return nil, nil
	}

	logger.Info("Connecting to Redis", "addr", redisConfig.Addr)

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
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

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 检查Redis连接
			if err := client.Ping(ctx).Err(); err != nil {
				return fmt.Errorf("failed to ping Redis: %w", err)
			}

			logger.Info("Redis connected successfully")
			return nil
		},
	})

	return client, nil
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(client *redis.Client, cfg *configs.Config) (cache.Cache, error) {
	if client == nil {
		logger.Info("Redis cache disabled")
		return nil, nil
	}

	logger.Info("Creating Redis cache")

	// 创建Redis缓存
	redisCache := cache.NewRedisCache(
		client,
		cache.WithExpiration(cfg.Redis.Cache.DefaultTTL),
		cache.WithKeyPrefix(cfg.Redis.Cache.KeyPrefix),
	)

	return redisCache, nil
}
