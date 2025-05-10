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

// AdminUserRepo 管理员用户仓库
type AdminUserRepo struct {
	DB          *gorm.DB
	GenericRepo Repository[model.AdminUser] // 使用接口而非具体实现
}

// NewAdminUserRepo 创建管理员用户仓库
func NewAdminUserRepo(params RepoParams) *AdminUserRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.AdminUser](params.DB).SetErrorCode(errorx.ErrorUserNotFoundCodeValue)

	repo := &AdminUserRepo{
		DB:          params.DB,
		GenericRepo: genericRepo,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.Info("AdminUserRepo initialized")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.Info("AdminUserRepo stopped")
			}
			return nil
		},
	})

	return repo
}

// GetByID 根据ID获取管理员用户
func (r *AdminUserRepo) GetByID(ctx context.Context, id int64) (*model.AdminUser, error) {
	// 使用仓库接口获取用户，并预加载角色
	opts := &QueryOptions{
		Preloads: []string{"Roles"},
	}

	user, err := r.GenericRepo.Get(ctx, id, opts)
	if err != nil {
		return nil, err
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return user, nil
}

// GetByUsername 根据用户名获取管理员用户
func (r *AdminUserRepo) GetByUsername(ctx context.Context, username string) (*model.AdminUser, error) {
	// 使用仓库接口获取用户，并预加载角色
	opts := &QueryOptions{
		Condition: "username = ?",
		Args:      []any{username},
		Preloads:  []string{"Roles"},
	}

	user, err := r.GenericRepo.Get(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return user, nil
}

// Create 创建管理员用户
func (r *AdminUserRepo) Create(ctx context.Context, user *model.AdminUser) error {
	// 使用仓库接口
	return r.GenericRepo.Create(ctx, user)
}

// Update 更新管理员用户
func (r *AdminUserRepo) Update(ctx context.Context, user *model.AdminUser) error {
	// 使用仓库接口
	return r.GenericRepo.Update(ctx, user)
}

// UpdateFields 更新管理员用户字段
func (r *AdminUserRepo) UpdateFields(ctx context.Context, id int64, fields map[string]any) error {
	// 使用仓库接口
	return r.GenericRepo.UpdateFields(ctx, id, fields)
}

// Delete 删除管理员用户
func (r *AdminUserRepo) Delete(ctx context.Context, id int64) error {
	// 使用仓库接口
	return r.GenericRepo.Delete(ctx, id)
}

// List 获取管理员用户列表
func (r *AdminUserRepo) List(ctx context.Context, query *model.AdminUserQuery) ([]model.AdminUser, int64, error) {
	// 构建查询条件
	conditions := []string{}
	args := []any{}

	if query.Username != "" {
		conditions = append(conditions, "username LIKE ?")
		args = append(args, "%"+query.Username+"%")
	}
	if query.Nickname != "" {
		conditions = append(conditions, "nickname LIKE ?")
		args = append(args, "%"+query.Nickname+"%")
	}
	if query.Email != "" {
		conditions = append(conditions, "email LIKE ?")
		args = append(args, "%"+query.Email+"%")
	}
	if query.Phone != "" {
		conditions = append(conditions, "phone LIKE ?")
		args = append(args, "%"+query.Phone+"%")
	}
	if query.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *query.Status)
	}

	// 组合条件
	condition := ""
	if len(conditions) > 0 {
		condition = conditions[0]
		for i := 1; i < len(conditions); i++ {
			condition += " AND " + conditions[i]
		}
	}

	// 添加排序
	if query.OrderBy == "" {
		// 设置默认排序
		condition += " ORDER BY id DESC"
	} else {
		direction := "ASC"
		if query.OrderDesc {
			direction = "DESC"
		}
		condition += " ORDER BY " + query.OrderBy + " " + direction
	}

	// 使用仓库接口的分页查询
	return r.GenericRepo.GetPage(ctx, int(query.Page), int(query.PageSize), condition, args...)
}

// UpdateAvatar 更新管理员用户头像
func (r *AdminUserRepo) UpdateAvatar(ctx context.Context, userID int64, fileID uint) error {
	// 使用仓库接口的事务支持
	return r.GenericRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// 创建事务中的仓库
		txRepo := r.GenericRepo.WithTx(tx)

		// 查找用户
		opts := &QueryOptions{
			Preloads: []string{"Roles"},
		}
		user, err := txRepo.Get(ctx, userID, opts)
		if err != nil {
			return err
		}

		// 更新用户头像
		user.AvatarFileID = fileID

		// 保存用户
		return txRepo.Update(ctx, user)
	})
}

// AssignRolesToAdminUser 为管理员用户分配角色
func (r *AdminUserRepo) AssignRolesToAdminUser(ctx context.Context, adminUserID int64, roleIDs []uint) error {
	// 使用泛型仓库的事务支持
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 创建AdminUserRole的泛型仓库
		adminUserRoleRepo := NewGenericRepo[model.AdminUserRole](tx)

		// 使用查询选项
		opts := &QueryOptions{
			Condition: "admin_user_id = ?",
			Args:      []any{adminUserID},
		}

		// 获取所有关联记录
		adminUserRoles, err := adminUserRoleRepo.List(ctx, 1, 1000, opts)
		if err != nil {
			return errorx.WrapError(err, fmt.Sprintf("查询管理员用户角色关联失败: adminUserID=%d", adminUserID))
		}

		// 删除所有关联记录
		for _, adminUserRole := range adminUserRoles {
			if err := adminUserRoleRepo.Delete(ctx, adminUserRole.ID); err != nil {
				return errorx.WrapError(err, fmt.Sprintf("删除管理员用户角色关联失败: adminUserID=%d, roleID=%d", adminUserID, adminUserRole.RoleID))
			}
		}

		// 添加新的管理员用户角色关联
		if len(roleIDs) > 0 {
			// 批量创建管理员用户角色关联
			for _, roleID := range roleIDs {
				adminUserRole := model.AdminUserRole{
					AdminUserID: adminUserID,
					RoleID:      roleID,
				}
				if err := adminUserRoleRepo.Create(ctx, &adminUserRole); err != nil {
					return errorx.WrapError(err, fmt.Sprintf("创建管理员用户角色关联失败: adminUserID=%d, roleID=%d", adminUserID, roleID))
				}
			}
		}

		return nil
	})
}
