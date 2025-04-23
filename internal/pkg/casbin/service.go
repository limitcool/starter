package casbin

import (
	"context"
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service Casbin服务接口
type Service interface {
	// 初始化
	Initialize() error

	// 获取执行器
	GetEnforcer() *casbin.Enforcer

	// 权限检查
	CheckPermission(ctx context.Context, sub, obj, act string) (bool, error)

	// 角色管理
	AddRoleForUser(ctx context.Context, userID, role string) (bool, error)
	DeleteRoleForUser(ctx context.Context, userID, role string) (bool, error)
	GetRolesForUser(ctx context.Context, userID string) ([]string, error)
	HasRoleForUser(ctx context.Context, userID, role string) (bool, error)
	GetUsersForRole(ctx context.Context, role string) ([]string, error)
	GetAllRoles(ctx context.Context) ([]string, error)
	DeleteRole(ctx context.Context, role string) (bool, error)

	// 批量角色管理
	AddRolesForUser(ctx context.Context, userID string, roles []string) (bool, error)
	DeleteRolesForUser(ctx context.Context, userID string) (bool, error)

	// 权限管理
	AddPermissionForRole(ctx context.Context, role, obj, act string) (bool, error)
	DeletePermissionForRole(ctx context.Context, role, obj, act string) (bool, error)
	GetPermissionsForUser(ctx context.Context, userID string) ([][]string, error)
	GetPermissionsForRole(ctx context.Context, role string) ([][]string, error)

	// 批量权限管理
	AddPolicies(ctx context.Context, policies [][]string) (bool, error)
	RemovePolicies(ctx context.Context, policies [][]string) (bool, error)
}

// DefaultService Casbin服务默认实现
type DefaultService struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
	config   *configs.Config
	mutex    sync.RWMutex // 保护enforcer的并发访问
}

// NewService 创建Casbin服务
func NewService(db *gorm.DB, config *configs.Config) Service {
	s := &DefaultService{
		db:     db,
		config: config,
	}

	// 初始化
	_ = s.Initialize()

	return s
}

// Initialize 初始化Casbin服务
func (s *DefaultService) Initialize() error {
	return s.InitializeWithContext(context.Background())
}

// InitializeWithContext 使用上下文初始化Casbin服务
func (s *DefaultService) InitializeWithContext(ctx context.Context) error {
	if s.db == nil {
		logger.ErrorContext(ctx, "数据库连接未初始化")
		return errorx.Errorf(errorx.ErrInternal, "数据库连接未初始化")
	}

	// 如果权限系统未启用，直接返回
	if s.config != nil && !s.config.Casbin.Enabled {
		logger.InfoContext(ctx, "Casbin权限系统未启用")
		return nil
	}

	logger.InfoContext(ctx, "初始化Casbin服务")

	// 使用gorm适配器
	adapter, err := gormadapter.NewAdapterByDB(s.db)
	if err != nil {
		logger.ErrorContext(ctx, "创建Casbin适配器失败", "error", err)
		return errorx.WrapError(err, "创建Casbin适配器失败")
	}

	// 获取模型文件路径
	modelPath := "configs/rbac_model.conf"
	if s.config != nil && s.config.Casbin.ModelPath != "" {
		modelPath = s.config.Casbin.ModelPath
	}

	logger.DebugContext(ctx, "加载Casbin模型", "path", modelPath)

	// 创建enforcer
	e, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		logger.ErrorContext(ctx, "创建Casbin Enforcer失败", "error", err)
		return errorx.WrapError(err, "创建Casbin Enforcer失败")
	}

	// 加载策略
	logger.DebugContext(ctx, "加载Casbin策略")
	if err := e.LoadPolicy(); err != nil {
		logger.ErrorContext(ctx, "加载Casbin策略失败", "error", err)
		return errorx.WrapError(err, "加载Casbin策略失败")
	}

	// 启用自动保存
	e.EnableAutoSave(true)

	s.mutex.Lock()
	s.enforcer = e
	s.mutex.Unlock()

	logger.InfoContext(ctx, "Casbin服务初始化成功")
	return nil
}

// GetEnforcer 获取Casbin执行器
func (s *DefaultService) GetEnforcer() *casbin.Enforcer {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.enforcer
}

