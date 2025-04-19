package repository

import (
	"errors"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// 已移除 UserRepository 接口定义

// UserRepo 用户仓库
type UserRepo struct {
	DB *gorm.DB
}

// NewUserRepo 创建用户仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

// GetByID 根据ID获取用户
func (r *UserRepo) GetByID(id int64) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		notFoundErr := errorx.Errorf(errorx.ErrUserNotFound, "用户ID %d 不存在", id)
		return nil, errorx.WrapError(notFoundErr, "")
	}
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询用户失败: id=%d", id))
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		notFoundErr := errorx.Errorf(errorx.ErrUserNotFound, "用户名 %s 不存在", username)
		return nil, errorx.WrapError(notFoundErr, "")
	}
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询用户失败: username=%s", username))
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepo) Create(user *model.User) error {
	err := r.DB.Create(user).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("创建用户失败: username=%s", user.Username))
	}
	return nil
}

// Update 更新用户
func (r *UserRepo) Update(user *model.User) error {
	err := r.DB.Save(user).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新用户失败: id=%d, username=%s", user.ID, user.Username))
	}
	return nil
}

// UpdateFields 更新用户字段
func (r *UserRepo) UpdateFields(id int64, fields map[string]any) error {
	err := r.DB.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新用户字段失败: id=%d", id))
	}
	return nil
}

// Delete 删除用户
func (r *UserRepo) Delete(id int64) error {
	err := r.DB.Delete(&model.User{}, id).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除用户失败: id=%d", id))
	}
	return nil
}

// IsExist 判断用户是否存在
func (r *UserRepo) IsExist(username string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, errorx.WrapError(err, fmt.Sprintf("检查用户是否存在失败: username=%s", username))
	}
	return count > 0, nil
}

// UpdateAvatar 更新用户头像
func (r *UserRepo) UpdateAvatar(userID int64, fileID uint) error {
	// 开始事务
	tx := r.DB.Begin()
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
	return &UserRepo{DB: tx}
}
