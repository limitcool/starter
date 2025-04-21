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
	genericRepo *GenericRepo[model.Role] // 泛型仓库
}

// NewRoleRepo 创建角色仓库
func NewRoleRepo(params RepoParams) *RoleRepo {
	genericRepo := NewGenericRepo[model.Role](params.DB)
	genericRepo.SetErrorCode(errorx.ErrorNotFoundCode) // 设置错误码

	repo := &RoleRepo{
		DB:          params.DB,
		genericRepo: genericRepo,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.Info("RoleRepo initialized")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.Info("RoleRepo stopped")
			}
			return nil
		},
	})

	return repo
}

// GetByID 根据ID获取角色
func (r *RoleRepo) GetByID(ctx context.Context, id uint) (*model.Role, error) {
	// 使用泛型仓库
	return r.genericRepo.GetByID(ctx, id)
}

// GetByCode 根据编码获取角色
func (r *RoleRepo) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	// 使用泛型仓库的高级查询
	return r.genericRepo.FindByField(ctx, "code", code)
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
	// 使用泛型仓库
	return r.genericRepo.Create(ctx, role)
}

// Update 更新角色
func (r *RoleRepo) Update(ctx context.Context, role *model.Role) error {
	// 使用泛型仓库
	return r.genericRepo.Update(ctx, role)
}

// Delete 删除角色
func (r *RoleRepo) Delete(ctx context.Context, id uint) error {
	// 使用泛型仓库
	return r.genericRepo.Delete(ctx, id)
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
func (r *RoleRepo) GetMenuIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error) {
	var roleMenus []model.RoleMenu
	err := r.DB.WithContext(ctx).Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取角色菜单ID列表失败: roleID=%d", roleID))
	}

	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	return menuIDs, nil
}

// GetRoleIDsByUserID 获取用户角色ID列表
func (r *RoleRepo) GetRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	var userRoles []model.UserRole
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取用户角色ID列表失败: userID=%d", userID))
	}

	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
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

// GetRolesByMenuID 获取拥有指定菜单的所有角色
func (r *RoleRepo) GetRolesByMenuID(ctx context.Context, menuID uint) ([]*model.Role, error) {
	// 查询菜单关联的角色ID
	var roleMenus []model.RoleMenu
	err := r.DB.WithContext(ctx).Where("menu_id = ?", menuID).Find(&roleMenus).Error
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单关联的角色失败: menuID=%d", menuID))
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
	err = r.DB.WithContext(ctx).Where("id IN ?", roleIDs).Find(&roles).Error
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询角色失败: roleIDs=%v", roleIDs))
	}
	return roles, nil
}
