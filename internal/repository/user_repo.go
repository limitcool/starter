package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id int64) (*model.User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*model.User, error)

	// Create 创建用户
	Create(ctx context.Context, user *model.User) error

	// Update 更新用户
	Update(ctx context.Context, user *model.User) error

	// UpdateFields 更新用户字段
	UpdateFields(ctx context.Context, id int64, fields map[string]any) error

	// Delete 删除用户
	Delete(ctx context.Context, id int64) error

	// IsExist 判断用户是否存在
	IsExist(ctx context.Context, username string) (bool, error)

	// UpdateAvatar 更新用户头像
	UpdateAvatar(ctx context.Context, userID int64, fileID uint) error

	// WithTx 使用事务
	WithTx(tx *gorm.DB) *UserRepo
}

// UserRepo 用户仓库
type UserRepo struct {
	DB          *gorm.DB
	GenericRepo *GenericRepo[model.User] // 泛型仓库
}

// NewUserRepo 创建用户仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
	genericRepo := NewGenericRepo[model.User](db)
	genericRepo.SetErrorCode(errorx.ErrorUserNotFoundCode) // 设置错误码

	return &UserRepo{
		DB:          db,
		GenericRepo: genericRepo,
	}
}

// GetByID 根据ID获取用户
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	// 使用泛型仓库
	return r.GenericRepo.GetByID(ctx, id)
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	// 使用泛型仓库的高级查询
	return r.GenericRepo.FindByField(ctx, "username", username)
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
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFoundErr := errorx.Errorf(errorx.ErrUserNotFound, "用户ID %d 不存在", userID)
			return errorx.WrapError(notFoundErr, "")
		}
		return errorx.WrapError(err, fmt.Sprintf("查询用户失败: id=%d", userID))
	}

	// 更新用户头像
	user.AvatarFileID = fileID
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return errorx.WrapError(err, fmt.Sprintf("更新用户头像失败: id=%d, fileID=%d", userID, fileID))
	}

	if err := tx.Commit().Error; err != nil {
		return errorx.WrapError(err, fmt.Sprintf("提交事务失败: 更新用户头像, userID=%d, fileID=%d", userID, fileID))
	}
	return nil
}

// WithTx 使用事务
func (r *UserRepo) WithTx(tx *gorm.DB) *UserRepo {
	genericRepo := NewGenericRepo[model.User](tx)
	genericRepo.SetErrorCode(errorx.ErrorUserNotFoundCode)

	return &UserRepo{
		DB:          tx,
		GenericRepo: genericRepo,
	}
}
