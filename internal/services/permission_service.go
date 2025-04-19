package services

import (
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/log"
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
	// 获取角色
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	// 开始数据库事务
	tx := s.permissionRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务中执行数据库操作
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新的角色权限关联
	if len(permissionIDs) > 0 {
		var rolePermissions []model.RolePermission
		for _, permID := range permissionIDs {
			rolePermissions = append(rolePermissions, model.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
			})
		}
		if err := tx.Create(&rolePermissions).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交数据库事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 使用 Casbin 事务批量更新权限策略
	// 注意：这里我们使用 Casbin 的批量操作 API
	// 1. 删除 Casbin 中的角色权限
	_, err = s.casbinService.DeleteRole(role.Code)
	if err != nil {
		// 如果 Casbin 操作失败，记录错误但不回滚数据库事务
		// 因为数据库事务已经提交，我们可以在后续的同步操作中修复 Casbin 状态
		log.Error("删除 Casbin 角色权限失败", "error", err)
		// 这里可以添加重试逻辑或触发异步修复任务
	}

	// 2. 批量添加新的权限策略
	if len(permissionIDs) > 0 {
		// 准备批量添加的策略
		var policies [][]string
		for _, permID := range permissionIDs {
			perm, err := s.permissionRepo.GetByID(permID)
			if err != nil {
				log.Warn("获取权限信息失败", "permission_id", permID, "error", err)
				continue
			}

			// 添加策略 [role.Code, perm.Code, "*"]
			policies = append(policies, []string{role.Code, perm.Code, "*"})
		}

		// 批量添加策略
		if len(policies) > 0 {
			success, err := s.casbinService.AddPolicies(policies)
			if err != nil || !success {
				log.Error("批量添加 Casbin 权限策略失败", "error", err)
				// 这里可以添加重试逻辑或触发异步修复任务
			}
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

	// 将 userID 转换为 int64 用于数据库操作
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}

	// 开始数据库事务
	tx := s.roleRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务中删除原有的用户角色关联
	if err := tx.Where("user_id = ?", userIDInt).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新的用户角色关联
	if len(roleIDs) > 0 {
		var userRoles []model.UserRole
		for _, roleID := range roleIDs {
			userRoles = append(userRoles, model.UserRole{
				UserID: userIDInt,
				RoleID: roleID,
			})
		}
		if err := tx.Create(&userRoles).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交数据库事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 使用 Casbin 批量操作 API 更新权限
	// 1. 批量删除当前角色
	if len(currentRoles) > 0 {
		success, err := s.casbinService.DeleteRolesForUser(userID)
		if err != nil || !success {
			log.Error("批量删除用户角色失败", "error", err)
			// 这里可以添加重试逻辑或触发异步修复任务
		}
	}

	// 2. 批量添加新角色
	if len(roleIDs) > 0 {
		// 准备批量添加的角色
		var roleCodes []string
		for _, roleID := range roleIDs {
			role, err := s.roleRepo.GetByID(roleID)
			if err != nil {
				log.Warn("获取角色信息失败", "role_id", roleID, "error", err)
				continue
			}
			roleCodes = append(roleCodes, role.Code)
		}

		// 批量添加角色
		if len(roleCodes) > 0 {
			success, err := s.casbinService.AddRolesForUser(userID, roleCodes)
			if err != nil || !success {
				log.Error("批量添加用户角色失败", "error", err)
				// 这里可以添加重试逻辑或触发异步修复任务
			}
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
