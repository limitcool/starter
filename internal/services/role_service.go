package services

import (
	"errors"
	"strconv"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/repository"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo      *repository.RoleRepo
	casbinService *CasbinService
}

// NewRoleService 创建角色服务
func NewRoleService(roleRepo *repository.RoleRepo, casbinService *CasbinService) *RoleService {
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
	// 检查角色是否已分配给用户
	isAssigned, err := s.roleRepo.IsAssignedToUser(id)
	if err != nil {
		return err
	}
	if isAssigned {
		return errors.New("该角色已分配给用户，不能删除")
	}

	// 删除角色菜单关联
	if err := s.roleRepo.DeleteRoleMenus(id); err != nil {
		return err
	}

	// 查询角色信息
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除Casbin中的角色策略
	_, err = s.casbinService.DeleteRole(role.Code)
	if err != nil {
		return err
	}

	// 删除角色
	return s.roleRepo.Delete(id)
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	return s.roleRepo.GetByID(id)
}

// GetRoles 获取角色列表
func (s *RoleService) GetRoles() ([]model.Role, error) {
	return s.roleRepo.GetAll()
}

// AssignRolesToUser 为用户分配角色
func (s *RoleService) AssignRolesToUser(userID int64, roleIDs []uint) error {
	// 使用 roleRepo 的 AssignRolesToUser 方法
	if err := s.roleRepo.AssignRolesToUser(userID, roleIDs); err != nil {
		return err
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
		for _, roleID := range roleIDs {
			// 查询角色编码
			roleObj, err := s.roleRepo.GetByID(roleID)
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
	return s.roleRepo.GetRoleIDsByUserID(userID)
}

// GetRoleMenuIDs 获取角色菜单ID列表
func (s *RoleService) GetRoleMenuIDs(roleID uint) ([]uint, error) {
	return s.roleRepo.GetMenuIDsByRoleID(roleID)
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
