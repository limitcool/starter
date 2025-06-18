package filestore

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/limitcool/starter/internal/pkg/idgen"
)

// FileUsage 文件用途类型
type FileUsage string

const (
	// 用户相关
	FileUsageAvatar  FileUsage = "avatar"  // 用户头像
	FileUsageProfile FileUsage = "profile" // 用户资料图片

	// 内容相关
	FileUsageBanner  FileUsage = "banner"  // 横幅图片
	FileUsageCover   FileUsage = "cover"   // 封面图片
	FileUsageGallery FileUsage = "gallery" // 相册图片
	FileUsagePost    FileUsage = "post"    // 文章图片

	// 文档相关
	FileUsageDocument FileUsage = "document" // 文档文件
	FileUsageReport   FileUsage = "report"   // 报告文件
	FileUsageContract FileUsage = "contract" // 合同文件

	// 媒体相关
	FileUsageVideo FileUsage = "video" // 视频文件
	FileUsageAudio FileUsage = "audio" // 音频文件

	// 系统相关
	FileUsageTemp    FileUsage = "temp"    // 临时文件
	FileUsageBackup  FileUsage = "backup"  // 备份文件
	FileUsageGeneral FileUsage = "general" // 通用文件
)

// PathManager 路径管理器
type PathManager struct {
	// 路径配置
	pathConfig map[FileUsage]PathConfig
}

// PathConfig 路径配置
type PathConfig struct {
	BaseDir     string   // 基础目录
	UseDate     bool     // 是否使用日期分组
	DateFormat  string   // 日期格式
	MaxFileSize int64    // 最大文件大小 (字节)
	AllowedExts []string // 允许的文件扩展名
}

// NewPathManager 创建路径管理器
func NewPathManager() *PathManager {
	return &PathManager{
		pathConfig: map[FileUsage]PathConfig{
			// 用户相关 - 按用户ID分组
			FileUsageAvatar: {
				BaseDir:     "users/avatars",
				UseDate:     false,
				MaxFileSize: 5 * 1024 * 1024, // 5MB
				AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp"},
			},
			FileUsageProfile: {
				BaseDir:     "users/profiles",
				UseDate:     false,
				MaxFileSize: 10 * 1024 * 1024, // 10MB
				AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp"},
			},

			// 内容相关 - 按日期分组
			FileUsageBanner: {
				BaseDir:     "content/banners",
				UseDate:     true,
				DateFormat:  "2006/01",
				MaxFileSize: 20 * 1024 * 1024, // 20MB
				AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp", ".svg"},
			},
			FileUsageCover: {
				BaseDir:     "content/covers",
				UseDate:     true,
				DateFormat:  "2006/01",
				MaxFileSize: 15 * 1024 * 1024, // 15MB
				AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp"},
			},
			FileUsageGallery: {
				BaseDir:     "content/gallery",
				UseDate:     true,
				DateFormat:  "2006/01/02",
				MaxFileSize: 10 * 1024 * 1024, // 10MB
				AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp"},
			},
			FileUsagePost: {
				BaseDir:     "content/posts",
				UseDate:     true,
				DateFormat:  "2006/01",
				MaxFileSize: 10 * 1024 * 1024, // 10MB
				AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp", ".gif"},
			},

			// 文档相关 - 按年月分组
			FileUsageDocument: {
				BaseDir:     "documents/general",
				UseDate:     true,
				DateFormat:  "2006/01",
				MaxFileSize: 50 * 1024 * 1024, // 50MB
				AllowedExts: []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt"},
			},
			FileUsageReport: {
				BaseDir:     "documents/reports",
				UseDate:     true,
				DateFormat:  "2006/01",
				MaxFileSize: 100 * 1024 * 1024, // 100MB
				AllowedExts: []string{".pdf", ".doc", ".docx", ".xls", ".xlsx"},
			},
			FileUsageContract: {
				BaseDir:     "documents/contracts",
				UseDate:     true,
				DateFormat:  "2006",
				MaxFileSize: 50 * 1024 * 1024, // 50MB
				AllowedExts: []string{".pdf", ".doc", ".docx"},
			},

			// 媒体相关 - 按年月分组
			FileUsageVideo: {
				BaseDir:     "media/videos",
				UseDate:     true,
				DateFormat:  "2006/01",
				MaxFileSize: 500 * 1024 * 1024, // 500MB
				AllowedExts: []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm"},
			},
			FileUsageAudio: {
				BaseDir:     "media/audio",
				UseDate:     true,
				DateFormat:  "2006/01",
				MaxFileSize: 100 * 1024 * 1024, // 100MB
				AllowedExts: []string{".mp3", ".wav", ".flac", ".aac", ".ogg"},
			},

			// 系统相关
			FileUsageTemp: {
				BaseDir:     "system/temp",
				UseDate:     true,
				DateFormat:  "2006/01/02",
				MaxFileSize: 100 * 1024 * 1024, // 100MB
				AllowedExts: []string{},        // 允许所有类型
			},
			FileUsageBackup: {
				BaseDir:     "system/backups",
				UseDate:     true,
				DateFormat:  "2006/01",
				MaxFileSize: 1024 * 1024 * 1024, // 1GB
				AllowedExts: []string{".zip", ".tar", ".gz", ".sql"},
			},
			FileUsageGeneral: {
				BaseDir:     "general",
				UseDate:     true,
				DateFormat:  "2006/01/02",
				MaxFileSize: 50 * 1024 * 1024, // 50MB
				AllowedExts: []string{},       // 允许所有类型
			},
		},
	}
}

