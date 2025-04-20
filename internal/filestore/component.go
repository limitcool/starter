package filestore

import (
	"fmt"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/types"
)

// Component 文件存储组件
type Component struct {
	config  *configs.Config
	storage *Storage
	enabled bool
}

// NewComponent 创建文件存储组件
func NewComponent(cfg *configs.Config) *Component {
	return &Component{
		config:  cfg,
		enabled: cfg.Storage.Enabled,
	}
}

// Name 返回组件名称
func (c *Component) Name() string {
	return "FileStorage"
}

// Initialize 初始化文件存储组件
func (c *Component) Initialize() error {
	if !c.enabled {
		logger.Info("File storage component disabled")
		return nil
	}

	logger.Info("Initializing file storage component")

	// 创建存储配置
	storageConfig := Config{Type: c.config.Storage.Type}

	// 根据存储类型设置配置
	switch c.config.Storage.Type {
	case types.StorageTypeLocal:
		storageConfig.Path = c.config.Storage.Local.Path
		storageConfig.URL = c.config.Storage.Local.URL
	case types.StorageTypeS3:
		storageConfig.AccessKey = c.config.Storage.S3.AccessKey
		storageConfig.SecretKey = c.config.Storage.S3.SecretKey
		storageConfig.Region = c.config.Storage.S3.Region
		storageConfig.Bucket = c.config.Storage.S3.Bucket
		storageConfig.Endpoint = c.config.Storage.S3.Endpoint
	default:
		return fmt.Errorf("unsupported storage type: %s", c.config.Storage.Type)
	}

	// 创建存储服务
	storage, err := New(storageConfig)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	c.storage = storage

	logger.Info("File storage component initialized successfully")
	return nil
}

// Cleanup 清理文件存储资源
func (c *Component) Cleanup() {
	if c.enabled {
		logger.Info("Cleaning up file storage resources")
	}
}

// IsEnabled 检查组件是否启用
func (c *Component) IsEnabled() bool {
	return c.enabled
}

// GetStorage 获取存储服务
func (c *Component) GetStorage() *Storage {
	return c.storage
}
