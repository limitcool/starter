package handler

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/limitcool/starter/internal/model"
)

// FileUtil 文件工具类
type FileUtil struct {
	baseURL string // 基础URL，如 "/uploads" 或 "https://cdn.example.com"
}

// NewFileUtil 创建文件工具实例
func NewFileUtil(baseURL string) *FileUtil {
	return &FileUtil{
		baseURL: baseURL,
	}
}

// BuildFileURL 构建文件访问URL
func (f *FileUtil) BuildFileURL(path string) string {
	if path == "" {
		return ""
	}
	
	// 确保路径以正斜杠开头
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	
	// 如果baseURL已经是完整URL，直接拼接
	if strings.HasPrefix(f.baseURL, "http") {
		return f.baseURL + path
	}
	
	// 否则作为相对路径处理
	return f.baseURL + path
}

// GenerateFileName 生成唯一的文件名
func (f *FileUtil) GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randString(8), ext)
	return name
}

// GetStoragePath 获取存储路径
func (f *FileUtil) GetStoragePath(fileType, fileName string, isPublic bool) string {
	// 根据是否公开选择根目录
	var rootDir string
	if isPublic {
		rootDir = "public"
	} else {
		rootDir = "private"
	}

	// 根据文件类型选择子目录
	var typeDir string
	switch fileType {
	case model.FileTypeImage:
		typeDir = "images"
	case model.FileTypeDocument:
		typeDir = "documents"
	case model.FileTypeVideo:
		typeDir = "videos"
	case model.FileTypeAudio:
		typeDir = "audios"
	default:
		typeDir = "others"
	}

	// 添加日期子目录
	dateDir := time.Now().Format("2006/01/02")

	// 构建完整路径：public/documents/2025/06/17/filename.txt
	path := filepath.Join(rootDir, typeDir, dateDir, fileName)
	// 确保Windows上也使用正斜杠
	return strings.ReplaceAll(path, "\\", "/")
}

// SetFileURL 为文件模型设置URL字段
func (f *FileUtil) SetFileURL(file *model.File) {
	if file != nil && file.Path != "" {
		file.URL = f.BuildFileURL(file.Path)
	}
}

// SetFileURLs 为文件列表设置URL字段
func (f *FileUtil) SetFileURLs(files []model.File) {
	for i := range files {
		f.SetFileURL(&files[i])
	}
}

// IsAllowedFileType 检查文件类型是否允许
func (f *FileUtil) IsAllowedFileType(ext, fileType string) bool {
	ext = strings.ToLower(ext)
	
	switch fileType {
	case model.FileTypeImage:
		return isImageExt(ext)
	case model.FileTypeDocument:
		return isDocumentExt(ext)
	case model.FileTypeVideo:
		return isVideoExt(ext)
	case model.FileTypeAudio:
		return isAudioExt(ext)
	default:
		return true // 其他类型默认允许
	}
}

// IsAllowedFileSize 检查文件大小是否允许
func (f *FileUtil) IsAllowedFileSize(size int64, fileType string) bool {
	const (
		MB = 1024 * 1024
		GB = 1024 * MB
	)
	
	switch fileType {
	case model.FileTypeImage:
		return size <= 10*MB // 图片最大10MB
	case model.FileTypeDocument:
		return size <= 50*MB // 文档最大50MB
	case model.FileTypeVideo:
		return size <= 500*MB // 视频最大500MB
	case model.FileTypeAudio:
		return size <= 100*MB // 音频最大100MB
	default:
		return size <= 100*MB // 其他类型最大100MB
	}
}

// randString 生成随机字符串
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[int(time.Now().UnixNano()%int64(len(letters)))]
		time.Sleep(1 * time.Nanosecond) // 确保随机性
	}
	return string(b)
}

// isImageExt 检查是否为图片扩展名
func isImageExt(ext string) bool {
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}
	for _, allowedExt := range imageExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// isDocumentExt 检查是否为文档扩展名
func isDocumentExt(ext string) bool {
	docExts := []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".rtf"}
	for _, allowedExt := range docExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// isVideoExt 检查是否为视频扩展名
func isVideoExt(ext string) bool {
	videoExts := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm", ".mkv"}
	for _, allowedExt := range videoExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// isAudioExt 检查是否为音频扩展名
func isAudioExt(ext string) bool {
	audioExts := []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma"}
	for _, allowedExt := range audioExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
