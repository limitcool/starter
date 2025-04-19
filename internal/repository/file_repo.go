package repository

import (
	"errors"
	"fmt"

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
	err := r.DB.Create(file).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("创建文件记录失败: %s", file.Name))
	}
	return nil
}

// GetByID 根据ID获取文件
func (r *FileRepo) GetByID(id string) (*model.File, error) {
	var file model.File
	err := r.DB.First(&file, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		notFoundErr := errorx.Errorf(errorx.ErrNotFound, "文件ID %s 不存在", id)
		return nil, errorx.WrapError(notFoundErr, "")
	}
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询文件失败: id=%s", id))
	}
	return &file, nil
}

// Delete 删除文件记录
func (r *FileRepo) Delete(id string) error {
	err := r.DB.Delete(&model.File{}, id).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除文件记录失败: id=%s", id))
	}
	return nil
}

// Update 更新文件记录
func (r *FileRepo) Update(file *model.File) error {
	err := r.DB.Save(file).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新文件记录失败: id=%s, name=%s", file.ID, file.Name))
	}
	return nil
}

// UpdateFileUsage 更新文件用途
func (r *FileRepo) UpdateFileUsage(file *model.File, usage string) error {
	file.Usage = usage
	err := r.DB.Save(file).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新文件用途失败: id=%s, usage=%s", file.ID, usage))
	}
	return nil
}
