package services

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/storage/database"
)

// SystemService 系统服务
type SystemService struct {
	db     database.DB
	config *configs.Config
}

// NewSystemService 创建系统服务
func NewSystemService(db database.DB) *SystemService {
	return &SystemService{
		db:     db,
		config: core.Instance().Config(), // 在初始化时获取配置，避免全局访问
	}
}

// GetSystemSettings 获取系统设置
func (s *SystemService) GetSystemSettings() map[string]any {
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
