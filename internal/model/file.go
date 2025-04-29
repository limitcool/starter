package model

import (
	"context"
	"time"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// 文件类型枚举
const (
	FileTypeImage    = "image"    // 图片
	FileTypeDocument = "document" // 文档
	FileTypeVideo    = "video"    // 视频
	FileTypeAudio    = "audio"    // 音频
	FileTypeOther    = "other"    // 其他
)

// 文件用途枚举
const (
	FileUsageAvatar  = "avatar"  // 头像
	FileUsageAttach  = "attach"  // 附件
	FileUsageBanner  = "banner"  // 横幅
	FileUsageCover   = "cover"   // 封面
	FileUsageGallery = "gallery" // 相册
	FileUsageGeneral = "general" // 通用
)

// File 文件模型
type File struct {
	BaseModel

	Name           string    `json:"name" gorm:"size:255;comment:文件名称"`
	OriginalName   string    `json:"original_name" gorm:"size:255;comment:原始文件名"`
	Path           string    `json:"path" gorm:"size:500;comment:存储路径"`
	URL            string    `json:"url" gorm:"size:500;comment:访问URL"`
	Type           string    `json:"type" gorm:"size:50;comment:文件类型"`
	Usage          string    `json:"usage" gorm:"size:50;default:'general';comment:文件用途"`
	Size           int64     `json:"size" gorm:"comment:文件大小(字节)"`
	MimeType       string    `json:"mime_type" gorm:"size:100;comment:MIME类型"`
	Extension      string    `json:"extension" gorm:"size:20;comment:扩展名"`
	StorageType    string    `json:"storage_type" gorm:"size:20;comment:存储类型(local/s3/oss)"`
	UploadedBy     int64     `json:"uploaded_by" gorm:"type:bigint;comment:上传者ID"`
	UploadedByType uint8     `json:"uploaded_by_type" gorm:"size:20;default:1;comment:上传者类型(1:系统用户,2:普通用户)"`
	UploadedAt     time.Time `json:"uploaded_at" gorm:"comment:上传时间"`
	Status         int       `json:"status" gorm:"comment:状态(1:正常,0:禁用,-1:删除)"`
	IsPublic       bool      `json:"is_public" gorm:"default:false;comment:是否公开访问"`
}

func (File) TableName() string {
	return "file"
}

// FileRepo 文件仓库
type FileRepo struct {
	DB *gorm.DB
}

// NewFileRepo 创建文件仓库
func NewFileRepo(db *gorm.DB) *FileRepo {
	return &FileRepo{
		DB: db,
	}
}

// Create 创建文件记录
func (r *FileRepo) Create(ctx context.Context, file *File) error {
	return r.DB.WithContext(ctx).Create(file).Error
}

// GetByID 根据ID获取文件
func (r *FileRepo) GetByID(ctx context.Context, id uint) (*File, error) {
	var file File
	if err := r.DB.WithContext(ctx).First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrFileNotFound
		}
		return nil, errorx.WrapError(err, "查询文件失败")
	}
	return &file, nil
}

// Update 更新文件记录
func (r *FileRepo) Update(ctx context.Context, file *File) error {
	return r.DB.WithContext(ctx).Save(file).Error
}

// Delete 删除文件记录
func (r *FileRepo) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&File{}, id).Error
}

// UpdateFileUsage 更新文件用途
func (r *FileRepo) UpdateFileUsage(ctx context.Context, file *File, usage string) error {
	file.Usage = usage
	return r.DB.WithContext(ctx).Save(file).Error
}

// ListByUser 获取用户的文件列表
func (r *FileRepo) ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]File, error) {
	var files []File
	offset := (page - 1) * pageSize
	if err := r.DB.WithContext(ctx).Where("uploaded_by = ? AND uploaded_by_type = ?", userID, 2).
		Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, errorx.WrapError(err, "查询用户文件列表失败")
	}
	return files, nil
}

// CountByUser 获取用户的文件总数
func (r *FileRepo) CountByUser(ctx context.Context, userID int64) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&File{}).
		Where("uploaded_by = ? AND uploaded_by_type = ?", userID, 2).Count(&count).Error; err != nil {
		return 0, errorx.WrapError(err, "查询用户文件总数失败")
	}
	return count, nil
}