// CheckPermission 检查权限
func (s *DefaultService) CheckPermission(ctx context.Context, sub, obj, act string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.enforcer == nil {
		return false, errorx.Errorf(errorx.ErrCasbinService, "Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "检查权限", "subject", sub, "object", obj, "action", act)
	result, err := s.enforcer.Enforce(sub, obj, act)
	if err != nil {
		return false, errorx.WrapError(err, "权限检查失败")
	}

	return result, nil
}

// AddRoleForUser 为用户添加角色
func (s *DefaultService) AddRoleForUser(ctx context.Context, userID string, role string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.enforcer == nil {
		return false, errorx.Errorf(errorx.ErrCasbinService, "Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "为用户添加角色", "user", userID, "role", role)
	return s.enforcer.AddGroupingPolicy(userID, role)
}

// DeleteRoleForUser 删除用户角色
func (s *DefaultService) DeleteRoleForUser(ctx context.Context, userID string, role string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "删除用户角色", "user", userID, "role", role)
	return s.enforcer.RemoveGroupingPolicy(userID, role)
}

// AddPermissionForRole 为角色添加权限
func (s *DefaultService) AddPermissionForRole(ctx context.Context, role string, obj string, act string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "为角色添加权限", "role", role, "object", obj, "action", act)
	return s.enforcer.AddPolicy(role, obj, act)
}

// DeletePermissionForRole 删除角色权限
func (s *DefaultService) DeletePermissionForRole(ctx context.Context, role string, obj string, act string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "删除角色权限", "role", role, "object", obj, "action", act)
	return s.enforcer.RemovePolicy(role, obj, act)
}

// GetRolesForUser 获取用户角色列表
func (s *DefaultService) GetRolesForUser(ctx context.Context, userID string) ([]string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.enforcer == nil {
		return nil, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "获取用户角色", "user", userID)
	return s.enforcer.GetRolesForUser(userID)
}

// HasRoleForUser 判断用户是否拥有指定角色
func (s *DefaultService) HasRoleForUser(ctx context.Context, userID string, role string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "检查用户角色", "user", userID, "role", role)
	return s.enforcer.HasRoleForUser(userID, role)
}

// GetUsersForRole 获取角色的所有用户
func (s *DefaultService) GetUsersForRole(ctx context.Context, role string) ([]string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.enforcer == nil {
		return nil, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "获取角色用户", "role", role)
	return s.enforcer.GetUsersForRole(role)
}

// GetPermissionsForUser 获取用户所有权限
func (s *DefaultService) GetPermissionsForUser(ctx context.Context, userID string) ([][]string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.enforcer == nil {
		return nil, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "获取用户权限", "user", userID)
	permissions, err := s.enforcer.GetImplicitPermissionsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户权限失败: %w", err)
	}
	return permissions, nil
}

// GetPermissionsForRole 获取角色所有权限
func (s *DefaultService) GetPermissionsForRole(ctx context.Context, role string) ([][]string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.enforcer == nil {
		return nil, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "获取角色权限", "role", role)
	permissions, err := s.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		return nil, fmt.Errorf("获取角色权限失败: %w", err)
	}
	return permissions, nil
}

// GetAllRoles 获取所有角色
func (s *DefaultService) GetAllRoles(ctx context.Context) ([]string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.enforcer == nil {
		return nil, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "获取所有角色")
	roles, err := s.enforcer.GetAllRoles()
	if err != nil {
		return nil, fmt.Errorf("获取所有角色失败: %w", err)
	}
	return roles, nil
}

// DeleteRole 删除角色
func (s *DefaultService) DeleteRole(ctx context.Context, role string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "删除角色", "role", role)

	// 删除角色所有权限
	_, err := s.enforcer.RemoveFilteredPolicy(0, role)
	if err != nil {
		return false, fmt.Errorf("删除角色权限失败: %w", err)
	}

	// 删除所有用户与该角色的关联
	_, err = s.enforcer.RemoveFilteredGroupingPolicy(1, role)
	if err != nil {
		return false, fmt.Errorf("删除用户与角色关联失败: %w", err)
	}

	logger.InfoContext(ctx, "角色删除成功", "role", role)
	return true, nil
}

// AddRolesForUser 为用户批量添加角色
func (s *DefaultService) AddRolesForUser(ctx context.Context, userID string, roles []string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "为用户批量添加角色", "user", userID, "roles", roles)

	// 准备批量添加的角色策略
	var rules [][]string
	for _, role := range roles {
		rules = append(rules, []string{userID, role})
	}

	// 批量添加角色
	success, err := s.enforcer.AddGroupingPolicies(rules)
	if err != nil {
		return false, fmt.Errorf("批量添加用户角色失败: %w", err)
	}

	return success, nil
}

// DeleteRolesForUser 删除用户所有角色
func (s *DefaultService) DeleteRolesForUser(ctx context.Context, userID string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "删除用户所有角色", "user", userID)

	// 删除用户所有角色
	success, err := s.enforcer.DeleteRolesForUser(userID)
	if err != nil {
		return false, fmt.Errorf("删除用户角色失败: %w", err)
	}

	return success, nil
}

// AddPolicies 批量添加权限策略
func (s *DefaultService) AddPolicies(ctx context.Context, policies [][]string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "批量添加权限策略", "count", len(policies))

	// 批量添加策略
	success, err := s.enforcer.AddPolicies(policies)
	if err != nil {
		return false, fmt.Errorf("批量添加权限策略失败: %w", err)
	}

	return success, nil
}

// RemovePolicies 批量删除权限策略
func (s *DefaultService) RemovePolicies(ctx context.Context, policies [][]string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer 未初始化")
	}

	logger.DebugContext(ctx, "批量删除权限策略", "count", len(policies))

	// 批量删除策略
	success, err := s.enforcer.RemovePolicies(policies)
	if err != nil {
		return false, fmt.Errorf("批量删除权限策略失败: %w", err)
	}

	return success, nil
}
