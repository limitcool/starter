package services

import (
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/storage/database"
)

// SystemService 系统服务
type SystemService struct {
	db database.DB
}

// NewSystemService 创建系统服务
func NewSystemService(db database.DB) *SystemService {
	return &SystemService{
		db: db,
	}
}

// GetSystemSettings 获取系统设置
func (s *SystemService) GetSystemSettings() map[string]any {
	// 获取配置
	config := core.Instance().Config()

	// 返回当前权限系统设置
	settings := map[string]any{
		"permission": map[string]any{
			"enabled":       config.Casbin.Enabled,
			"default_allow": config.Casbin.DefaultAllow,
		},
	}

	return settings
}
