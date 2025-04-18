package repository

import (
	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

// RoleRepo 角色仓库
type RoleRepo struct {
	DB *gorm.DB
}

// NewRoleRepo 创建角色仓库
func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{DB: db}
}

// GetByID 根据ID获取角色
func (r *RoleRepo) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.DB.Where("id = ?", id).First(&role).Error
	return &role, err
}

// GetByCode 根据编码获取角色
func (r *RoleRepo) GetByCode(code string) (*model.Role, error) {
	var role model.Role
	err := r.DB.Where("code = ?", code).First(&role).Error
	return &role, err
}

// GetAll 获取所有角色
func (r *RoleRepo) GetAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.DB.Order("sort").Find(&roles).Error
	return roles, err
}

// Create 创建角色
func (r *RoleRepo) Create(role *model.Role) error {
	return r.DB.Create(role).Error
}

// Update 更新角色
func (r *RoleRepo) Update(role *model.Role) error {
	return r.DB.Model(&model.Role{}).Where("id = ?", role.ID).Updates(role).Error
}

// Delete 删除角色
func (r *RoleRepo) Delete(id uint) error {
	return r.DB.Delete(&model.Role{}, id).Error
}

// IsAssignedToUser 检查角色是否已分配给用户
func (r *RoleRepo) IsAssignedToUser(id uint) (bool, error) {
	var count int64
	err := r.DB.Model(&model.UserRole{}).Where("role_id = ?", id).Count(&count).Error
	return count > 0, err
}

// DeleteRoleMenus 删除角色的菜单关联
func (r *RoleRepo) DeleteRoleMenus(roleID uint) error {
	return r.DB.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error
}

// GetMenuIDsByRoleID 获取角色菜单ID列表
func (r *RoleRepo) GetMenuIDsByRoleID(roleID uint) ([]uint, error) {
	var roleMenus []model.RoleMenu
	err := r.DB.Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	return menuIDs, nil
}

// GetRoleIDsByUserID 获取用户角色ID列表
func (r *RoleRepo) GetRoleIDsByUserID(userID uint) ([]uint, error) {
	var userRoles []model.UserRole
	err := r.DB.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	return roleIDs, nil
}

// AssignRolesToUser 为用户分配角色
func (r *RoleRepo) AssignRolesToUser(userID int64, roleIDs []uint) error {
	// 开始事务
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的用户角色关联
	if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// BatchCreateRoleMenus 批量创建角色菜单关联
func (r *RoleRepo) BatchCreateRoleMenus(roleMenus []model.RoleMenu) error {
	return r.DB.Create(&roleMenus).Error
}

// DeleteUserRolesByUserID 删除用户的角色关联
func (r *RoleRepo) DeleteUserRolesByUserID(userID int64) error {
	return r.DB.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error
}

// AssignMenusToRole 为角色分配菜单
func (r *RoleRepo) AssignMenusToRole(roleID uint, menuIDs []uint) error {
	// 开始事务
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的角色菜单关联
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新的角色菜单关联
	if len(menuIDs) > 0 {
		var roleMenus []model.RoleMenu
		for _, menuID := range menuIDs {
			roleMenus = append(roleMenus, model.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			})
		}
		if err := tx.Create(&roleMenus).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetRolesByMenuID 获取拥有指定菜单的所有角色
func (r *RoleRepo) GetRolesByMenuID(menuID uint) ([]*model.Role, error) {
	// 查询菜单关联的角色ID
	var roleMenus []model.RoleMenu
	err := r.DB.Where("menu_id = ?", menuID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取角色ID
	var roleIDs []uint
	for _, rm := range roleMenus {
		roleIDs = append(roleIDs, rm.RoleID)
	}

	if len(roleIDs) == 0 {
		return []*model.Role{}, nil
	}

	// 查询角色
	var roles []*model.Role
	err = r.DB.Where("id IN ?", roleIDs).Find(&roles).Error
	return roles, err
}
