package repository

import (
	"context"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// RoleRepo 角色仓库
type RoleRepo struct {
	DB          *gorm.DB
	GenericRepo Repository[model.Role] // 使用接口而非具体实现
}

// NewRoleRepo 创建角色仓库
func NewRoleRepo(params RepoParams) *RoleRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.Role](params.DB).SetErrorCode(errorx.ErrorNotFoundCode)

	repo := &RoleRepo{
		DB:          params.DB,
		GenericRepo: genericRepo,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.InfoContext(ctx, "RoleRepo initialized")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.InfoContext(ctx, "RoleRepo stopped")
			}
			return nil
		},
	})

	return repo
}

// GetByID 根据ID获取角色
func (r *RoleRepo) GetByID(ctx context.Context, id uint) (*model.Role, error) {
	// 使用仓库接口
	return r.GenericRepo.GetByID(ctx, id)
}

// GetByCode 根据编码获取角色
func (r *RoleRepo) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	// 使用仓库接口的高级查询
	return r.GenericRepo.FindByField(ctx, "code", code)
}

// GetAll 获取所有角色
func (r *RoleRepo) GetAll(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	err := r.DB.WithContext(ctx).Order("sort").Find(&roles).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询所有角色失败")
	}
	return roles, nil
}

// Create 创建角色
func (r *RoleRepo) Create(ctx context.Context, role *model.Role) error {
	// 使用仓库接口
	return r.GenericRepo.Create(ctx, role)
}

// Update 更新角色
func (r *RoleRepo) Update(ctx context.Context, role *model.Role) error {
	// 使用仓库接口
	return r.GenericRepo.Update(ctx, role)
}

// Delete 删除角色
func (r *RoleRepo) Delete(ctx context.Context, id uint) error {
	// 使用仓库接口
	return r.GenericRepo.Delete(ctx, id)
}

// IsAssignedToUser 检查角色是否已分配给用户
func (r *RoleRepo) IsAssignedToUser(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.UserRole{}).Where("role_id = ?", id).Count(&count).Error
	if err != nil {
		return false, errorx.WrapError(err, fmt.Sprintf("检查角色是否已分配给用户失败: roleID=%d", id))
	}
	return count > 0, nil
}

// DeleteRoleMenus 删除角色的菜单关联
func (r *RoleRepo) DeleteRoleMenus(ctx context.Context, roleID uint) error {
	err := r.DB.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除角色的菜单关联失败: roleID=%d", roleID))
	}
	return nil
}

// GetMenuIDsByRoleID 获取角色菜单ID列表
// 使用直接查询避免N+1查询问题
func (r *RoleRepo) GetMenuIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error) {
	// 直接查询菜单ID，避免中间对象转换
	var menuIDs []uint
	err := r.DB.WithContext(ctx).
		Model(&model.RoleMenu{}).
		Select("menu_id").
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取角色菜单ID列表失败: roleID=%d", roleID))
	}

	return menuIDs, nil
}

// GetRoleIDsByUserID 获取用户角色ID列表
// 使用直接查询避免N+1查询问题
func (r *RoleRepo) GetRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	// 直接查询角色ID，避免中间对象转换
	var roleIDs []uint
	err := r.DB.WithContext(ctx).
		Model(&model.UserRole{}).
		Select("role_id").
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取用户角色ID列表失败: userID=%d", userID))
	}

	return roleIDs, nil
}

// AssignRolesToUser 为用户分配角色
func (r *RoleRepo) AssignRolesToUser(ctx context.Context, userID int64, roleIDs []uint) error {
	// 开始事务
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的用户角色关联
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return errorx.WrapError(err, fmt.Sprintf("删除原有的用户角色关联失败: userID=%d", userID))
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
		if err := tx.WithContext(ctx).Create(&userRoles).Error; err != nil {
			tx.Rollback()
			return errorx.WrapError(err, fmt.Sprintf("创建用户角色关联失败: userID=%d, roleIDs=%v", userID, roleIDs))
		}
	}

	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		return errorx.WrapError(err, fmt.Sprintf("提交事务失败: 为用户分配角色, userID=%d, roleIDs=%v", userID, roleIDs))
	}
	return nil
}

