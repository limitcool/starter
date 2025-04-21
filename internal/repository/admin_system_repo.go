package repository

import (
	"context"

	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// AdminSystemRepo 系统仓库
type AdminSystemRepo struct {
	DB *gorm.DB
}

// NewAdminSystemRepo 创建系统仓库
func NewAdminSystemRepo(params RepoParams) *AdminSystemRepo {
	repo := &AdminSystemRepo{
		DB: params.DB,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("AdminSystemRepo initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("AdminSystemRepo stopped")
			return nil
		},
	})

	return repo
}

// GetSystemSettings 获取系统设置
func (r *AdminSystemRepo) GetSystemSettings(ctx context.Context) (map[string]any, error) {
	// 这里可以从数据库中获取系统设置
	// 目前我们返回一个空的map，因为AdminSystemService会使用配置文件中的设置
	return map[string]any{}, nil
}
