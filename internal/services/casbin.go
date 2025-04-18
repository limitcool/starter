package services

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/storage/database"
)

// 注意: 我们已经移除了全局实例和单例模式

// CasbinService Casbin权限服务
type CasbinService struct {
	db       database.DB
	enforcer *casbin.Enforcer
}

// NewCasbinService 创建Casbin服务
func NewCasbinService(db database.DB) *CasbinService {
	s := &CasbinService{
		db: db,
	}

	// 初始化
	_ = s.Initialize()

	return s
}

// GetEnforcer 获取Casbin执行器
func (s *CasbinService) GetEnforcer() *casbin.Enforcer {
	return s.enforcer
}

// Initialize 初始Casbin服务
func (s *CasbinService) Initialize() error {
	if s.db == nil || s.db.GetDB() == nil {
		log.Error("数据库连接未初始化")
		return nil
	}

	// 如果权限系统未启用，直接返回
	config := core.Instance().Config()
	if config != nil && !config.Casbin.Enabled {
		return nil
	}

	// 使用gorm适配器
	adapter, err := gormadapter.NewAdapterByDB(s.db.GetDB())
	if err != nil {
		log.Error("创建Casbin适配器失败", "error", err)
		return err
	}

	// 获取模型文件路径
	modelPath := "configs/rbac_model.conf"
	if config != nil && config.Casbin.ModelPath != "" {
		modelPath = config.Casbin.ModelPath
	}

	// 创建enforcer
	e, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		log.Error("创建Casbin Enforcer失败", "error", err)
		return err
	}

	// 加载策略
	if err := e.LoadPolicy(); err != nil {
		log.Error("加载Casbin策略失败", "error", err)
		return err
	}

	// 启用自动保存
	e.EnableAutoSave(true)

	s.enforcer = e

	return nil
}

// InitCasbin 初始化Casbin
// 已弃用: 请使用依赖注入创建Casbin服务
func InitCasbin() (*casbin.Enforcer, error) {
	return nil, fmt.Errorf("请使用依赖注入创建Casbin服务")
}

// 注意: 这个方法已经在上面声明过了

// CheckPermission 检查权限
func (s *CasbinService) CheckPermission(userID string, obj string, act string) (bool, error) {
	return s.enforcer.Enforce(userID, obj, act)
}

// AddRoleForUser 为用户添加角色
func (s *CasbinService) AddRoleForUser(userID string, role string) (bool, error) {
	return s.enforcer.AddGroupingPolicy(userID, role)
}

// DeleteRoleForUser 删除用户角色
func (s *CasbinService) DeleteRoleForUser(userID string, role string) (bool, error) {
	return s.enforcer.RemoveGroupingPolicy(userID, role)
}

// AddPermissionForRole 为角色添加权限
func (s *CasbinService) AddPermissionForRole(role string, obj string, act string) (bool, error) {
	return s.enforcer.AddPolicy(role, obj, act)
}

// DeletePermissionForRole 删除角色权限
func (s *CasbinService) DeletePermissionForRole(role string, obj string, act string) (bool, error) {
	return s.enforcer.RemovePolicy(role, obj, act)
}

// GetRolesForUser 获取用户角色列表
func (s *CasbinService) GetRolesForUser(userID string) ([]string, error) {
	return s.enforcer.GetRolesForUser(userID)
}

// HasRoleForUser 判断用户是否拥有指定角色
func (s *CasbinService) HasRoleForUser(userID string, role string) (bool, error) {
	return s.enforcer.HasRoleForUser(userID, role)
}

// GetUsersForRole 获取角色的所有用户
func (s *CasbinService) GetUsersForRole(role string) ([]string, error) {
	return s.enforcer.GetUsersForRole(role)
}

// GetPermissionsForUser 获取用户所有权限
func (s *CasbinService) GetPermissionsForUser(userID string) [][]string {
	permissions, _ := s.enforcer.GetImplicitPermissionsForUser(userID)
	return permissions
}

// GetPermissionsForRole 获取角色所有权限
func (s *CasbinService) GetPermissionsForRole(role string) [][]string {
	permissions, _ := s.enforcer.GetFilteredPolicy(0, role)
	return permissions
}

// GetAllRoles 获取所有角色
func (s *CasbinService) GetAllRoles() []string {
	roles, _ := s.enforcer.GetAllRoles()
	return roles
}

// DeleteRole 删除角色
func (s *CasbinService) DeleteRole(role string) (bool, error) {
	// 删除角色所有权限
	_, err := s.enforcer.RemoveFilteredPolicy(0, role)
	if err != nil {
		return false, err
	}

	// 删除所有用户与该角色的关联
	_, err = s.enforcer.RemoveFilteredGroupingPolicy(1, role)
	if err != nil {
		return false, err
	}

	return true, nil
}
