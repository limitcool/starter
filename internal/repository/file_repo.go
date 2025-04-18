package repository

import (
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// FileRepo 文件仓库
type FileRepo struct {
	DB *gorm.DB
}

// NewFileRepo 创建文件仓库
func NewFileRepo(db *gorm.DB) *FileRepo {
	return &FileRepo{DB: db}
}

// Create 创建文件记录
func (r *FileRepo) Create(file *model.File) error {
	return r.DB.Create(file).Error
}

// GetByID 根据ID获取文件
func (r *FileRepo) GetByID(id string) (*model.File, error) {
	var file model.File
	if err := r.DB.First(&file, id).Error; err != nil {
		return nil, errorx.ErrNotFound.WithError(err)
	}
	return &file, nil
}

// Delete 删除文件记录
func (r *FileRepo) Delete(id string) error {
	return r.DB.Delete(&model.File{}, id).Error
}

// Update 更新文件记录
func (r *FileRepo) Update(file *model.File) error {
	return r.DB.Save(file).Error
}

// UpdateUserAvatar 更新用户头像
func (r *FileRepo) UpdateUserAvatar(userID int64, fileID uint) error {
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

// UpdateFileUsage 更新文件用途
func (r *FileRepo) UpdateFileUsage(file *model.File, usage string) error {
	file.Usage = usage
	return r.DB.Save(file).Error
}

// UpdateSysUserAvatar 更新系统用户头像
func (r *FileRepo) UpdateSysUserAvatar(userID int64, fileID uint) error {
	// 开始事务
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 查找用户
	sysUser := model.SysUser{}
	if err := tx.First(&sysUser, userID).Error; err != nil {
		tx.Rollback()
		return errorx.ErrNotFound.WithError(err)
	}

	// 更新用户头像
	sysUser.AvatarFileID = fileID
	if err := tx.Save(&sysUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
