package handler

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/filestore"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// FileHandler 文件处理器
type FileHandler struct {
	db      *gorm.DB
	config  *configs.Config
	logger  *logger.Logger
	storage *filestore.Storage
}

// NewFileHandler 创建文件处理器
func NewFileHandler(db *gorm.DB, config *configs.Config, storage *filestore.Storage) *FileHandler {
	handler := &FileHandler{
		db:      db,
		config:  config,
		storage: storage,
	}

	logger.Info("FileHandler initialized")
	return handler
}

// UploadFile 上传文件
func (h *FileHandler) UploadFile(ctx *gin.Context) {
	// 获取请求上下文
	reqCtx := ctx.Request.Context()

	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.WarnContext(reqCtx, "UploadFile 未找到用户ID")
		response.Error(ctx, errorx.ErrUserNoLogin)
		return
	}

	// 转换用户ID
	id := cast.ToInt64(userID)

	// 获取用户类型
	userTypeValue, exists := ctx.Get("user_type")
	if !exists {
		logger.WarnContext(reqCtx, "UploadFile 未找到用户类型")
		response.Error(ctx, errorx.ErrUserNoLogin)
		return
	}

	// 获取用户类型（简化为数字）
	userType := cast.ToUint8(userTypeValue)

	// 获取文件
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		logger.WarnContext(reqCtx, "UploadFile 获取文件失败",
			"error", err,
			"user_id", id)
		response.Error(ctx, errorx.ErrInvalidParams.WithError(err))
		return
	}

	// 获取文件类型
	fileType := ctx.DefaultPostForm("type", model.FileTypeOther)

	// 获取是否公开
	isPublic := ctx.DefaultPostForm("is_public", "false") == "true"

	// 获取文件基本信息
	originalName := fileHeader.Filename
	ext := strings.ToLower(filepath.Ext(originalName))
	mimeType := filestore.GetMimeType(originalName)
	size := fileHeader.Size

	// 验证文件类型
	if !isAllowedFileType(ext, fileType) {
		logger.WarnContext(reqCtx, "UploadFile 文件类型不允许",
			"user_id", id,
			"file_type", fileType,
			"ext", ext)
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件类型不允许"))
		return
	}

	// 验证文件大小
	if !isAllowedFileSize(size, fileType) {
		logger.WarnContext(reqCtx, "UploadFile 文件大小超出限制",
			"user_id", id,
			"file_type", fileType,
			"size", size)
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件大小超出限制"))
		return
	}

	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		logger.ErrorContext(reqCtx, "UploadFile 打开文件失败",
			"error", err,
			"user_id", id)
		response.Error(ctx, errorx.ErrInvalidParams.WithError(err))
		return
	}
	defer file.Close()

	// 生成文件名和存储路径
	fileName := generateFileName(originalName)
	storagePath := getStoragePath(fileType, fileName)

	// 上传文件到存储
	err = h.storage.Put(reqCtx, storagePath, file)
	if err != nil {
		logger.ErrorContext(reqCtx, "UploadFile 上传文件到存储失败",
			"error", err,
			"user_id", id,
			"storage_path", storagePath)
		response.Error(ctx, errorx.WrapError(err, fmt.Sprintf("上传文件到存储失败: %s", storagePath)))
		return
	}

	// 获取文件访问URL
	fileURL, err := h.storage.GetURL(reqCtx, storagePath)
	if err != nil {
		logger.ErrorContext(reqCtx, "UploadFile 获取文件访问URL失败",
			"error", err,
			"user_id", id,
			"storage_path", storagePath)
		response.Error(ctx, errorx.WrapError(err, fmt.Sprintf("获取文件访问URL失败: %s", storagePath)))
		return
	}

	// 创建文件仓库
	fileRepo := model.NewFileRepo(h.db)

	// 记录到数据库
	fileModel := &model.File{
		Name:           fileName,
		OriginalName:   originalName,
		Path:           storagePath,
		URL:            fileURL,
		Type:           fileType,
		Usage:          model.FileUsageGeneral, // 默认用途为通用
		Size:           size,
		MimeType:       mimeType,
		Extension:      ext,
		StorageType:    string(h.storage.Config.Type),
		UploadedBy:     id,
		UploadedByType: uint8(userType),
		UploadedAt:     time.Now(),
		Status:         1,        // 状态正常
		IsPublic:       isPublic, // 设置是否公开
	}

	if err := fileRepo.Create(reqCtx, fileModel); err != nil {
		logger.ErrorContext(reqCtx, "UploadFile 保存文件记录失败",
			"error", err,
			"user_id", id)
		response.Error(ctx, err)
		return
	}

	logger.InfoContext(reqCtx, "UploadFile 上传文件成功",
		"user_id", id,
		"file_id", fileModel.ID,
		"file_name", fileModel.Name)

	response.Success(ctx, fileModel)
}

// 生成文件名
func generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randString(8), ext)
	return name
}

// 随机字符串
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[int(time.Now().UnixNano()%int64(len(letters)))]
		time.Sleep(1 * time.Nanosecond) // 确保随机性
	}
	return string(b)
}

// 获取存储路径
func getStoragePath(fileType, fileName string) string {
	var baseDir string
	switch fileType {
	case model.FileTypeImage:
		baseDir = "images"
	case model.FileTypeDocument:
		baseDir = "documents"
	case model.FileTypeVideo:
		baseDir = "videos"
	case model.FileTypeAudio:
		baseDir = "audios"
	default:
		baseDir = "others"
	}

	// 添加日期子目录
	dateDir := time.Now().Format("2006/01/02")

	// 使用path包而不是filepath包，确保始终使用/作为分隔符
	path := filepath.Join(baseDir, dateDir, fileName)
	// 确保Windows上也使用正斜杠
	return strings.ReplaceAll(path, "\\", "/")
}

// 检查文件类型是否允许
func isAllowedFileType(ext string, fileType string) bool {
	ext = strings.ToLower(ext)
	switch fileType {
	case model.FileTypeImage:
		return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" || ext == ".svg"
	case model.FileTypeDocument:
		return ext == ".pdf" || ext == ".doc" || ext == ".docx" || ext == ".xls" || ext == ".xlsx" || ext == ".txt"
	case model.FileTypeVideo:
		return ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".wmv" || ext == ".flv"
	case model.FileTypeAudio:
		return ext == ".mp3" || ext == ".wav" || ext == ".ogg" || ext == ".flac"
	default:
		return true
	}
}

// 检查文件大小是否允许
func isAllowedFileSize(size int64, fileType string) bool {
	const (
		MB = 1024 * 1024
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
		return size <= 50*MB // 其他类型最大50MB
	}
}
