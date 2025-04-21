package services

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"go.uber.org/fx"
)

// AdminService 管理系统服务
type AdminService struct {
	adminRepo *repository.AdminRepo
	config    *configs.Config
}

// NewAdminService 创建管理系统服务
func NewAdminService(params ServiceParams) *AdminService {
	service := &AdminService{
		adminRepo: params.AdminRepo,
		config:    params.Config,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("AdminService initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("AdminService stopped")
			return nil
		},
	})

	return service
}

// GetSystemSettings 获取系统设置
func (s *AdminService) GetSystemSettings(ctx context.Context) map[string]any {
	// 使用服务实例中的配置
	config := s.config

	// 从仓库中获取系统设置
	dbSettings, err := s.adminRepo.GetSystemSettings(ctx)
	if err != nil {
		logger.Error("Failed to get system settings from database", "error", err)
		// 如果出错，使用配置文件中的设置
		dbSettings = map[string]any{}
	}

	// 合并配置文件和数据库中的设置
	// 默认使用配置文件中的设置
	settings := map[string]any{
		"permission": map[string]any{
			"enabled":       config.Casbin.Enabled,
			"default_allow": config.Casbin.DefaultAllow,
		},
	}

	// 如果数据库中有设置，则使用数据库中的设置
	for k, v := range dbSettings {
		settings[k] = v
	}

	return settings
}
