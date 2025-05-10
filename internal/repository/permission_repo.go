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
	GenericRepo Repository[model.Permission] // 使用接口而非具体实现
}

// NewPermissionRepo 创建权限仓库
func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.Permission](db).SetErrorCode(errorx.ErrNotFoundCodeValue)

	return &PermissionRepo{
		DB:          db,
		GenericRepo: genericRepo,
	}
}

// GetByID 根据ID获取权限
func (r *PermissionRepo) GetByID(ctx context.Context, id uint) (*model.Permission, error) {
	// 使用仓库接口
	return r.GenericRepo.GetByID(ctx, id)
}

// GetAll 获取所有权限
func (r *PermissionRepo) GetAll(ctx context.Context) ([]model.Permission, error) {
	// 使用泛型仓库的List方法
	permissions, err := r.GenericRepo.List(ctx, 1, 1000, nil)
	if err != nil {
		return nil, errorx.WrapError(err, "查询所有权限失败")
	}
	return permissions, nil
}

// Create 创建权限
func (r *PermissionRepo) Create(ctx context.Context, permission *model.Permission) error {
	// 使用仓库接口
	return r.GenericRepo.Create(ctx, permission)
}

// Update 更新权限
func (r *PermissionRepo) Update(ctx context.Context, permission *model.Permission) error {
	// 使用仓库接口
	return r.GenericRepo.Update(ctx, permission)
}

// Delete 删除权限
func (r *PermissionRepo) Delete(ctx context.Context, id uint) error {
	// 使用仓库接口
	return r.GenericRepo.Delete(ctx, id)
}

// GetByRoleID 获取角色的权限列表
func (r *PermissionRepo) GetByRoleID(ctx context.Context, roleID uint) ([]model.Permission, error) {
	// 创建RolePermission的泛型仓库
	rolePermRepo := NewGenericRepo[model.RolePermission](r.DB)

	// 使用查询选项
	rolePermOpts := &QueryOptions{
		Condition: "role_id = ?",
		Args:      []any{roleID},
	}

	// 获取所有关联记录
	rolePerms, err := rolePermRepo.List(ctx, 1, 1000, rolePermOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询角色权限关联失败: roleID=%d", roleID))
	}

	// 提取权限ID
	var permissionIDs []uint
	for _, rolePerm := range rolePerms {
		permissionIDs = append(permissionIDs, rolePerm.PermissionID)
	}

	if len(permissionIDs) == 0 {
		return []model.Permission{}, nil
	}

	// 使用泛型仓库的List方法
	opts := &QueryOptions{
		Condition: "id IN ?",
		Args:      []any{permissionIDs},
	}

	permissions, err := r.GenericRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询角色权限失败: roleID=%d", roleID))
	}

	return permissions, nil
}

// GetByUserID 获取用户的权限列表
func (r *PermissionRepo) GetByUserID(ctx context.Context, userID uint) ([]model.Permission, error) {
	// 创建UserRole的泛型仓库
	userRoleRepo := NewGenericRepo[model.UserRole](r.DB)

	// 使用查询选项
	userRoleOpts := &QueryOptions{
		Condition: "user_id = ?",
		Args:      []any{userID},
	}

	// 获取所有关联记录
	userRoles, err := userRoleRepo.List(ctx, 1, 1000, userRoleOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询用户角色关联失败: userID=%d", userID))
	}

	// 提取角色ID
	var roleIDs []uint
	for _, userRole := range userRoles {
		roleIDs = append(roleIDs, userRole.RoleID)
	}

	if len(roleIDs) == 0 {
		return []model.Permission{}, nil
	}

	// 创建RolePermission的泛型仓库
	rolePermRepo := NewGenericRepo[model.RolePermission](r.DB)

	// 使用查询选项
	rolePermOpts := &QueryOptions{
		Condition: "role_id IN ?",
		Args:      []any{roleIDs},
	}

	// 获取所有关联记录
	rolePerms, err := rolePermRepo.List(ctx, 1, 1000, rolePermOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询角色权限关联失败: roleIDs=%v", roleIDs))
	}

	// 提取权限ID（去重）
	permMap := make(map[uint]struct{})
	for _, rolePerm := range rolePerms {
		permMap[rolePerm.PermissionID] = struct{}{}
	}

	var permissionIDs []uint
	for permID := range permMap {
		permissionIDs = append(permissionIDs, permID)
	}

	if len(permissionIDs) == 0 {
		return []model.Permission{}, nil
	}

	// 使用泛型仓库的List方法
	opts := &QueryOptions{
		Condition: "id IN ?",
		Args:      []any{permissionIDs},
	}

	permissions, err := r.GenericRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询用户权限失败: userID=%d", userID))
	}

	return permissions, nil
}

// AssignPermissionToRole 为角色分配权限
func (r *PermissionRepo) AssignPermissionToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	// 创建Role的泛型仓库
	roleRepo := NewGenericRepo[model.Role](r.DB)

	// 检查角色是否存在
	_, err := roleRepo.Get(ctx, roleID, nil)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("查询角色失败: roleID=%d", roleID))
	}

	// 使用泛型仓库的事务支持
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 创建RolePermission的泛型仓库
		rolePermRepo := NewGenericRepo[model.RolePermission](tx)

		// 使用查询选项
		opts := &QueryOptions{
			Condition: "role_id = ?",
			Args:      []any{roleID},
		}

		// 获取所有关联记录
		rolePerms, err := rolePermRepo.List(ctx, 1, 1000, opts)
		if err != nil {
			return errorx.WrapError(err, fmt.Sprintf("查询角色权限关联失败: roleID=%d", roleID))
		}

		// 删除所有关联记录
		for _, rolePerm := range rolePerms {
			if err := rolePermRepo.Delete(ctx, rolePerm.ID); err != nil {
				return errorx.WrapError(err, fmt.Sprintf("删除角色权限关联失败: roleID=%d, permissionID=%d", roleID, rolePerm.PermissionID))
			}
		}

		// 添加新的角色权限关联
		if len(permissionIDs) > 0 {

			// 逐个创建角色权限关联
			for _, permID := range permissionIDs {
				rolePermission := model.RolePermission{
					RoleID:       roleID,
					PermissionID: permID,
				}
				if err := rolePermRepo.Create(ctx, &rolePermission); err != nil {
					return errorx.WrapError(err, fmt.Sprintf("创建角色权限关联失败: roleID=%d, permissionID=%d", roleID, permID))
				}
			}
		}

		return nil
	})
}
