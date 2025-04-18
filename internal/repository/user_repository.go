package repository

import (
	"errors"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// 已移除 UserRepository 接口定义

// GormUserRepository 基于Gorm的用户仓库实现
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// GetByID 根据ID获取用户
func (r *GormUserRepository) GetByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	return &user, err
}

// GetByUsername 根据用户名获取用户
func (r *GormUserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	return &user, err
}

// Create 创建用户
func (r *GormUserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *GormUserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// UpdateFields 更新用户字段
func (r *GormUserRepository) UpdateFields(id int64, fields map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除用户
func (r *GormUserRepository) Delete(id int64) error {
	return r.db.Delete(&model.User{}, id).Error
}

// IsExist 判断用户是否存在
func (r *GormUserRepository) IsExist(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// UpdateAvatar 更新用户头像
func (r *GormUserRepository) UpdateAvatar(userID int64, fileID uint) error {
	// 开始事务
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
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
func (r *GormUserRepository) WithTx(tx *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: tx}
}
