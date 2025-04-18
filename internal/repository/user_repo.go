package repository

import (
	"errors"

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
		return nil, errorx.ErrUserNotFound
	}
	return &user, err
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	return &user, err
}

// Create 创建用户
func (r *UserRepo) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

// Update 更新用户
func (r *UserRepo) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

// UpdateFields 更新用户字段
func (r *UserRepo) UpdateFields(id int64, fields map[string]any) error {
	return r.DB.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除用户
func (r *UserRepo) Delete(id int64) error {
	return r.DB.Delete(&model.User{}, id).Error
}

// IsExist 判断用户是否存在
func (r *UserRepo) IsExist(username string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
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
		return errorx.ErrNotFound.WithError(err)
	}

	// 更新用户头像
	user.AvatarFileID = fileID
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// WithTx 使用事务
func (r *UserRepo) WithTx(tx *gorm.DB) *UserRepo {
	return &UserRepo{DB: tx}
}