// BatchCreateRoleMenus 批量创建角色菜单关联
func (r *RoleRepo) BatchCreateRoleMenus(ctx context.Context, roleMenus []model.RoleMenu) error {
	err := r.DB.WithContext(ctx).Create(&roleMenus).Error
	if err != nil {
		return errorx.WrapError(err, "批量创建角色菜单关联失败")
	}
	return nil
}

// DeleteUserRolesByUserID 删除用户的角色关联
func (r *RoleRepo) DeleteUserRolesByUserID(ctx context.Context, userID int64) error {
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.UserRole{}).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除用户的角色关联失败: userID=%d", userID))
	}
	return nil
}

// AssignMenusToRole 为角色分配菜单
func (r *RoleRepo) AssignMenusToRole(ctx context.Context, roleID uint, menuIDs []uint) error {
	// 开始事务
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的角色菜单关联
	if err := tx.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return errorx.WrapError(err, fmt.Sprintf("删除原有的角色菜单关联失败: roleID=%d", roleID))
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
		if err := tx.WithContext(ctx).Create(&roleMenus).Error; err != nil {
			tx.Rollback()
			return errorx.WrapError(err, fmt.Sprintf("创建角色菜单关联失败: roleID=%d, menuIDs=%v", roleID, menuIDs))
		}
	}

	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		return errorx.WrapError(err, fmt.Sprintf("提交事务失败: 为角色分配菜单, roleID=%d, menuIDs=%v", roleID, menuIDs))
	}
	return nil
}

// GetRoleCodesByUserID 获取用户角色编码列表
// 使用JOIN查询避免N+1查询问题
func (r *RoleRepo) GetRoleCodesByUserID(ctx context.Context, userID int64) ([]string, error) {
	// 使用JOIN查询直接获取角色编码
	var roleCodes []string
	err := r.DB.WithContext(ctx).
		Model(&model.Role{}).
		Select("sys_role.code").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ?", userID).
		Pluck("code", &roleCodes).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取用户角色编码失败: userID=%d", userID))
	}

	// 如果没有角色，返回默认角色
	if len(roleCodes) == 0 {
		return []string{"user"}, nil // 默认角色
	}

	return roleCodes, nil
}

// GetRoleCodesByAdminUserID 获取管理员用户角色编码列表
// 使用JOIN查询避免N+1查询问题
func (r *RoleRepo) GetRoleCodesByAdminUserID(ctx context.Context, adminUserID int64) ([]string, error) {
	// 使用JOIN查询直接获取角色编码
	var roleCodes []string
	err := r.DB.WithContext(ctx).
		Model(&model.Role{}).
		Select("sys_role.code").
		Joins("JOIN admin_user_role ON admin_user_role.role_id = sys_role.id").
		Where("admin_user_role.admin_user_id = ?", adminUserID).
		Pluck("code", &roleCodes).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取管理员用户角色编码失败: adminUserID=%d", adminUserID))
	}

	// 如果没有角色，返回默认角色
	if len(roleCodes) == 0 {
		return []string{"admin"}, nil // 默认管理员角色
	}

	return roleCodes, nil
}

// GetRolesByMenuID 获取拥有指定菜单的所有角色
// 使用JOIN查询避免N+1查询问题
func (r *RoleRepo) GetRolesByMenuID(ctx context.Context, menuID uint) ([]*model.Role, error) {
	// 使用JOIN查询直接获取角色
	var roles []*model.Role
	err := r.DB.WithContext(ctx).
		Model(&model.Role{}).
		Joins("JOIN sys_role_menu ON sys_role_menu.role_id = sys_role.id").
		Where("sys_role_menu.menu_id = ?", menuID).
		Find(&roles).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单关联的角色失败: menuID=%d", menuID))
	}

	return roles, nil
}