// GenerateFilePath 生成文件路径
func (pm *PathManager) GenerateFilePath(usage FileUsage, originalName string, userID ...int64) (string, string, error) {
	config, exists := pm.pathConfig[usage]
	if !exists {
		config = pm.pathConfig[FileUsageGeneral] // 默认使用通用配置
	}

	// 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(originalName))
	if len(config.AllowedExts) > 0 && !pm.isExtensionAllowed(ext, config.AllowedExts) {
		return "", "", fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 生成唯一文件名
	fileID := idgen.GenerateUUID()
	filename := fmt.Sprintf("%s%s", fileID, ext)

	// 构建路径
	var pathParts []string
	pathParts = append(pathParts, config.BaseDir)

	// 添加日期路径
	if config.UseDate {
		datePath := time.Now().Format(config.DateFormat)
		pathParts = append(pathParts, datePath)
	}

	// 对于用户相关文件，添加用户ID路径
	if pm.isUserRelated(usage) && len(userID) > 0 {
		userPath := fmt.Sprintf("user_%d", userID[0])
		pathParts = append(pathParts, userPath)
	}

	// 组合完整路径
	filePath := filepath.Join(pathParts...)
	filePath = filepath.Join(filePath, filename)

	// 转换为Unix风格路径（用于URL）
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	return filePath, fileID, nil
}

// ValidateFile 验证文件
func (pm *PathManager) ValidateFile(usage FileUsage, filename string, size int64) error {
	config, exists := pm.pathConfig[usage]
	if !exists {
		config = pm.pathConfig[FileUsageGeneral]
	}

	// 验证文件大小
	if config.MaxFileSize > 0 && size > config.MaxFileSize {
		return fmt.Errorf("文件大小超过限制: %d bytes (最大: %d bytes)", size, config.MaxFileSize)
	}

	// 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	if len(config.AllowedExts) > 0 && !pm.isExtensionAllowed(ext, config.AllowedExts) {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}

	return nil
}

// GetUsageFromString 从字符串获取文件用途
func (pm *PathManager) GetUsageFromString(usageStr string) FileUsage {
	usage := FileUsage(usageStr)
	if _, exists := pm.pathConfig[usage]; exists {
		return usage
	}
	return FileUsageGeneral
}

// GetAllowedExtensions 获取允许的扩展名
func (pm *PathManager) GetAllowedExtensions(usage FileUsage) []string {
	if config, exists := pm.pathConfig[usage]; exists {
		return config.AllowedExts
	}
	return []string{}
}

// GetMaxFileSize 获取最大文件大小
func (pm *PathManager) GetMaxFileSize(usage FileUsage) int64 {
	if config, exists := pm.pathConfig[usage]; exists {
		return config.MaxFileSize
	}
	return 50 * 1024 * 1024 // 默认50MB
}

// 辅助方法
func (pm *PathManager) isExtensionAllowed(ext string, allowedExts []string) bool {
	if len(allowedExts) == 0 {
		return true // 允许所有类型
	}
	for _, allowed := range allowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (pm *PathManager) isUserRelated(usage FileUsage) bool {
	return usage == FileUsageAvatar || usage == FileUsageProfile
}
