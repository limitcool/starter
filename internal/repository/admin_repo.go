package repository

import (
	"context"

	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// AdminRepo 管理系统仓库
type AdminRepo struct {
	DB *gorm.DB
}

// NewAdminRepo 创建管理系统仓库
func NewAdminRepo(params RepoParams) *AdminRepo {
	repo := &AdminRepo{
		DB: params.DB,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("AdminRepo initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("AdminRepo stopped")
			return nil
		},
	})

	return repo
}

// GetSystemSettings 获取系统设置
func (r *AdminRepo) GetSystemSettings(ctx context.Context) (map[string]any, error) {
	// 这里可以从数据库中获取系统设置
	// 目前我们返回一个空的map，因为AdminService会使用配置文件中的设置
	return map[string]any{}, nil
}
