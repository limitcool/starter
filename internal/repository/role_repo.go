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
	genericRepo := NewGenericRepo[model.Role](params.DB).SetErrorCode(errorx.ErrNotFoundCodeValue)

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
	return r.GenericRepo.Get(ctx, id, nil)
}

// GetByCode 根据编码获取角色
func (r *RoleRepo) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	// 使用仓库接口的高级查询
	opts := &QueryOptions{
		Condition: "code = ?",
		Args:      []any{code},
	}
	return r.GenericRepo.Get(ctx, nil, opts)
}

// GetAll 获取所有角色
func (r *RoleRepo) GetAll(ctx context.Context) ([]model.Role, error) {
	// 使用泛型仓库的List方法，按照sort字段排序
	opts := &QueryOptions{
		Condition: "1=1 ORDER BY sort",
	}

	// 获取所有角色（不分页，传入一个很大的页大小）
	roles, err := r.GenericRepo.List(ctx, 1, 1000, opts)
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
	// 创建UserRole的泛型仓库
	userRoleRepo := NewGenericRepo[model.UserRole](r.DB)

	// 使用泛型仓库的Count方法
	opts := &QueryOptions{
		Condition: "role_id = ?",
		Args:      []any{id},
	}

	count, err := userRoleRepo.Count(ctx, opts)
	if err != nil {
		return false, errorx.WrapError(err, fmt.Sprintf("检查角色是否已分配给用户失败: roleID=%d", id))
	}
	return count > 0, nil
}

// DeleteRoleMenus 删除角色的菜单关联
func (r *RoleRepo) DeleteRoleMenus(ctx context.Context, roleID uint) error {
	// 创建RoleMenu的泛型仓库
	roleMenuRepo := NewGenericRepo[model.RoleMenu](r.DB)

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "role_id = ?",
		Args:      []any{roleID},
	}

	// 获取所有关联记录
	roleMenus, err := roleMenuRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("查询角色菜单关联失败: roleID=%d", roleID))
	}

	// 删除所有关联记录
	for _, roleMenu := range roleMenus {
		if err := roleMenuRepo.Delete(ctx, roleMenu.ID); err != nil {
			return errorx.WrapError(err, fmt.Sprintf("删除角色菜单关联失败: roleID=%d, menuID=%d", roleID, roleMenu.MenuID))
		}
	}

	return nil
}

// GetMenuIDsByRoleID 获取角色菜单ID列表
// 使用直接查询避免N+1查询问题
func (r *RoleRepo) GetMenuIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error) {
	// 创建RoleMenu的泛型仓库
	roleMenuRepo := NewGenericRepo[model.RoleMenu](r.DB)

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "role_id = ?",
		Args:      []any{roleID},
	}

	// 获取所有关联记录
	roleMenus, err := roleMenuRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取角色菜单ID列表失败: roleID=%d", roleID))
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, roleMenu := range roleMenus {
		menuIDs = append(menuIDs, roleMenu.MenuID)
	}

	return menuIDs, nil
}

// GetRoleIDsByUserID 获取用户角色ID列表
// 使用直接查询避免N+1查询问题
func (r *RoleRepo) GetRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	// 创建UserRole的泛型仓库
	userRoleRepo := NewGenericRepo[model.UserRole](r.DB)

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "user_id = ?",
		Args:      []any{userID},
	}

	// 获取所有关联记录
	userRoles, err := userRoleRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取用户角色ID列表失败: userID=%d", userID))
	}

	// 提取角色ID
	var roleIDs []uint
	for _, userRole := range userRoles {
		roleIDs = append(roleIDs, userRole.RoleID)
	}

	return roleIDs, nil
}

// AssignRolesToUser 为用户分配角色
func (r *RoleRepo) AssignRolesToUser(ctx context.Context, userID int64, roleIDs []uint) error {
	// 使用泛型仓库的事务支持
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 创建UserRole的泛型仓库
		userRoleRepo := NewGenericRepo[model.UserRole](tx)

		// 使用查询选项
		opts := &QueryOptions{
			Condition: "user_id = ?",
			Args:      []any{userID},
		}

		// 获取所有关联记录
		userRoles, err := userRoleRepo.List(ctx, 1, 1000, opts)
		if err != nil {
			return errorx.WrapError(err, fmt.Sprintf("查询用户角色关联失败: userID=%d", userID))
		}

		// 删除所有关联记录
		for _, userRole := range userRoles {
			if err := userRoleRepo.Delete(ctx, userRole.ID); err != nil {
				return errorx.WrapError(err, fmt.Sprintf("删除用户角色关联失败: userID=%d, roleID=%d", userID, userRole.RoleID))
			}
		}

		// 添加新的用户角色关联
		if len(roleIDs) > 0 {
			// 批量创建用户角色关联
			for _, roleID := range roleIDs {
				userRole := model.UserRole{
					UserID: userID,
					RoleID: roleID,
				}
				if err := userRoleRepo.Create(ctx, &userRole); err != nil {
					return errorx.WrapError(err, fmt.Sprintf("创建用户角色关联失败: userID=%d, roleID=%d", userID, roleID))
				}
			}
		}

		return nil
	})
}

// BatchCreateRoleMenus 批量创建角色菜单关联
func (r *RoleRepo) BatchCreateRoleMenus(ctx context.Context, roleMenus []model.RoleMenu) error {
	// 创建RoleMenu的泛型仓库
	roleMenuRepo := NewGenericRepo[model.RoleMenu](r.DB)

	// 使用事务批量创建
	return r.DB.Transaction(func(tx *gorm.DB) error {
		txRepo := roleMenuRepo.WithTx(tx)

		// 逐个创建角色菜单关联
		for i := range roleMenus {
			if err := txRepo.Create(ctx, &roleMenus[i]); err != nil {
				return errorx.WrapError(err, "批量创建角色菜单关联失败")
			}
		}

		return nil
	})
}

