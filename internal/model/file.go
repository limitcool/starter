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
	URL            string    `json:"url" gorm:"-"` // 计算字段，不存储到数据库
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
	DB          *gorm.DB
	GenericRepo *GenericRepo[File]
}

// NewFileRepo 创建文件仓库
func NewFileRepo(db *gorm.DB) *FileRepo {
	genericRepo := NewGenericRepo[File](db)
	genericRepo.ErrorCode = errorx.ErrorFileNotFoundCode

	return &FileRepo{
		DB:          db,
		GenericRepo: genericRepo,
	}
}

// Create 创建文件记录
func (r *FileRepo) Create(ctx context.Context, file *File) error {
	return r.GenericRepo.Create(ctx, file)
}

// GetByID 根据ID获取文件
func (r *FileRepo) GetByID(ctx context.Context, id uint) (*File, error) {
	file, err := r.GenericRepo.Get(ctx, id, nil)
	if err != nil {
		return nil, errorx.WrapError(err, "查询文件失败")
	}

	// 设置URL字段
	if file.Path != "" {
		file.URL = r.buildFileURL(file.Path)
	}

	return file, nil
}

// buildFileURL 构建文件URL
func (r *FileRepo) buildFileURL(path string) string {
	// 这里可以根据实际情况构建URL，例如添加域名前缀等
	// 简单示例：
	return "/uploads/" + path
}

// Update 更新文件记录
func (r *FileRepo) Update(ctx context.Context, file *File) error {
	return r.GenericRepo.Update(ctx, file)
}

// Delete 删除文件记录
func (r *FileRepo) Delete(ctx context.Context, id uint) error {
	return r.GenericRepo.Delete(ctx, id)
}

// UpdateFileUsage 更新文件用途
func (r *FileRepo) UpdateFileUsage(ctx context.Context, file *File, usage string) error {
	file.Usage = usage
	return r.GenericRepo.Update(ctx, file)
}

// ListByUser 获取用户的文件列表
func (r *FileRepo) ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]File, error) {
	// 使用统一的List方法
	opts := &QueryOptions{
		Condition: "uploaded_by = ? AND uploaded_by_type = ?",
		Args:      []any{userID, 2},
	}
	files, err := r.GenericRepo.List(ctx, page, pageSize, opts)
	if err != nil {
		return nil, errorx.WrapError(err, "查询用户文件列表失败")
	}

	// 为所有文件设置URL
	for i := range files {
		if files[i].Path != "" {
			files[i].URL = r.buildFileURL(files[i].Path)
		}
	}

	return files, nil
}

// CountByUser 获取用户的文件总数
func (r *FileRepo) CountByUser(ctx context.Context, userID int64) (int64, error) {
	count, err := r.GenericRepo.Count(ctx, &QueryOptions{
		Condition: "uploaded_by = ? AND uploaded_by_type = ?",
		Args:      []any{userID, 2},
	})
	if err != nil {
		return 0, errorx.WrapError(err, "查询用户文件总数失败")
	}
	return count, nil
}

// ListFiles 获取文件列表，支持多种查询条件和预加载
func (r *FileRepo) ListFiles(ctx context.Context, page, pageSize int, fileType, usage string, preloads []string) ([]File, int64, error) {
	var conditions []string
	var args []any

	// 添加查询条件
	if fileType != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, fileType)
	}

	if usage != "" {
		conditions = append(conditions, "usage = ?")
		args = append(args, usage)
	}

	// 构建查询选项
	opts := &QueryOptions{
		Preloads: preloads,
	}

	// 如果有条件，设置条件
	if len(conditions) > 0 {
		opts.Condition = conditions[0]
		opts.Args = args[:1]

		// 如果有多个条件，使用AND连接
		for i := 1; i < len(conditions); i++ {
			opts.Condition += " AND " + conditions[i]
			opts.Args = append(opts.Args, args[i])
		}
	}

	// 获取文件列表
	files, err := r.GenericRepo.List(ctx, page, pageSize, opts)
	if err != nil {
		return nil, 0, errorx.WrapError(err, "查询文件列表失败")
	}

	// 为所有文件设置URL
	for i := range files {
		if files[i].Path != "" {
			files[i].URL = r.buildFileURL(files[i].Path)
		}
	}

	// 获取总数
	total, err := r.GenericRepo.Count(ctx, opts)
	if err != nil {
		return nil, 0, errorx.WrapError(err, "查询文件总数失败")
	}

	return files, total, nil
}
