package casbin

import (
	"context"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service Casbin权限服务
type Service struct {
	enforcer *casbin.Enforcer
	config   *configs.Config
}

// NewService 创建Casbin服务
func NewService(db *gorm.DB, config *configs.Config) (*Service, error) {
	if !config.Casbin.Enabled {
		return nil, nil
	}

	// 创建GORM适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, errorx.WrapError(err, "创建Casbin适配器失败")
	}

	// 创建执行器
	enforcer, err := casbin.NewEnforcer(config.Casbin.ModelPath, adapter)
	if err != nil {
		return nil, errorx.WrapError(err, "创建Casbin执行器失败")
	}

	// 启用日志
	if config.Casbin.EnableLog {
		enforcer.EnableLog(true)
	}

	// 启用自动保存
	if config.Casbin.EnableAutoSave {
		enforcer.EnableAutoSave(true)
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, errorx.WrapError(err, "加载Casbin策略失败")
	}

	service := &Service{
		enforcer: enforcer,
		config:   config,
	}

	// 启动自动加载策略
	if config.Casbin.AutoLoadInterval > 0 {
		go service.autoLoadPolicy()
	}

	logger.InfoContext(context.Background(), "Casbin权限服务初始化成功",
		"model_path", config.Casbin.ModelPath,
		"auto_load_interval", config.Casbin.AutoLoadInterval)

	return service, nil
}

// autoLoadPolicy 自动加载策略
func (s *Service) autoLoadPolicy() {
	ticker := time.NewTicker(time.Duration(s.config.Casbin.AutoLoadInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := s.enforcer.LoadPolicy()
		if err != nil {
			logger.ErrorContext(context.Background(), "自动加载Casbin策略失败", "error", err)
		}
	}
}

// Enforce 权限验证
func (s *Service) Enforce(ctx context.Context, sub, obj, act string) (bool, error) {
	if s.enforcer == nil {
		return false, errorx.WrapError(nil, "Casbin服务未初始化")
	}

	result, err := s.enforcer.Enforce(sub, obj, act)
	if err != nil {
		logger.ErrorContext(ctx, "权限验证失败",
			"subject", sub,
			"object", obj,
			"action", act,
			"error", err)
		return false, errorx.WrapError(err, "权限验证失败")
	}

	logger.DebugContext(ctx, "权限验证结果",
		"subject", sub,
		"object", obj,
		"action", act,
		"result", result)

	return result, nil
}

// AddPolicy 添加策略
func (s *Service) AddPolicy(ctx context.Context, sub, obj, act string) error {
	if s.enforcer == nil {
		return errorx.WrapError(nil, "Casbin服务未初始化")
	}

	added, err := s.enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		return errorx.WrapError(err, "添加策略失败")
	}

	if !added {
		logger.WarnContext(ctx, "策略已存在", "subject", sub, "object", obj, "action", act)
	}

	return nil
}

// RemovePolicy 删除策略
func (s *Service) RemovePolicy(ctx context.Context, sub, obj, act string) error {
	if s.enforcer == nil {
		return errorx.WrapError(nil, "Casbin服务未初始化")
	}

	removed, err := s.enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		return errorx.WrapError(err, "删除策略失败")
	}

	if !removed {
		logger.WarnContext(ctx, "策略不存在", "subject", sub, "object", obj, "action", act)
	}

	return nil
}

// AddRoleForUser 为用户添加角色
func (s *Service) AddRoleForUser(ctx context.Context, user, role string) error {
	if s.enforcer == nil {
		return errorx.WrapError(nil, "Casbin服务未初始化")
	}

	added, err := s.enforcer.AddRoleForUser(user, role)
	if err != nil {
		return errorx.WrapError(err, "为用户添加角色失败")
	}

	if !added {
		logger.WarnContext(ctx, "用户角色已存在", "user", user, "role", role)
	}

	return nil
}

// DeleteRoleForUser 删除用户角色
func (s *Service) DeleteRoleForUser(ctx context.Context, user, role string) error {
	if s.enforcer == nil {
		return errorx.WrapError(nil, "Casbin服务未初始化")
	}

	deleted, err := s.enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		return errorx.WrapError(err, "删除用户角色失败")
	}

	if !deleted {
		logger.WarnContext(ctx, "用户角色不存在", "user", user, "role", role)
	}

	return nil
}

// DeleteRolesForUser 删除用户的所有角色
func (s *Service) DeleteRolesForUser(ctx context.Context, user string) error {
	if s.enforcer == nil {
		return errorx.WrapError(nil, "Casbin服务未初始化")
	}

	deleted, err := s.enforcer.DeleteRolesForUser(user)
	if err != nil {
		return errorx.WrapError(err, "删除用户所有角色失败")
	}

	if !deleted {
		logger.WarnContext(ctx, "用户没有角色", "user", user)
	}

	return nil
}

// GetRolesForUser 获取用户的角色
func (s *Service) GetRolesForUser(ctx context.Context, user string) ([]string, error) {
	if s.enforcer == nil {
		return nil, errorx.WrapError(nil, "Casbin服务未初始化")
	}

	roles, err := s.enforcer.GetRolesForUser(user)
	if err != nil {
		return nil, errorx.WrapError(err, "获取用户角色失败")
	}

	return roles, nil
}

// GetUsersForRole 获取角色的用户
func (s *Service) GetUsersForRole(ctx context.Context, role string) ([]string, error) {
	if s.enforcer == nil {
		return nil, errorx.WrapError(nil, "Casbin服务未初始化")
	}

	users, err := s.enforcer.GetUsersForRole(role)
	if err != nil {
		return nil, errorx.WrapError(err, "获取角色用户失败")
	}

	return users, nil
}

// HasRoleForUser 检查用户是否有角色
func (s *Service) HasRoleForUser(ctx context.Context, user, role string) (bool, error) {
	if s.enforcer == nil {
		return false, errorx.WrapError(nil, "Casbin服务未初始化")
	}

	hasRole, err := s.enforcer.HasRoleForUser(user, role)
	if err != nil {
		return false, errorx.WrapError(err, "检查用户角色失败")
	}

	return hasRole, nil
}

// GetImplicitPermissionsForUser 获取用户的所有权限（包括继承的）
func (s *Service) GetImplicitPermissionsForUser(ctx context.Context, user string) ([][]string, error) {
	if s.enforcer == nil {
		return nil, errorx.WrapError(nil, "Casbin服务未初始化")
	}

	permissions, err := s.enforcer.GetImplicitPermissionsForUser(user)
	if err != nil {
		return nil, errorx.WrapError(err, "获取用户权限失败")
	}
	return permissions, nil
}

// GetUserKey 获取用户的Casbin标识
func GetUserKey(userID int64) string {
	return fmt.Sprintf("user:%d", userID)
}

// GetRoleKey 获取角色的Casbin标识
func GetRoleKey(roleID uint) string {
	return fmt.Sprintf("role:%d", roleID)
}

// GetResourceKey 获取资源的Casbin标识
func GetResourceKey(resource string) string {
	return resource
}
