package model

import (
	"time"

	"github.com/limitcool/starter/internal/pkg/enum"
)

// 文件类型枚举
const (
	FileTypeImage    = "image"    // 图片
	FileTypeAvatar   = "avatar"   // 头像
	FileTypeDocument = "document" // 文档
	FileTypeVideo    = "video"    // 视频
	FileTypeAudio    = "audio"    // 音频
	FileTypeOther    = "other"    // 其他
	FileTypePublic   = "public"   // 公共文件，所有用户可访问
)

// File 文件模型
type File struct {
	BaseModel

	Name           string        `json:"name" gorm:"size:255;comment:文件名称"`
	OriginalName   string        `json:"original_name" gorm:"size:255;comment:原始文件名"`
	Path           string        `json:"path" gorm:"size:500;comment:存储路径"`
	URL            string        `json:"url" gorm:"size:500;comment:访问URL"`
	Type           string        `json:"type" gorm:"size:50;comment:文件类型"`
	Size           int64         `json:"size" gorm:"comment:文件大小(字节)"`
	MimeType       string        `json:"mime_type" gorm:"size:100;comment:MIME类型"`
	Extension      string        `json:"extension" gorm:"size:20;comment:扩展名"`
	StorageType    string        `json:"storage_type" gorm:"size:20;comment:存储类型(local/s3/oss)"`
	UploadedBy     int64         `json:"uploaded_by" gorm:"type:bigint;comment:上传者ID"`
	UploadedByType enum.UserType `json:"uploaded_by_type" gorm:"default:1;comment:上传者类型(1:系统用户,2:普通用户)"`
	UploadedAt     time.Time     `json:"uploaded_at" gorm:"comment:上传时间"`
	Status         int           `json:"status" gorm:"comment:状态(1:正常,0:禁用,-1:删除)"`
}

func (File) TableName() string {
	return "file"
}
