package services

import (
	"path/filepath"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/repository"
	"github.com/spf13/viper"
)

// PermissionService 权限服务
type PermissionService struct {
	permissionRepo *repository.PermissionRepo
	roleRepo       *repository.RoleRepo
	menuRepo       *repository.MenuRepo
	casbinService  casbin.Service
	config         *configs.Config
}

// NewPermissionService 创建权限服务
func NewPermissionService(
	permissionRepo *repository.PermissionRepo,
	roleRepo *repository.RoleRepo,
	menuRepo *repository.MenuRepo,
	casbinService casbin.Service,
	config *configs.Config,
) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
		menuRepo:       menuRepo,
		casbinService:  casbinService,
		config:         config,
	}
}

// UpdatePermissionSettings 更新权限系统设置
func (s *PermissionService) UpdatePermissionSettings(enabled, defaultAllow bool) error {
	// 更新内存中的配置
	s.config.Casbin.Enabled = enabled
	s.config.Casbin.DefaultAllow = defaultAllow

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
	// 数据库层面的关联
	err := s.permissionRepo.AssignPermissionToRole(roleID, permissionIDs)
	if err != nil {
		return err
	}

	// 获取角色
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	// 删除Casbin中的角色权限
	_, err = s.casbinService.DeleteRole(role.Code)
	if err != nil {
		return err
	}

	// 添加新的权限
	for _, permID := range permissionIDs {
		perm, err := s.permissionRepo.GetByID(permID)
		if err != nil {
			continue
		}

		// 添加到Casbin
		_, err = s.casbinService.AddPermissionForRole(role.Code, perm.Code, "*")
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPermissionsByRoleID 获取角色的权限列表
func (s *PermissionService) GetPermissionsByRoleID(roleID uint) ([]model.Permission, error) {
	return s.permissionRepo.GetByRoleID(roleID)
}

// GetPermissionsByUserID 获取用户的权限列表
func (s *PermissionService) GetPermissionsByUserID(userID uint) ([]model.Permission, error) {
	return s.permissionRepo.GetByUserID(userID)
}

// CheckPermission 检查权限
func (s *PermissionService) CheckPermission(userID string, obj string, act string) (bool, error) {
	return s.casbinService.CheckPermission(userID, obj, act)
}

// AssignRolesToUser 为用户分配角色
func (s *PermissionService) AssignRolesToUser(userID string, roleIDs []uint) error {
	// 获取用户当前角色
	currentRoles, err := s.casbinService.GetRolesForUser(userID)
	if err != nil {
		return err
	}

	// 删除当前角色
	for _, role := range currentRoles {
		_, err := s.casbinService.DeleteRoleForUser(userID, role)
		if err != nil {
			return err
		}
	}

	// 添加新角色
	for _, roleID := range roleIDs {
		role, err := s.roleRepo.GetByID(roleID)
		if err != nil {
			continue
		}

		_, err = s.casbinService.AddRoleForUser(userID, role.Code)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetUserRoles 获取用户角色
func (s *PermissionService) GetUserRoles(userID string) ([]model.Role, error) {
	// 获取用户角色编码
	roleCodes, err := s.casbinService.GetRolesForUser(userID)
	if err != nil {
		return nil, err
	}

	// 获取角色详情
	var roles []model.Role
	for _, code := range roleCodes {
		role, err := s.roleRepo.GetByCode(code)
		if err != nil {
			continue
		}
		roles = append(roles, *role)
	}

	return roles, nil
}

// GetUserMenus 获取用户菜单
func (s *PermissionService) GetUserMenus(userID string) ([]*model.MenuTree, error) {
	// 获取用户角色
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	// 获取角色菜单
	var menuIDs []uint
	for _, role := range roles {
		// 管理员角色获取所有菜单
		if role.Code == "admin" {
			allMenus, err := s.menuRepo.GetAll()
			if err != nil {
				return nil, err
			}
			return s.menuRepo.BuildMenuTree(allMenus), nil
		}

		// 获取角色菜单
		roleMenus, err := s.menuRepo.GetByRoleID(role.ID)
		if err != nil {
			continue
		}

		for _, menu := range roleMenus {
			menuIDs = append(menuIDs, menu.ID)
		}
	}

	// 去重
	uniqueMenuIDs := make(map[uint]bool)
	var uniqueIDs []uint
	for _, id := range menuIDs {
		if !uniqueMenuIDs[id] {
			uniqueMenuIDs[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	// 获取菜单详情
	var menus []*model.Menu
	for _, id := range uniqueIDs {
		menu, err := s.menuRepo.GetByIDWithRelations(id)
		if err != nil {
			continue
		}
		menus = append(menus, menu)
	}

	// 构建菜单树
	return s.menuRepo.BuildMenuTree(menus), nil
}
