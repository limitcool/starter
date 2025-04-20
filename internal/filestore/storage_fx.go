package filestore

import (
	"context"
	"fmt"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
)

// NewStorage 创建存储服务
func NewStorage(lc fx.Lifecycle, cfg *configs.Config) (*Storage, error) {
	if !cfg.Storage.Enabled {
		logger.Info("File storage disabled")
		return nil, nil
	}

	logger.Info("Initializing storage service", "type", cfg.Storage.Type)

	// 创建存储配置
	storageConfig := Config{
		Type: cfg.Storage.Type,
	}

	// 根据存储类型设置配置
	switch cfg.Storage.Type {
	case "local":
		storageConfig.Path = cfg.Storage.Local.Path
		storageConfig.URL = cfg.Storage.Local.URL
	case "s3":
		storageConfig.AccessKey = cfg.Storage.S3.AccessKey
		storageConfig.SecretKey = cfg.Storage.S3.SecretKey
		storageConfig.Region = cfg.Storage.S3.Region
		storageConfig.Bucket = cfg.Storage.S3.Bucket
		storageConfig.Endpoint = cfg.Storage.S3.Endpoint
	}

	// 创建存储服务
	storage, err := New(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage service: %w", err)
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Storage service initialized successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Cleaning up storage service")
			return nil
		},
	})

	return storage, nil
}
