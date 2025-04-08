package services

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"gorm.io/gorm"
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

// CasbinService Casbin权限服务
type CasbinService struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
}

// NewCasbinService 创建Casbin服务
func NewCasbinService(db *gorm.DB) *CasbinService {
	e, err := InitCasbin(db)
	if err != nil {
		log.Error("初始化Casbin失败", "error", err)
		return nil
	}
	return &CasbinService{
		db:       db,
		enforcer: e,
	}
}

// InitCasbin 初始化Casbin
func InitCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	var err error

	// 如果权限系统未启用，直接返回nil
	if global.Config != nil && !global.Config.Permission.Enabled {
		return nil, nil
	}

	// 如果已经初始化，直接返回
	if enforcer != nil {
		return enforcer, nil
	}

	once.Do(func() {
		// 检查数据库连接是否为nil
		if db == nil {
			// 尝试从SQL组件获取数据库连接
			db = sqldb.GetDB()
			if db == nil {
				err = fmt.Errorf("数据库未初始化")
				return
			}
		}

		// 使用gorm适配器
		adapter, adapterErr := gormadapter.NewAdapterByDB(db)
		if adapterErr != nil {
			err = adapterErr
			return
		}

		// 获取模型文件路径
		modelPath := "configs/rbac_model.conf"
		if global.Config != nil && global.Config.Permission.ModelPath != "" {
			modelPath = global.Config.Permission.ModelPath
		}

		// 创建enforcer
		e, casbinErr := casbin.NewEnforcer(modelPath, adapter)
		if casbinErr != nil {
			err = casbinErr
			return
		}

		// 加载策略
		if loadErr := e.LoadPolicy(); loadErr != nil {
			err = loadErr
			return
		}

		// 启用自动保存
		e.EnableAutoSave(true)

		enforcer = e
	})

	return enforcer, err
}

// GetEnforcer 获取Casbin实例
func (s *CasbinService) GetEnforcer() *casbin.Enforcer {
	return s.enforcer
}

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
