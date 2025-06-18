package filestore

import (
	"context"
	"io"
)

// FileStorage 文件存储接口
type FileStorage interface {
	// GetUploadURL 获取上传URL
	// filePath: 文件路径（不包含public/private前缀）
	// contentType: 文件MIME类型
	// isPublic: 是否公开文件
	// returns: 上传URL, 上传方法(PUT/POST), error
	GetUploadURL(ctx context.Context, filePath string, contentType string, isPublic bool) (string, string, error)

	// GetDownloadURL 获取下载URL
	// filePath: 文件路径（不包含public/private前缀）
	// isPublic: 是否公开文件
	// returns: 下载URL, error
	GetDownloadURL(ctx context.Context, filePath string, isPublic bool) (string, error)

	// UploadFile 直接上传文件（仅用于本地存储或特殊情况）
	// filePath: 文件路径（不包含public/private前缀）
	// reader: 文件内容
	// isPublic: 是否公开文件
	UploadFile(ctx context.Context, filePath string, reader io.Reader, isPublic bool) error

	// FileExists 检查文件是否存在
	// filePath: 文件路径（不包含public/private前缀）
	// isPublic: 是否公开文件
	FileExists(ctx context.Context, filePath string, isPublic bool) (bool, error)

	// DeleteFile 删除文件
	// filePath: 文件路径（不包含public/private前缀）
	// isPublic: 是否公开文件
	DeleteFile(ctx context.Context, filePath string, isPublic bool) error

	// GetStorageType 获取存储类型
	GetStorageType() string

	// BuildFullPath 构建完整路径（包含public/private前缀）
	BuildFullPath(filePath string, isPublic bool) string
}

// UploadResponse 上传响应
type UploadResponse struct {
	UploadURL string `json:"upload_url"`
	Method    string `json:"method"`    // PUT 或 POST
	Headers   map[string]string `json:"headers,omitempty"`
}

// DownloadResponse 下载响应
type DownloadResponse struct {
	DownloadURL string `json:"download_url"`
	IsPublic    bool   `json:"is_public"`
	ExpiresAt   *int64 `json:"expires_at,omitempty"` // 过期时间戳，公开文件为nil
}
