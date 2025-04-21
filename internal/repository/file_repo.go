package repository

import (
	"context"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// FileRepo 文件仓库
type FileRepo struct {
	DB          *gorm.DB
	genericRepo *GenericRepo[model.File] // 泛型仓库
}

// NewFileRepo 创建文件仓库
func NewFileRepo(db *gorm.DB) *FileRepo {
	genericRepo := NewGenericRepo[model.File](db)
	genericRepo.SetErrorCode(errorx.ErrorNotFoundCode) // 设置错误码

	return &FileRepo{
		DB:          db,
		genericRepo: genericRepo,
	}
}

// Create 创建文件记录
func (r *FileRepo) Create(ctx context.Context, file *model.File) error {
	// 使用泛型仓库
	return r.genericRepo.Create(ctx, file)
}

// GetByID 根据ID获取文件
func (r *FileRepo) GetByID(ctx context.Context, id string) (*model.File, error) {
	// 使用泛型仓库
	return r.genericRepo.GetByID(ctx, id)
}

// Delete 删除文件记录
func (r *FileRepo) Delete(ctx context.Context, id string) error {
	// 使用泛型仓库
	return r.genericRepo.Delete(ctx, id)
}

// Update 更新文件记录
func (r *FileRepo) Update(ctx context.Context, file *model.File) error {
	// 使用泛型仓库
	return r.genericRepo.Update(ctx, file)
}

// UpdateFileUsage 更新文件用途
func (r *FileRepo) UpdateFileUsage(ctx context.Context, file *model.File, usage string) error {
	file.Usage = usage
	// 使用泛型仓库
	return r.genericRepo.Update(ctx, file)
}

// UpdateUserAvatar 更新用户头像
func (r *FileRepo) UpdateUserAvatar(ctx context.Context, userID int64, fileID uint) error {
	// 这里应该实现更新用户头像的逻辑
	// 当前只是一个占位实现
	return nil
}

// UpdateAdminUserAvatar 更新管理员用户头像
func (r *FileRepo) UpdateAdminUserAvatar(ctx context.Context, userID int64, fileID uint) error {
	// 这里应该实现更新管理员用户头像的逻辑
	// 当前只是一个占位实现
	return nil
}
