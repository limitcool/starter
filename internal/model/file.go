package model

import (
	"context"
	"time"

	"github.com/limitcool/starter/internal/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 使用辅助方法获取集合
func getFileCollection() *mongo.Collection {
	return mongodb.Collection("file")
}

// 文件类型枚举
const (
	FileTypeImage    = "image"    // 图片
	FileTypeAvatar   = "avatar"   // 头像
	FileTypeDocument = "document" // 文档
	FileTypeVideo    = "video"    // 视频
	FileTypeAudio    = "audio"    // 音频
	FileTypeOther    = "other"    // 其他
)

// File 文件模型
type File struct {
	BaseModel

	Name         string    `json:"name" gorm:"size:255;comment:文件名称"`
	OriginalName string    `json:"original_name" gorm:"size:255;comment:原始文件名"`
	Path         string    `json:"path" gorm:"size:500;comment:存储路径"`
	URL          string    `json:"url" gorm:"size:500;comment:访问URL"`
	Type         string    `json:"type" gorm:"size:50;comment:文件类型"`
	Size         int64     `json:"size" gorm:"comment:文件大小(字节)"`
	MimeType     string    `json:"mime_type" gorm:"size:100;comment:MIME类型"`
	Extension    string    `json:"extension" gorm:"size:20;comment:扩展名"`
	StorageType  string    `json:"storage_type" gorm:"size:20;comment:存储类型(local/s3/oss)"`
	UploadedBy   int64     `json:"uploaded_by" gorm:"type:bigint;comment:上传者ID"`
	UploadedAt   time.Time `json:"uploaded_at" gorm:"comment:上传时间"`
	Status       int       `json:"status" gorm:"comment:状态(1:正常,0:禁用,-1:删除)"`
}

func (File) TableName() string {
	return "file"
}

func (File) Registry() {
	var ctx = context.Background()
	coll := getFileCollection()
	if coll == nil {
		return
	}
	coll.FindOne(ctx, bson.M{"name": "file"})
	coll.InsertOne(ctx, bson.M{"name": "file"})
}
