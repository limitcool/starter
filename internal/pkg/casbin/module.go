package casbin

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/enum"
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

	// 获取用户模式
	userMode := enum.GetUserMode(cfg.Admin.UserMode)

	// 如果是简单模式或者Casbin未启用，则不初始化Casbin
	if userMode == enum.UserModeSimple || !cfg.Casbin.Enabled {
		logger.InfoContext(ctx, "Casbin disabled", "user_mode", userMode)
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

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load Casbin policy: %w", err)
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

	// 获取用户模式
	userMode := enum.GetUserMode(cfg.Admin.UserMode)

	// 如果是简单模式或者Casbin未启用，则不初始化Casbin
	if userMode == enum.UserModeSimple || !cfg.Casbin.Enabled {
		logger.InfoContext(ctx, "Casbin service disabled", "user_mode", userMode)
		return nil, nil
	}

	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// 创建Casbin服务
	return NewService(db, cfg), nil
}
