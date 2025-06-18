package filestore

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage 本地存储实现
type LocalStorage struct {
	basePath string // 存储根路径，如 "./uploads"
	baseURL  string // 访问基础URL，如 "http://localhost:8081/uploads"
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// GetUploadURL 获取上传URL（本地存储返回应用服务器的统一上传接口）
func (l *LocalStorage) GetUploadURL(ctx context.Context, filePath string, contentType string, isPublic bool) (string, string, error) {
	// 本地存储返回应用服务器的统一上传接口
	fullPath := l.BuildFullPath(filePath, isPublic)
	uploadURL := fmt.Sprintf("/api/v1/upload/file?path=%s&public=%t", fullPath, isPublic)
	return uploadURL, "POST", nil
}

// GetDownloadURL 获取下载URL（本地存储返回HTTP静态文件URL）
func (l *LocalStorage) GetDownloadURL(ctx context.Context, filePath string, isPublic bool) (string, error) {
	fullPath := l.BuildFullPath(filePath, isPublic)
	downloadURL := fmt.Sprintf("%s/%s", strings.TrimRight(l.baseURL, "/"), strings.TrimLeft(fullPath, "/"))
	return downloadURL, nil
}

// UploadFile 直接上传文件到本地存储
func (l *LocalStorage) UploadFile(ctx context.Context, filePath string, reader io.Reader, isPublic bool) error {
	fullPath := l.BuildFullPath(filePath, isPublic)
	absolutePath := filepath.Join(l.basePath, fullPath)

	// 确保目录存在
	dir := filepath.Dir(absolutePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	file, err := os.Create(absolutePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 复制内容
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// FileExists 检查文件是否存在
func (l *LocalStorage) FileExists(ctx context.Context, filePath string, isPublic bool) (bool, error) {
	fullPath := l.BuildFullPath(filePath, isPublic)
	absolutePath := filepath.Join(l.basePath, fullPath)

	_, err := os.Stat(absolutePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// DeleteFile 删除文件
func (l *LocalStorage) DeleteFile(ctx context.Context, filePath string, isPublic bool) error {
	fullPath := l.BuildFullPath(filePath, isPublic)
	absolutePath := filepath.Join(l.basePath, fullPath)
	return os.Remove(absolutePath)
}

// GetStorageType 获取存储类型
func (l *LocalStorage) GetStorageType() string {
	return "local"
}

// BuildFullPath 构建完整路径（包含public/private前缀）
func (l *LocalStorage) BuildFullPath(filePath string, isPublic bool) string {
	if isPublic {
		return fmt.Sprintf("public/%s", filePath)
	}
	return fmt.Sprintf("private/%s", filePath)
}
