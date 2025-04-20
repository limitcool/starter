package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/limitcool/starter/internal/datastore/redisdb"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/cache"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
)

// RedisCacheExample 展示如何使用Redis缓存
func RedisCacheExample(redisComponent *redisdb.Component, userRepo repository.UserRepository) {
	// 获取Redis缓存
	redisCache, err := redisComponent.GetCache()
	if err != nil {
		logger.Error("Failed to get Redis cache", "error", err)
		return
	}

	ctx := context.Background()

	// 示例1: 基本的缓存操作
	basicCacheOperations(ctx, redisCache)

	// 示例2: 缓存穿透保护
	cachePenetrationProtection(ctx, redisCache, userRepo)

	// 示例3: 缓存预热
	cacheWarmUp(ctx, redisCache, userRepo)
}

// 基本的缓存操作
func basicCacheOperations(ctx context.Context, c cache.Cache) {
	// 设置缓存
	key := "example:basic:1"
	value := map[string]any{
		"name":  "Example User",
		"email": "example@example.com",
		"age":   30,
	}

	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("Failed to marshal value", "error", err)
		return
	}

	// 设置缓存，过期时间5分钟
	if err := c.Set(ctx, key, data, 5*time.Minute); err != nil {
		logger.Error("Failed to set cache", "error", err)
		return
	}

	logger.Info("Cache set successfully", "key", key)

	// 获取缓存
	data, err = c.Get(ctx, key)
	if err != nil {
		logger.Error("Failed to get cache", "error", err)
		return
	}

	// 解析值
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		logger.Error("Failed to unmarshal value", "error", err)
		return
	}

	logger.Info("Cache get successfully", "key", key, "value", result)

	// 检查缓存是否存在
	exists, err := c.Exists(ctx, key)
	if err != nil {
		logger.Error("Failed to check cache existence", "error", err)
		return
	}

	logger.Info("Cache exists", "key", key, "exists", exists)

	// 获取过期时间
	ttl, err := c.TTL(ctx, key)
	if err != nil {
		logger.Error("Failed to get cache TTL", "error", err)
		return
	}

	logger.Info("Cache TTL", "key", key, "ttl", ttl)

	// 删除缓存
	if err := c.Delete(ctx, key); err != nil {
		logger.Error("Failed to delete cache", "error", err)
		return
	}

	logger.Info("Cache deleted successfully", "key", key)

	// 批量操作
	items := map[string][]byte{
		"example:batch:1": []byte(`{"id": 1, "name": "User 1"}`),
		"example:batch:2": []byte(`{"id": 2, "name": "User 2"}`),
		"example:batch:3": []byte(`{"id": 3, "name": "User 3"}`),
	}

	// 批量设置缓存
	if err := c.SetMulti(ctx, items, 5*time.Minute); err != nil {
		logger.Error("Failed to set multiple caches", "error", err)
		return
	}

	logger.Info("Multiple caches set successfully")

	// 批量获取缓存
	keys := []string{"example:batch:1", "example:batch:2", "example:batch:3"}
	results, err := c.GetMulti(ctx, keys)
	if err != nil {
		logger.Error("Failed to get multiple caches", "error", err)
		return
	}

	logger.Info("Multiple caches get successfully", "count", len(results))

	// 批量删除缓存
	if err := c.DeleteMulti(ctx, keys); err != nil {
		logger.Error("Failed to delete multiple caches", "error", err)
		return
	}

	logger.Info("Multiple caches deleted successfully")
}

// 缓存穿透保护
func cachePenetrationProtection(ctx context.Context, c cache.Cache, userRepo repository.UserRepository) {
	// 检查是否为Redis缓存
	redisCache, ok := c.(*cache.RedisCache)
	if !ok {
		logger.Error("Cache is not a Redis cache")
		return
	}

	// 使用缓存穿透保护获取用户
	userID := int64(1)
	cacheKey := fmt.Sprintf("user:%d", userID)

	// 定义加载函数
	loader := func(ctx context.Context) (any, error) {
		logger.Info("Loading user from database", "id", userID)
		// 从数据库加载用户
		user, err := userRepo.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	// 使用缓存穿透保护获取用户
	result, err := redisCache.GetWithProtection(ctx, cacheKey, loader, 5*time.Minute)
	if err != nil {
		logger.Error("Failed to get user with cache protection", "error", err)
		return
	}

	// 转换为用户对象
	user, ok := result.(*model.User)
	if !ok {
		logger.Error("Failed to convert result to user")
		return
	}

	logger.Info("Got user with cache protection", "id", user.ID, "username", user.Username)
}

// 缓存预热
func cacheWarmUp(ctx context.Context, c cache.Cache, userRepo repository.UserRepository) {
	// 检查是否为Redis缓存
	redisCache, ok := c.(*cache.RedisCache)
	if !ok {
		logger.Error("Cache is not a Redis cache")
		return
	}

	// 获取需要预热的用户ID列表
	userIDs := []int64{1, 2, 3, 4, 5}

	// 准备缓存键
	keys := make([]string, len(userIDs))
	for i, id := range userIDs {
		keys[i] = fmt.Sprintf("user:%d", id)
	}

	// 定义加载函数
	loader := func(ctx context.Context, key string) (any, error) {
		// 从键中提取用户ID
		var id int64
		_, err := fmt.Sscanf(key, "user:%d", &id)
		if err != nil {
			return nil, err
		}

		logger.Info("Preloading user from database", "id", id)
		// 从数据库加载用户
		user, err := userRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	// 预热缓存
	if err := redisCache.WarmUp(ctx, keys, loader); err != nil {
		logger.Error("Failed to warm up cache", "error", err)
		return
	}

	logger.Info("Cache warm up completed successfully")
}
