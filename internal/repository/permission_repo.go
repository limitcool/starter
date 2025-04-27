package repository

import (
	"context"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// PermissionRepo 权限仓库
type PermissionRepo struct {
	DB          *gorm.DB
	genericRepo Repository[model.Permission] // 使用接口而非具体实现
}

// NewPermissionRepo 创建权限仓库
func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.Permission](db).SetErrorCode(errorx.ErrorNotFoundCode)

	return &PermissionRepo{
		DB:          db,
		genericRepo: genericRepo,
	}
}

// GetByID 根据ID获取权限
func (r *PermissionRepo) GetByID(ctx context.Context, id uint) (*model.Permission, error) {
	// 使用仓库接口
	return r.genericRepo.GetByID(ctx, id)
}

// GetAll 获取所有权限
func (r *PermissionRepo) GetAll(ctx context.Context) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.DB.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询所有权限失败")
	}
	return permissions, nil
}

// Create 创建权限
func (r *PermissionRepo) Create(ctx context.Context, permission *model.Permission) error {
	// 使用仓库接口
	return r.genericRepo.Create(ctx, permission)
}

// Update 更新权限
func (r *PermissionRepo) Update(ctx context.Context, permission *model.Permission) error {
	// 使用仓库接口
	return r.genericRepo.Update(ctx, permission)
}

// Delete 删除权限
func (r *PermissionRepo) Delete(ctx context.Context, id uint) error {
	// 使用仓库接口
	return r.genericRepo.Delete(ctx, id)
}

// GetByRoleID 获取角色的权限列表
func (r *PermissionRepo) GetByRoleID(ctx context.Context, roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission

	// 通过关联表查询
	err := r.DB.WithContext(ctx).Joins("JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Where("sys_role_permission.role_id = ?", roleID).
		Find(&permissions).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询角色权限失败: roleID=%d", roleID))
	}
	return permissions, nil
}

// GetByUserID 获取用户的权限列表
func (r *PermissionRepo) GetByUserID(ctx context.Context, userID uint) ([]model.Permission, error) {
	var permissions []model.Permission

	// 通过用户角色关联查询权限
	err := r.DB.WithContext(ctx).Joins("JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_permission.role_id").
		Where("sys_user_role.user_id = ?", userID).
		Find(&permissions).Error

	return permissions, err
}

// AssignPermissionToRole 为角色分配权限
func (r *PermissionRepo) AssignPermissionToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	// 获取角色
	var role model.Role
	if err := r.DB.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return errorx.ErrNotFound.WithError(err)
	}

	// 开始事务
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的角色权限关联
	if err := tx.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
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
		if err := tx.WithContext(ctx).Create(&rolePermissions).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
