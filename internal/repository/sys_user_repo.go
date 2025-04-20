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

// SysUserRepo 系统用户仓库
type SysUserRepo struct {
	DB          *gorm.DB
	genericRepo *GenericRepo[model.SysUser] // 泛型仓库
}

// NewSysUserRepo 创建系统用户仓库
func NewSysUserRepo(params RepoParams) *SysUserRepo {
	genericRepo := NewGenericRepo[model.SysUser](params.DB)
	genericRepo.SetErrorCode(errorx.ErrorUserNotFoundCode) // 设置错误码

	repo := &SysUserRepo{
		DB:          params.DB,
		genericRepo: genericRepo,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.Info("SysUserRepo initialized")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.Info("SysUserRepo stopped")
			}
			return nil
		},
	})

	return repo
}

// GetByID 根据ID获取系统用户
func (r *SysUserRepo) GetByID(ctx context.Context, id int64) (*model.SysUser, error) {
	// 使用泛型仓库获取用户
	user, err := r.genericRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取用户的角色
	if err := r.DB.WithContext(ctx).Model(user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询系统用户角色失败: id=%d", id))
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return user, nil
}

// GetByUsername 根据用户名获取系统用户
func (r *SysUserRepo) GetByUsername(ctx context.Context, username string) (*model.SysUser, error) {
	// 使用泛型仓库获取用户
	user, err := r.genericRepo.FindByField(ctx, "username", username)
	if err != nil {
		return nil, err
	}

	// 获取用户的角色
	if err := r.DB.WithContext(ctx).Model(user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询系统用户角色失败: username=%s", username))
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return user, nil
}

// Create 创建系统用户
func (r *SysUserRepo) Create(ctx context.Context, user *model.SysUser) error {
	// 使用泛型仓库
	return r.genericRepo.Create(ctx, user)
}

// Update 更新系统用户
func (r *SysUserRepo) Update(ctx context.Context, user *model.SysUser) error {
	// 使用泛型仓库
	return r.genericRepo.Update(ctx, user)
}

// UpdateFields 更新系统用户字段
func (r *SysUserRepo) UpdateFields(ctx context.Context, id int64, fields map[string]any) error {
	// 使用泛型仓库
	return r.genericRepo.UpdateFields(ctx, id, fields)
}

// Delete 删除系统用户
func (r *SysUserRepo) Delete(ctx context.Context, id int64) error {
	// 使用泛型仓库
	return r.genericRepo.Delete(ctx, id)
}

// List 获取系统用户列表
func (r *SysUserRepo) List(ctx context.Context, query *model.SysUserQuery) ([]model.SysUser, int64, error) {
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

	// 设置默认排序
	if query.OrderBy == "" {
		// 在查询前设置默认排序
		r.DB = r.DB.Order("id DESC")
	} else {
		direction := "ASC"
		if query.OrderDesc {
			direction = "DESC"
		}
		r.DB = r.DB.Order(query.OrderBy + " " + direction)
	}

	// 使用泛型仓库的分页查询
	return r.genericRepo.GetPage(ctx, int(query.Page), int(query.PageSize), condition, args...)
}

// UpdateAvatar 更新系统用户头像
func (r *SysUserRepo) UpdateAvatar(ctx context.Context, userID int64, fileID uint) error {
	// 使用泛型仓库的事务支持
	return r.genericRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// 创建事务中的泛型仓库
		txRepo := r.genericRepo.WithTx(tx).(*GenericRepo[model.SysUser])

		// 查找用户
		user, err := txRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		// 更新用户头像
		user.AvatarFileID = fileID

		// 保存用户
		return txRepo.Update(ctx, user)
	})
}
