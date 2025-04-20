package sqldb

import (
	"context"
	"fmt"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module 数据库模块
var Module = fx.Options(
	// 提供数据库连接
	fx.Provide(NewDB),
)

// NewDB 创建数据库连接
func NewDB(lc fx.Lifecycle, cfg *configs.Config) (*gorm.DB, error) {
	if !cfg.Database.Enabled {
		logger.Info("Database disabled")
		return nil, nil
	}

	logger.Info("Connecting to database", "driver", cfg.Driver)

	// 创建数据库连接
	db := newDbConn(cfg)

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 检查数据库连接
			sqlDB, err := db.DB()
			if err != nil {
				return fmt.Errorf("failed to get database connection: %w", err)
			}

			if err := sqlDB.Ping(); err != nil {
				return fmt.Errorf("failed to ping database: %w", err)
			}

			logger.Info("Database connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection")
			sqlDB, err := db.DB()
			if err != nil {
				return fmt.Errorf("failed to get database connection: %w", err)
			}
			return sqlDB.Close()
		},
	})

	return db, nil
}
