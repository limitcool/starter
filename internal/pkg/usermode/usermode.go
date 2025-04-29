package usermode

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// Service 用户模式服务
type Service struct {
	config *configs.Config
	mode   enum.UserMode
}

// NewService 创建用户模式服务
func NewService(config *configs.Config) *Service {
	// 在lite版本中，始终使用简单模式
	mode := enum.UserModeSimple
	logger.Info("初始化用户模式服务", "user_mode", mode)
	return &Service{
		config: config,
		mode:   mode,
	}
}

// GetMode 获取用户模式
func (s *Service) GetMode() enum.UserMode {
	return s.mode
}

// IsSimpleMode 是否为简单模式
func (s *Service) IsSimpleMode() bool {
	return s.mode == enum.UserModeSimple
}

// IsSeparateMode 是否为分离模式
func (s *Service) IsSeparateMode() bool {
	return s.mode == enum.UserModeSeparate
}

// LogMode 记录用户模式
func (s *Service) LogMode(ctx context.Context) {
	logger.InfoContext(ctx, "当前用户模式", "mode", s.mode)
}
