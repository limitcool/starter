package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/repository"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo      *repository.RoleRepo
	casbinService casbin.Service
}

// NewRoleService 创建角色服务
func NewRoleService(roleRepo *repository.RoleRepo, casbinService casbin.Service) *RoleService {
	return &RoleService{
		roleRepo:      roleRepo,
		casbinService: casbinService,
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, role *model.Role) error {
	return s.roleRepo.Create(ctx, role)
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, role *model.Role) error {
	return s.roleRepo.Update(ctx, role)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, id uint) error {
	// 检查角色是否已分配给用户
	isAssigned, err := s.roleRepo.IsAssignedToUser(ctx, id)
	if err != nil {
		return err
	}
	if isAssigned {
		return errors.New("该角色已分配给用户，不能删除")
	}

	// 删除角色菜单关联
	if err := s.roleRepo.DeleteRoleMenus(ctx, id); err != nil {
		return err
	}

	// 查询角色信息
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 删除Casbin中的角色策略
	_, err = s.casbinService.DeleteRole(ctx, role.Code)
	if err != nil {
		return err
	}

	// 删除角色
	return s.roleRepo.Delete(ctx, id)
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}

// GetRoles 获取角色列表
func (s *RoleService) GetRoles(ctx context.Context) ([]model.Role, error) {
	return s.roleRepo.GetAll(ctx)
}

// AssignRolesToUser 为用户分配角色
func (s *RoleService) AssignRolesToUser(ctx context.Context, userID int64, roleIDs []uint) error {
	// 使用 roleRepo 的 AssignRolesToUser 方法
	if err := s.roleRepo.AssignRolesToUser(ctx, userID, roleIDs); err != nil {
		return err
	}

	// 更新Casbin中的用户角色
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 获取用户当前角色
	roles, err := s.casbinService.GetRolesForUser(ctx, userIDStr)
	if err != nil {
		return err
	}

	// 移除所有角色
	for _, role := range roles {
		_, err = s.casbinService.DeleteRoleForUser(ctx, userIDStr, role)
		if err != nil {
			return err
		}
	}

	// 添加新角色
	if len(roleIDs) > 0 {
		for _, roleID := range roleIDs {
			// 查询角色编码
			roleObj, err := s.roleRepo.GetByID(ctx, roleID)
			if err != nil {
				return err
			}

			// 添加用户角色关联
			_, err = s.casbinService.AddRoleForUser(ctx, userIDStr, roleObj.Code)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUserRoleIDs 获取用户角色ID列表
func (s *RoleService) GetUserRoleIDs(ctx context.Context, userID uint) ([]uint, error) {
	return s.roleRepo.GetRoleIDsByUserID(ctx, userID)
}

// GetRoleMenuIDs 获取角色菜单ID列表
func (s *RoleService) GetRoleMenuIDs(ctx context.Context, roleID uint) ([]uint, error) {
	return s.roleRepo.GetMenuIDsByRoleID(ctx, roleID)
}

// 为角色设置权限策略
func (s *RoleService) SetRolePermission(ctx context.Context, roleCode string, obj string, act string) error {
	_, err := s.casbinService.AddPermissionForRole(ctx, roleCode, obj, act)
	return err
}

// 删除角色的权限策略
func (s *RoleService) DeleteRolePermission(ctx context.Context, roleCode string, obj string, act string) error {
	_, err := s.casbinService.DeletePermissionForRole(ctx, roleCode, obj, act)
	return err
}
