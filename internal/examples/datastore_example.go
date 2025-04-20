package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/limitcool/starter/internal/datastore/redisdb"
	"github.com/limitcool/starter/internal/filestore"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

// DatastoreExample 展示如何使用文件存储和Redis缓存
func DatastoreExample(redisComponent *redisdb.Component, fileComponent *filestore.Component) {
	// 获取Redis缓存
	redisCache, err := redisComponent.GetCache()
	if err != nil {
		logger.Error("Failed to get Redis cache", "error", err)
		return
	}

	// 获取文件存储
	storage := fileComponent.GetStorage()
	if storage == nil {
		logger.Error("Failed to get file storage")
		return
	}

	ctx := context.Background()

	// 示例1: 使用Redis缓存存储文件元数据
	cacheFileMetadata(ctx, redisCache, storage)

	// 示例2: 使用Redis缓存实现文件访问计数
	fileAccessCounter(ctx, redisCache, storage)

	// 示例3: 使用Redis缓存实现文件热度排行
	fileHotRanking(ctx, redisCache, storage)
}

// 文件元数据
type FileMetadata struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 使用Redis缓存存储文件元数据
func cacheFileMetadata(ctx context.Context, cache interface{}, storage *filestore.Storage) {
	// 创建示例文件元数据
	metadata := FileMetadata{
		ID:        "file-123",
		Name:      "example.txt",
		Path:      "documents/example.txt",
		Size:      1024,
		MimeType:  "text/plain",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 序列化元数据
	data, err := json.Marshal(metadata)
	if err != nil {
		logger.Error("Failed to marshal file metadata", "error", err)
		return
	}

	// 设置缓存键
	key := fmt.Sprintf("file:metadata:%s", metadata.ID)

	// 将元数据存储到Redis缓存
	if c, ok := cache.(interface {
		Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	}); ok {
		if err := c.Set(ctx, key, data, 24*time.Hour); err != nil {
			logger.Error("Failed to cache file metadata", "error", err)
			return
		}
		logger.Info("File metadata cached successfully", "id", metadata.ID)
	}

	// 从Redis缓存获取元数据
	if c, ok := cache.(interface {
		Get(ctx context.Context, key string) ([]byte, error)
	}); ok {
		data, err := c.Get(ctx, key)
		if err != nil {
			logger.Error("Failed to get file metadata from cache", "error", err)
			return
		}

		// 解析元数据
		var cachedMetadata FileMetadata
		if err := json.Unmarshal(data, &cachedMetadata); err != nil {
			logger.Error("Failed to unmarshal file metadata", "error", err)
			return
		}

		logger.Info("File metadata retrieved from cache", "id", cachedMetadata.ID, "name", cachedMetadata.Name)
	}
}

// 使用Redis缓存实现文件访问计数
func fileAccessCounter(ctx context.Context, cache interface{}, storage *filestore.Storage) {
	// 文件ID
	fileID := "file-123"

	// 增加文件访问计数
	if c, ok := cache.(interface {
		Incr(ctx context.Context, key string, delta int64) (int64, error)
	}); ok {
		// 设置缓存键
		key := fmt.Sprintf("file:access_count:%s", fileID)

		// 增加访问计数
		count, err := c.Incr(ctx, key, 1)
		if err != nil {
			logger.Error("Failed to increment file access count", "error", err)
			return
		}

		logger.Info("File access count incremented", "id", fileID, "count", count)
	}
}

// 使用Redis缓存实现文件热度排行
func fileHotRanking(ctx context.Context, cache interface{}, storage *filestore.Storage) {
	// 模拟文件访问
	fileIDs := []string{"file-1", "file-2", "file-3", "file-4", "file-5"}
	accessCounts := []int{10, 5, 20, 8, 15}

	// 使用Redis的Sorted Set实现热度排行
	if redisCache, ok := cache.(*redisdb.RedisCache); ok {
		// 获取Redis客户端
		client := redisCache.GetClient()
		if client == nil {
			logger.Error("Failed to get Redis client")
			return
		}

		// 设置排行榜键
		rankingKey := "file:hot_ranking"

		// 添加文件访问计数到排行榜
		for i, fileID := range fileIDs {
			if err := client.ZAdd(ctx, rankingKey, redis.Z{
				Score:  float64(accessCounts[i]),
				Member: fileID,
			}).Err(); err != nil {
				logger.Error("Failed to add file to hot ranking", "error", err)
				continue
			}
		}

		// 获取热度排行前3的文件
		result, err := client.ZRevRangeWithScores(ctx, rankingKey, 0, 2).Result()
		if err != nil {
			logger.Error("Failed to get hot ranking", "error", err)
			return
		}

		logger.Info("File hot ranking top 3:")
		for i, z := range result {
			logger.Info(fmt.Sprintf("  %d. %s (score: %.0f)", i+1, z.Member, z.Score))
		}
	}
}
