package casbin

import (
	"context"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module Casbin模块
var Module = fx.Options(
	// 提供Casbin执行器
	fx.Provide(NewEnforcer),
	// 提供Casbin服务
	fx.Provide(NewCasbinService),
)

// NewEnforcer 创建Casbin执行器
func NewEnforcer(lc fx.Lifecycle, cfg *configs.Config, db *gorm.DB) (*casbin.Enforcer, error) {
	// 使用默认上下文
	ctx := context.Background()

	if !cfg.Casbin.Enabled {
		logger.InfoContext(ctx, "Casbin disabled")
		return nil, nil
	}

	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	logger.InfoContext(ctx, "Initializing Casbin")

	// 创建适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin adapter: %w", err)
	}

	// 加载模型
	var m model.Model
	if cfg.Casbin.ModelPath != "" {
		m, err = model.NewModelFromFile(cfg.Casbin.ModelPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load Casbin model: %w", err)
		}
	} else {
		// 使用默认模型
		m, err = model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
		if err != nil {
			return nil, fmt.Errorf("failed to create default Casbin model: %w", err)
		}
	}

	// 创建执行器
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin enforcer: %w", err)
	}

	// 启用自8动加载
	if cfg.Casbin.AutoLoad {
		// 注意：Casbin v2 不支持 StartAutoLoadPolicy 方法
		// 这里可以使用定时器或其他方式实现自动加载
		go func() {
			for {
				time.Sleep(time.Duration(cfg.Casbin.AutoLoadInterval) * time.Second)
				err := enforcer.LoadPolicy()
				if err != nil {
					logger.ErrorContext(context.Background(), "Failed to auto load policy", "error", err)
				}
			}
		}()
	}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "Casbin initialized successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "Cleaning up Casbin")
			return nil
		},
	})

	return enforcer, nil
}

// NewCasbinService 创建Casbin服务
func NewCasbinService(db *gorm.DB, cfg *configs.Config) (Service, error) {
	// 使用默认上下文
	ctx := context.Background()

	if !cfg.Casbin.Enabled {
		logger.InfoContext(ctx, "Casbin disabled")
		return nil, nil
	}

	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// 创建Casbin服务
	return NewService(db, cfg), nil
}
