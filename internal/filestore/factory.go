package filestore

import (
	"fmt"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/types"
)

// NewFileStorage 创建文件存储实例
func NewFileStorage(config configs.Config) (FileStorage, error) {
	if !config.Storage.Enabled {
		return nil, fmt.Errorf("文件存储未启用")
	}

	switch config.Storage.Type {
	case types.StorageTypeLocal:
		return NewLocalStorage(
			config.Storage.Local.Path,
			config.Storage.Local.URL,
		), nil

	case types.StorageTypeS3:
		minioConfig := MinIOConfig{
			Endpoint:  config.Storage.S3.Endpoint,
			Bucket:    config.Storage.S3.Bucket,
			Region:    config.Storage.S3.Region,
			AccessKey: config.Storage.S3.AccessKey,
			SecretKey: config.Storage.S3.SecretKey,
		}
		return NewMinIOStorage(minioConfig)

	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", config.Storage.Type)
	}
}

// StorageManager 存储管理器（可选，用于管理多个存储实例）
type StorageManager struct {
	primary   FileStorage
	secondary FileStorage // 可选的备用存储
}

// NewStorageManager 创建存储管理器
func NewStorageManager(primary FileStorage, secondary ...FileStorage) *StorageManager {
	manager := &StorageManager{
		primary: primary,
	}
	if len(secondary) > 0 {
		manager.secondary = secondary[0]
	}
	return manager
}

// GetPrimary 获取主存储
func (sm *StorageManager) GetPrimary() FileStorage {
	return sm.primary
}

// GetSecondary 获取备用存储
func (sm *StorageManager) GetSecondary() FileStorage {
	return sm.secondary
}

// HasSecondary 是否有备用存储
func (sm *StorageManager) HasSecondary() bool {
	return sm.secondary != nil
}
