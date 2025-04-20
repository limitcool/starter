package services

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/datastore/database"
)

// SystemService 系统服务
type SystemService struct {
	db     database.Database
	config *configs.Config
}

// NewSystemService 创建系统服务
func NewSystemService(db database.Database, config *configs.Config) *SystemService {
	return &SystemService{
		db:     db,
		config: config,
	}
}

// GetSystemSettings 获取系统设置
func (s *SystemService) GetSystemSettings(ctx context.Context) map[string]any {
	// 使用服务实例中的配置
	config := s.config

	// 返回当前权限系统设置
	settings := map[string]any{
		"permission": map[string]any{
			"enabled":       config.Casbin.Enabled,
			"default_allow": config.Casbin.DefaultAllow,
		},
	}

	return settings
}
