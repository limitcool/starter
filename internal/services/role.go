package services

import (
	"errors"
	"strconv"

	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService struct {
	db            *gorm.DB
	casbinService *CasbinService
}

// NewRoleService 创建角色服务
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		db:            db,
		casbinService: NewCasbinService(db),
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(role *model.Role) error {
	return s.db.Create(role).Error
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(role *model.Role) error {
	return s.db.Model(&model.Role{}).Where("id = ?", role.ID).Updates(role).Error
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 检查角色是否已分配给用户
		var count int64
		if err := tx.Model(&model.UserRole{}).Where("role_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("该角色已分配给用户，不能删除")
		}

		// 删除角色菜单关联
		if err := tx.Where("role_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}

		// 查询角色信息
		var role model.Role
		if err := tx.Where("id = ?", id).First(&role).Error; err != nil {
			return err
		}

		// 删除Casbin中的角色策略
		_, err := s.casbinService.DeleteRole(role.Code)
		if err != nil {
			return err
		}

		// 删除角色
		return tx.Delete(&model.Role{}, id).Error
	})
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	err := s.db.Where("id = ?", id).First(&role).Error
	return &role, err
}

// GetRoles 获取角色列表
func (s *RoleService) GetRoles() ([]model.Role, error) {
	var roles []model.Role
	err := s.db.Order("sort").Find(&roles).Error
	return roles, err
}

// AssignRolesToUser 为用户分配角色
func (s *RoleService) AssignRolesToUser(userID uint, roleIDs []uint) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除原有的用户角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
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
			if err := tx.Create(&userRoles).Error; err != nil {
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
			for _, roleID := range roleIDs {
				// 查询角色编码
				var role model.Role
				if err := tx.Where("id = ?", roleID).First(&role).Error; err != nil {
					return err
				}

				// 添加用户角色关联
				_, err = s.casbinService.AddRoleForUser(userIDStr, role.Code)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetUserRoleIDs 获取用户角色ID列表
func (s *RoleService) GetUserRoleIDs(userID uint) ([]uint, error) {
	var userRoles []model.UserRole
	err := s.db.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	return roleIDs, nil
}

// GetRoleMenuIDs 获取角色菜单ID列表
func (s *RoleService) GetRoleMenuIDs(roleID uint) ([]uint, error) {
	var roleMenus []model.RoleMenu
	err := s.db.Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	return menuIDs, nil
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
