package repository

import (
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// PermissionRepo 权限仓库
type PermissionRepo struct {
	DB *gorm.DB
}

// NewPermissionRepo 创建权限仓库
func NewPermissionRepo(db *gorm.DB) PermissionRepository {
	return &PermissionRepo{DB: db}
}

// GetByID 根据ID获取权限
func (r *PermissionRepo) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.DB.First(&permission, id).Error
	return &permission, err
}

// GetAll 获取所有权限
func (r *PermissionRepo) GetAll() ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.DB.Find(&permissions).Error
	return permissions, err
}

// Create 创建权限
func (r *PermissionRepo) Create(permission *model.Permission) error {
	return r.DB.Create(permission).Error
}

// Update 更新权限
func (r *PermissionRepo) Update(permission *model.Permission) error {
	return r.DB.Model(&model.Permission{}).Where("id = ?", permission.ID).Updates(permission).Error
}

// Delete 删除权限
func (r *PermissionRepo) Delete(id uint) error {
	return r.DB.Delete(&model.Permission{}, id).Error
}

// GetByRoleID 获取角色的权限列表
func (r *PermissionRepo) GetByRoleID(roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission

	// 通过关联表查询
	err := r.DB.Joins("JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Where("sys_role_permission.role_id = ?", roleID).
		Find(&permissions).Error

	return permissions, err
}

// GetByUserID 获取用户的权限列表
func (r *PermissionRepo) GetByUserID(userID uint) ([]model.Permission, error) {
	var permissions []model.Permission

	// 通过用户角色关联查询权限
	err := r.DB.Joins("JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_permission.role_id").
		Where("sys_user_role.user_id = ?", userID).
		Find(&permissions).Error

	return permissions, err
}

// AssignPermissionToRole 为角色分配权限
func (r *PermissionRepo) AssignPermissionToRole(roleID uint, permissionIDs []uint) error {
	// 获取角色
	var role model.Role
	if err := r.DB.First(&role, roleID).Error; err != nil {
		return errorx.ErrNotFound.WithError(err)
	}

	// 开始事务
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的角色权限关联
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

	return tx.Commit().Error
}
