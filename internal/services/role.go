package services

import (
	"errors"
	"strconv"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/repository"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo      repository.RoleRepository
	casbinService *CasbinService
}

// NewRoleService 创建角色服务
func NewRoleService(roleRepo repository.RoleRepository, casbinService *CasbinService) *RoleService {
	return &RoleService{
		roleRepo:      roleRepo,
		casbinService: casbinService,
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(role *model.Role) error {
	return s.roleRepo.Create(role)
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(role *model.Role) error {
	return s.roleRepo.Update(role)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	role := &model.Role{}

	// 检查角色是否已分配给用户
	isAssigned, err := role.IsAssignedToUser(id)
	if err != nil {
		return err
	}
	if isAssigned {
		return errors.New("该角色已分配给用户，不能删除")
	}

	// 删除角色菜单关联
	if err := role.DeleteRoleMenus(id); err != nil {
		return err
	}

	// 查询角色信息
	role, err = role.GetByID(id)
	if err != nil {
		return err
	}

	// 删除Casbin中的角色策略
	_, err = s.casbinService.DeleteRole(role.Code)
	if err != nil {
		return err
	}

	// 删除角色
	return role.Delete()
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	role := &model.Role{}
	return role.GetByID(id)
}

// GetRoles 获取角色列表
func (s *RoleService) GetRoles() ([]model.Role, error) {
	role := &model.Role{}
	return role.GetAll()
}

// AssignRolesToUser 为用户分配角色
func (s *RoleService) AssignRolesToUser(userID int64, roleIDs []uint) error {
	userRole := &model.UserRole{}

	// 删除原有的用户角色关联
	if err := userRole.DeleteByUserID(userID); err != nil {
		return err
	}

	// 添加新的用户角色关联
	if len(roleIDs) > 0 {
		var userRoles []model.UserRole
		for _, roleID := range roleIDs {
			userRoles = append(userRoles, model.UserRole{
				UserID: userID,
				RoleID: roleID,
			})
		}
		if err := userRole.BatchCreate(userRoles); err != nil {
			return err
		}
	}

	// 更新Casbin中的用户角色
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 获取用户当前角色
	roles, err := s.casbinService.GetRolesForUser(userIDStr)
	if err != nil {
		return err
	}

	// 移除所有角色
	for _, role := range roles {
		_, err = s.casbinService.DeleteRoleForUser(userIDStr, role)
		if err != nil {
			return err
		}
	}

	// 添加新角色
	if len(roleIDs) > 0 {
		role := &model.Role{}
		for _, roleID := range roleIDs {
			// 查询角色编码
			roleObj, err := role.GetByID(roleID)
			if err != nil {
				return err
			}

			// 添加用户角色关联
			_, err = s.casbinService.AddRoleForUser(userIDStr, roleObj.Code)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUserRoleIDs 获取用户角色ID列表
func (s *RoleService) GetUserRoleIDs(userID uint) ([]uint, error) {
	userRole := &model.UserRole{}
	return userRole.GetRoleIDsByUserID(userID)
}

// GetRoleMenuIDs 获取角色菜单ID列表
func (s *RoleService) GetRoleMenuIDs(roleID uint) ([]uint, error) {
	roleMenu := &model.RoleMenu{}
	return roleMenu.GetMenuIDsByRoleID(roleID)
}

// 为角色设置权限策略
func (s *RoleService) SetRolePermission(roleCode string, obj string, act string) error {
	_, err := s.casbinService.AddPermissionForRole(roleCode, obj, act)
	return err
}

// 删除角色的权限策略
func (s *RoleService) DeleteRolePermission(roleCode string, obj string, act string) error {
	_, err := s.casbinService.DeletePermissionForRole(roleCode, obj, act)
	return err
}
