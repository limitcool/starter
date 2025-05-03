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

// UserRepo 用户仓库
// 提供用户相关的数据库操作
type UserRepo struct {
	DB          *gorm.DB
	GenericRepo Repository[model.User] // 使用接口而非具体实现
}

// NewUserRepo 创建用户仓库
func NewUserRepo(params RepoParams) *UserRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.User](params.DB).SetErrorCode(errorx.ErrorUserNotFoundCodeValue)

	repo := &UserRepo{
		DB:          params.DB,
		GenericRepo: genericRepo,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.Info("UserRepo initialized")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if params.Logger != nil {
				logger.Info("UserRepo stopped")
			}
			return nil
		},
	})

	return repo
}

// GetByID 根据ID获取用户
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	// 使用仓库接口
	user, err := r.GenericRepo.Get(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	// 返回用户
	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	// 使用仓库接口的查询
	opts := &QueryOptions{
		Condition: "username = ?",
		Args:      []any{username},
	}
	user, err := r.GenericRepo.Get(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	// 返回用户
	return user, nil
}

// Create 创建用户
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	// 使用仓库接口
	return r.GenericRepo.Create(ctx, user)
}

// Update 更新用户
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	// 使用仓库接口
	return r.GenericRepo.Update(ctx, user)
}

// UpdateFields 更新用户字段
func (r *UserRepo) UpdateFields(ctx context.Context, id int64, fields map[string]any) error {
	// 使用仓库接口
	return r.GenericRepo.UpdateFields(ctx, id, fields)
}

// Delete 删除用户
func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	// 使用仓库接口
	return r.GenericRepo.Delete(ctx, id)
}

// IsExist 判断用户是否存在
func (r *UserRepo) IsExist(ctx context.Context, username string) (bool, error) {
	// 使用泛型仓库的查询选项
	opts := &QueryOptions{
		Condition: "username = ?",
		Args:      []any{username},
	}

	// 使用泛型仓库的Count方法
	count, err := r.GenericRepo.Count(ctx, opts)
	if err != nil {
		return false, errorx.WrapError(err, fmt.Sprintf("检查用户是否存在失败: username=%s", username))
	}
	return count > 0, nil
}

// UpdateAvatar 更新用户头像
func (r *UserRepo) UpdateAvatar(ctx context.Context, userID int64, fileID uint) error {
	// 使用泛型仓库的事务支持
	return r.GenericRepo.Transaction(ctx, func(tx *gorm.DB) error {
		// 创建事务中的仓库
		txRepo := r.GenericRepo.WithTx(tx)

		// 查找用户
		user, err := txRepo.GetByID(ctx, userID)
		if err != nil {
			return errorx.WrapError(err, fmt.Sprintf("查询用户失败: id=%d", userID))
		}

		// 更新用户头像
		user.AvatarFileID = fileID

		// 保存用户
		return txRepo.Update(ctx, user)
	})
}

// WithTx 使用事务
func (r *UserRepo) WithTx(tx *gorm.DB) *UserRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.User](tx).SetErrorCode(errorx.ErrorUserNotFoundCodeValue)

	// 创建新的仓库实例，使用事务
	return &UserRepo{
		DB:          tx,
		GenericRepo: genericRepo,
	}
}

// GetUserRoles 获取用户角色
func (r *UserRepo) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	// 使用泛型仓库的Get方法，并预加载Roles关联
	opts := &QueryOptions{
		Preloads: []string{"Roles"},
	}

	// 查询用户及其角色
	user, err := r.GenericRepo.Get(ctx, userID, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询用户角色失败: id=%d", userID))
	}

	// 提取角色编码
	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Code)
	}

	// 如果没有角色，返回默认角色
	if len(roles) == 0 {
		return []string{"user"}, nil
	}

	return roles, nil
}

// List 获取用户列表
func (r *UserRepo) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	// 标准化分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 使用泛型仓库的GetPage方法
	users, total, err := r.GenericRepo.GetPage(ctx, page, pageSize, "")
	if err != nil {
		return nil, 0, errorx.WrapError(err, "获取用户列表失败")
	}

	// 转换为指针切片
	userPtrs := make([]*model.User, len(users))
	for i := range users {
		userPtrs[i] = &users[i]
	}

	return userPtrs, total, nil
}

// FindByStatus 根据状态查询用户
func (r *UserRepo) FindByStatus(ctx context.Context, status int, page, pageSize int) ([]*model.User, int64, error) {
	// 标准化分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 根据状态值设置查询条件
	var condition string
	var args []any

	if status == model.UserStatusActive {
		condition = "enabled = ?"
		args = []any{true}
	} else if status == model.UserStatusDisabled {
		condition = "enabled = ?"
		args = []any{false}
	}

	// 使用泛型仓库的GetPage方法
	users, total, err := r.GenericRepo.GetPage(ctx, page, pageSize, condition, args...)
	if err != nil {
		return nil, 0, errorx.WrapError(err, fmt.Sprintf("获取状态为 %d 的用户列表失败", status))
	}

	// 转换为指针切片
	userPtrs := make([]*model.User, len(users))
	for i := range users {
		userPtrs[i] = &users[i]
	}

	return userPtrs, total, nil
}
