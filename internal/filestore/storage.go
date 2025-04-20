package filestore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/types"
	"github.com/qor/oss"
	"github.com/qor/oss/filesystem"
	"github.com/qor/oss/s3"
)

// Config 存储配置
type Config struct {
	Type      types.StorageType `json:"type"`       // 存储类型
	Path      string            `json:"path"`       // 本地存储路径
	URL       string            `json:"url"`        // 访问URL
	AccessKey string            `json:"access_key"` // 访问密钥
	SecretKey string            `json:"secret_key"` // 访问密钥
	Region    string            `json:"region"`     // 区域
	Bucket    string            `json:"bucket"`     // 桶名称
	Endpoint  string            `json:"endpoint"`   // 端点
}

// Storage 存储服务
type Storage struct {
	Config Config
	oss    oss.StorageInterface
}

// New 创建存储服务
func New(config Config) (*Storage, error) {
	var ossStorage oss.StorageInterface
	var err error

	switch config.Type {
	case types.StorageTypeLocal:
		// 确保本地存储路径存在
		err = os.MkdirAll(config.Path, os.ModePerm)
		if err != nil {
			return nil, errorx.WrapError(err, fmt.Sprintf("创建存储路径失败: %s", config.Path))
		}
		ossStorage = filesystem.New(config.Path)
	case types.StorageTypeS3:
		// S3配置检查
		if config.AccessKey == "" || config.SecretKey == "" || config.Bucket == "" {
			return nil, errorx.Errorf(errorx.ErrFileStorage, "S3配置不完整")
		}
		ossStorage = s3.New(&s3.Config{
			AccessID:  config.AccessKey,
			AccessKey: config.SecretKey,
			Region:    config.Region,
			Bucket:    config.Bucket,
			Endpoint:  config.Endpoint,
		})
	default:
		return nil, errorx.Errorf(errorx.ErrFileStorage, "不支持的存储类型: %s", config.Type)
	}

	return &Storage{
		Config: config,
		oss:    ossStorage,
	}, nil
}

// Put 上传文件
func (s *Storage) Put(path string, reader io.Reader) error {
	_, err := s.oss.Put(path, reader)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("上传文件失败: %s", path))
	}
	return nil
}

// Get 获取文件
func (s *Storage) Get(path string) (*os.File, error) {
	file, err := s.oss.Get(path)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取文件失败: %s", path))
	}
	return file, nil
}

// GetStream 获取文件流
func (s *Storage) GetStream(path string) (io.ReadCloser, error) {
	stream, err := s.oss.GetStream(path)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取文件流失败: %s", path))
	}
	return stream, nil
}

// Delete 删除文件
func (s *Storage) Delete(path string) error {
	err := s.oss.Delete(path)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除文件失败: %s", path))
	}
	return nil
}

// List 列出目录下的文件
func (s *Storage) List(path string) ([]*oss.Object, error) {
	objects, err := s.oss.List(path)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取文件列表失败: %s", path))
	}
	return objects, nil
}

// GetURL 获取文件URL
func (s *Storage) GetURL(path string) (string, error) {
	// 确保路径使用正斜杠
	normalizedPath := strings.ReplaceAll(path, "\\", "/")

	// 本地存储特殊处理
	if s.Config.Type == types.StorageTypeLocal && s.Config.URL != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(s.Config.URL, "/"), normalizedPath), nil
	}

	url, err := s.oss.GetURL(normalizedPath)
	if err != nil {
		return "", errorx.WrapError(err, fmt.Sprintf("获取文件URL失败: %s", normalizedPath))
	}

	// 确保URL使用正斜杠
	return strings.ReplaceAll(url, "\\", "/"), nil
}

// GetEndpoint 获取存储端点
func (s *Storage) GetEndpoint() string {
	return s.oss.GetEndpoint()
}

// 生成存储路径
func GeneratePath(baseDir, fileName string) string {
	// 清理文件名，移除危险字符
	fileName = filepath.Base(fileName)

	// 文件路径格式：baseDir/fileName
	return filepath.Join(baseDir, fileName)
}

// 常用文件类型的MIME映射
var MimeTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".bmp":  "image/bmp",
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".mp3":  "audio/mpeg",
	".mp4":  "video/mp4",
	".zip":  "application/zip",
	".rar":  "application/x-rar-compressed",
	".txt":  "text/plain",
	".html": "text/html",
	".css":  "text/css",
	".js":   "application/javascript",
	".json": "application/json",
	".xml":  "application/xml",
}

// GetMimeType 获取文件MIME类型
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if mime, ok := MimeTypes[ext]; ok {
		return mime
	}
	return "application/octet-stream" // 默认二进制流
}
