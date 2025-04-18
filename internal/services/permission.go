package services

import (
	"path/filepath"

	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/repository"
	"github.com/spf13/viper"
)

// PermissionService 权限服务
type PermissionService struct {
	permissionRepo *repository.PermissionRepo
	casbinService  *CasbinService
}

// NewPermissionService 创建权限服务
func NewPermissionService(permissionRepo *repository.PermissionRepo) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
	}
}

// UpdatePermissionSettings 更新权限系统设置
func (s *PermissionService) UpdatePermissionSettings(enabled, defaultAllow bool) error {
	// 获取配置
	config := core.Instance().Config()

	// 更新内存中的配置
	config.Casbin.Enabled = enabled
	config.Casbin.DefaultAllow = defaultAllow

	// 更新配置文件
	v := viper.New()
	v.SetConfigFile(filepath.Join("configs", "config.yaml"))

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	v.Set("casbin.enabled", enabled)
	v.Set("casbin.default_allow", defaultAllow)

	return v.WriteConfig()
}

// GetPermissions 获取权限列表
func (s *PermissionService) GetPermissions() ([]model.Permission, error) {
	return s.permissionRepo.GetAll()
}

// GetPermission 获取权限详情
func (s *PermissionService) GetPermission(id uint64) (*model.Permission, error) {
	return s.permissionRepo.GetByID(uint(id))
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(permission *model.Permission) error {
	return s.permissionRepo.Create(permission)
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(id uint64, permission *model.Permission) error {
	permission.ID = uint(id)
	return s.permissionRepo.Update(permission)
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(id uint64) error {
	return s.permissionRepo.Delete(uint(id))
}

// AssignPermissionToRole 为角色分配权限
func (s *PermissionService) AssignPermissionToRole(roleID uint, permissionIDs []uint) error {
	return s.permissionRepo.AssignPermissionToRole(roleID, permissionIDs)
}

// GetPermissionsByRoleID 获取角色的权限列表
func (s *PermissionService) GetPermissionsByRoleID(roleID uint) ([]model.Permission, error) {
	return s.permissionRepo.GetByRoleID(roleID)
}

// GetPermissionsByUserID 获取用户的权限列表
func (s *PermissionService) GetPermissionsByUserID(userID uint) ([]model.Permission, error) {
	return s.permissionRepo.GetByUserID(userID)
}