// DeleteUserRolesByUserID 删除用户的角色关联
func (r *RoleRepo) DeleteUserRolesByUserID(ctx context.Context, userID int64) error {
	// 创建UserRole的泛型仓库
	userRoleRepo := NewGenericRepo[model.UserRole](r.DB)

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "user_id = ?",
		Args:      []any{userID},
	}

	// 获取所有关联记录
	userRoles, err := userRoleRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("查询用户角色关联失败: userID=%d", userID))
	}

	// 删除所有关联记录
	for _, userRole := range userRoles {
		if err := userRoleRepo.Delete(ctx, userRole.ID); err != nil {
			return errorx.WrapError(err, fmt.Sprintf("删除用户角色关联失败: userID=%d, roleID=%d", userID, userRole.RoleID))
		}
	}

	return nil
}

// AssignMenusToRole 为角色分配菜单
func (r *RoleRepo) AssignMenusToRole(ctx context.Context, roleID uint, menuIDs []uint) error {
	// 使用泛型仓库的事务支持
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 创建RoleMenu的泛型仓库
		roleMenuRepo := NewGenericRepo[model.RoleMenu](tx)

		// 使用查询选项
		opts := &QueryOptions{
			Condition: "role_id = ?",
			Args:      []any{roleID},
		}

		// 获取所有关联记录
		roleMenus, err := roleMenuRepo.List(ctx, 1, 1000, opts)
		if err != nil {
			return errorx.WrapError(err, fmt.Sprintf("查询角色菜单关联失败: roleID=%d", roleID))
		}

		// 删除所有关联记录
		for _, roleMenu := range roleMenus {
			if err := roleMenuRepo.Delete(ctx, roleMenu.ID); err != nil {
				return errorx.WrapError(err, fmt.Sprintf("删除角色菜单关联失败: roleID=%d, menuID=%d", roleID, roleMenu.MenuID))
			}
		}

		// 添加新的角色菜单关联
		if len(menuIDs) > 0 {
			// 逐个创建角色菜单关联
			for _, menuID := range menuIDs {
				roleMenu := model.RoleMenu{
					RoleID: roleID,
					MenuID: menuID,
				}
				if err := roleMenuRepo.Create(ctx, &roleMenu); err != nil {
					return errorx.WrapError(err, fmt.Sprintf("创建角色菜单关联失败: roleID=%d, menuID=%d", roleID, menuID))
				}
			}
		}

		return nil
	})
}

// GetRoleCodesByUserID 获取用户角色编码列表
// 使用JOIN查询避免N+1查询问题
func (r *RoleRepo) GetRoleCodesByUserID(ctx context.Context, userID int64) ([]string, error) {
	// 首先获取用户的角色ID
	roleIDs, err := r.GetRoleIDsByUserID(ctx, uint(userID))
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取用户角色ID失败: userID=%d", userID))
	}

	if len(roleIDs) == 0 {
		return []string{"user"}, nil // 默认角色
	}

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "id IN ?",
		Args:      []any{roleIDs},
	}

	// 获取所有角色
	roles, err := r.GenericRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取用户角色失败: userID=%d", userID))
	}

	// 提取角色编码
	var roleCodes []string
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
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
	// 创建AdminUserRole的泛型仓库
	adminUserRoleRepo := NewGenericRepo[model.AdminUserRole](r.DB)

	// 使用查询选项
	adminUserRoleOpts := &QueryOptions{
		Condition: "admin_user_id = ?",
		Args:      []any{adminUserID},
	}

	// 获取所有关联记录
	adminUserRoles, err := adminUserRoleRepo.List(ctx, 1, 1000, adminUserRoleOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询管理员用户角色关联失败: adminUserID=%d", adminUserID))
	}

	// 提取角色ID
	var roleIDs []uint
	for _, adminUserRole := range adminUserRoles {
		roleIDs = append(roleIDs, adminUserRole.RoleID)
	}

	if len(roleIDs) == 0 {
		return []string{"admin"}, nil // 默认管理员角色
	}

	// 使用查询选项
	roleOpts := &QueryOptions{
		Condition: "id IN ?",
		Args:      []any{roleIDs},
	}

	// 获取所有角色
	roles, err := r.GenericRepo.List(ctx, 1, 1000, roleOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取管理员用户角色失败: adminUserID=%d", adminUserID))
	}

	// 提取角色编码
	var roleCodes []string
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
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
	// 创建RoleMenu的泛型仓库
	roleMenuRepo := NewGenericRepo[model.RoleMenu](r.DB)

	// 使用查询选项
	roleMenuOpts := &QueryOptions{
		Condition: "menu_id = ?",
		Args:      []any{menuID},
	}

	// 获取所有关联记录
	roleMenus, err := roleMenuRepo.List(ctx, 1, 1000, roleMenuOpts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单角色关联失败: menuID=%d", menuID))
	}

	// 提取角色ID
	var roleIDs []uint
	for _, roleMenu := range roleMenus {
		roleIDs = append(roleIDs, roleMenu.RoleID)
	}

	if len(roleIDs) == 0 {
		return []*model.Role{}, nil
	}

	// 使用查询选项
	opts := &QueryOptions{
		Condition: "id IN ?",
		Args:      []any{roleIDs},
	}

	// 获取所有角色
	roles, err := r.GenericRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单关联的角色失败: menuID=%d", menuID))
	}

	// 转换为指针切片
	rolePtrs := make([]*model.Role, len(roles))
	for i := range roles {
		rolePtrs[i] = &roles[i]
	}

	return rolePtrs, nil
}
