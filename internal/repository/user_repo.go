package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// UserRepo 用户仓库
// 提供用户相关的数据库操作
type UserRepo struct {
	DB          *gorm.DB
	GenericRepo *GenericRepo[model.User] // 泛型仓库
}

// NewUserRepo 创建用户仓库
func NewUserRepo(params RepoParams) *UserRepo {
	genericRepo := NewGenericRepo[model.User](params.DB)
	genericRepo.SetErrorCode(errorx.ErrorUserNotFoundCode) // 设置错误码

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
	// 使用泛型仓库
	user, err := r.GenericRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 返回用户
	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	// 使用泛型仓库的高级查询
	user, err := r.GenericRepo.FindByField(ctx, "username", username)
	if err != nil {
		return nil, err
	}

	// 返回用户
	return user, nil
}

// Create 创建用户
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	// 使用泛型仓库
	return r.GenericRepo.Create(ctx, user)
}

// Update 更新用户
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	// 使用泛型仓库
	return r.GenericRepo.Update(ctx, user)
}

// UpdateFields 更新用户字段
func (r *UserRepo) UpdateFields(ctx context.Context, id int64, fields map[string]any) error {
	// 使用泛型仓库
	return r.GenericRepo.UpdateFields(ctx, id, fields)
}

// Delete 删除用户
func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	// 使用泛型仓库
	return r.GenericRepo.Delete(ctx, id)
}

// IsExist 判断用户是否存在
func (r *UserRepo) IsExist(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, errorx.WrapError(err, fmt.Sprintf("检查用户是否存在失败: username=%s", username))
	}
	return count > 0, nil
}

// UpdateAvatar 更新用户头像
func (r *UserRepo) UpdateAvatar(ctx context.Context, userID int64, fileID uint) error {
	// 开始事务
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 查找用户
	user := model.User{}
	if err := tx.WithContext(ctx).First(&user, userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFoundErr := errorx.Errorf(errorx.ErrUserNotFound, "用户ID %d 不存在", userID)
			return errorx.WrapError(notFoundErr, "")
		}
		return errorx.WrapError(err, fmt.Sprintf("查询用户失败: id=%d", userID))
	}

	// 更新用户头像
	user.AvatarFileID = fileID
	if err := tx.WithContext(ctx).Save(&user).Error; err != nil {
		tx.Rollback()
		return errorx.WrapError(err, fmt.Sprintf("更新用户头像失败: id=%d, fileID=%d", userID, fileID))
	}

	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		return errorx.WrapError(err, fmt.Sprintf("提交事务失败: 更新用户头像, userID=%d, fileID=%d", userID, fileID))
	}
	return nil
}

// WithTx 使用事务
func (r *UserRepo) WithTx(tx *gorm.DB) *UserRepo {
	genericRepo := NewGenericRepo[model.User](tx)
	genericRepo.SetErrorCode(errorx.ErrorUserNotFoundCode)

	// 创建新的仓库实例，使用事务
	return &UserRepo{
		DB:          tx,
		GenericRepo: genericRepo,
	}
}

// GetUserRoles 获取用户角色
func (r *UserRepo) GetUserRoles(ctx context.Context, userID int64, isAdmin bool, userMode string) ([]string, error) {
	// 如果是简单模式
	if userMode == string(enum.UserModeSimple) {
		// 如果是管理员
		if isAdmin {
			return []string{"admin"}, nil
		}
		// 普通用户
		return []string{"user"}, nil
	}

	// 分离模式，从数据库查询角色
	var roles []string
	user := &model.User{}

	// 查询用户及其角色
	err := r.DB.WithContext(ctx).Preload("Roles").First(user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFoundErr := errorx.Errorf(errorx.ErrUserNotFound, "用户ID %d 不存在", userID)
			return nil, errorx.WrapError(notFoundErr, "")
		}
		return nil, errorx.WrapError(err, fmt.Sprintf("查询用户角色失败: id=%d", userID))
	}

	// 提取角色编码
	for _, role := range user.Roles {
		roles = append(roles, role.Code)
	}

	// 如果没有角色，返回默认角色
	if len(roles) == 0 {
		return []string{"user"}, nil
	}

	return roles, nil
}
