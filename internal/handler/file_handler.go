package handler

import (
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/filestore"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// FileHandler 文件处理器（基于接口）
type FileHandler struct {
	db          *gorm.DB
	storage     filestore.FileStorage
	pathManager *filestore.PathManager
}

// NewFileHandler 创建文件处理器
func NewFileHandler(db *gorm.DB, storage filestore.FileStorage) *FileHandler {
	return &FileHandler{
		db:          db,
		storage:     storage,
		pathManager: filestore.NewPathManager(),
	}
}

// GetUploadURL 获取上传URL
func (h *FileHandler) GetUploadURL(ctx *gin.Context) {
	var req dto.FileUploadRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams.WithError(err))
		return
	}

	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, errorx.ErrUnauthorized.WithMsg("用户未登录"))
		return
	}

	// 获取文件用途
	usage := h.pathManager.GetUsageFromString(req.Usage)

	// 预验证文件（如果提供了大小）
	if req.Size > 0 {
		if err := h.pathManager.ValidateFile(usage, req.Filename, req.Size); err != nil {
			response.Error(ctx, errorx.ErrInvalidParams.WithMsg(err.Error()))
			return
		}
	}

	// 生成智能文件路径
	filePath, _, err := h.pathManager.GenerateFilePath(usage, req.Filename, cast.ToInt64(userID))
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	// 获取上传URL
	uploadURL, method, err := h.storage.GetUploadURL(ctx.Request.Context(), filePath, req.ContentType, req.IsPublic)
	if err != nil {
		logger.ErrorContext(ctx.Request.Context(), "获取上传URL失败", "error", err)
		response.Error(ctx, errorx.ErrFileStorage.WithMsg("获取上传URL失败"))
		return
	}

	// 创建文件记录
	ext := filepath.Ext(req.Filename)
	fileRecord := &model.File{
		Name:         filepath.Base(filePath),
		OriginalName: req.Filename,
		Path:         filePath,
		MimeType:     req.ContentType,
		Extension:    ext,
		Usage:        req.Usage,
		StorageType:  h.storage.GetStorageType(),
		UploadedBy:   cast.ToInt64(userID),
		IsPublic:     req.IsPublic,
		Status:       0, // 待上传
	}

	if err := h.db.Create(fileRecord).Error; err != nil {
		logger.ErrorContext(ctx.Request.Context(), "创建文件记录失败", "error", err)
		response.Error(ctx, errorx.ErrDatabase.WithMsg("创建文件记录失败"))
		return
	}

	response.Success(ctx, &dto.FileUploadResponse{
		FileID:      fileRecord.ID,
		UploadURL:   uploadURL,
		Method:      method,
		ExpiresIn:   15, // 分钟
		StorageType: h.storage.GetStorageType(),
		Usage:       req.Usage,
		PathInfo: dto.PathInfo{
			Category: string(usage),
			Path:     filePath,
		},
	})
}

// ConfirmUpload 确认上传完成
func (h *FileHandler) ConfirmUpload(ctx *gin.Context) {
	var req dto.FileConfirmRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams.WithError(err))
		return
	}

	// 获取文件记录
	var fileRecord model.File
	if err := h.db.Where("id = ?", req.FileID).First(&fileRecord).Error; err != nil {
		response.Error(ctx, errorx.ErrNotFound.WithMsg("文件记录不存在"))
		return
	}

	// 检查文件是否存在
	exists, err := h.storage.FileExists(ctx.Request.Context(), fileRecord.Path, fileRecord.IsPublic)
	if err != nil {
		logger.ErrorContext(ctx.Request.Context(), "检查文件存在性失败", "error", err)
		response.Error(ctx, errorx.ErrFileStorage.WithMsg("验证文件失败"))
		return
	}

	if !exists {
		response.Error(ctx, errorx.ErrNotFound.WithMsg("文件上传未完成"))
		return
	}

	// 生成下载URL
	downloadURL, err := h.storage.GetDownloadURL(ctx.Request.Context(), fileRecord.Path, fileRecord.IsPublic)
	if err != nil {
		logger.ErrorContext(ctx.Request.Context(), "生成下载URL失败", "error", err)
		response.Error(ctx, errorx.ErrFileStorage.WithMsg("生成下载URL失败"))
		return
	}

	// 更新文件记录
	fileRecord.Size = req.Size
	fileRecord.URL = downloadURL
	fileRecord.Status = 1 // 已完成
	fileRecord.UploadedAt = time.Now()

	if err := h.db.Save(&fileRecord).Error; err != nil {
		logger.ErrorContext(ctx.Request.Context(), "更新文件记录失败", "error", err)
		response.Error(ctx, errorx.ErrDatabase.WithMsg("更新文件记录失败"))
		return
	}

	response.Success(ctx, fileRecord)
}

// GetFileInfo 获取文件信息
func (h *FileHandler) GetFileInfo(ctx *gin.Context) {
	fileID := ctx.Param("id")
	if fileID == "" {
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件ID不能为空"))
		return
	}

	var fileRecord model.File
	if err := h.db.Where("id = ?", fileID).First(&fileRecord).Error; err != nil {
		response.Error(ctx, errorx.ErrNotFound.WithMsg("文件不存在"))
		return
	}

	response.Success(ctx, fileRecord)
}

// GetDownloadURL 获取下载URL
func (h *FileHandler) GetDownloadURL(ctx *gin.Context) {
	fileID := ctx.Param("id")
	if fileID == "" {
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件ID不能为空"))
		return
	}

	var fileRecord model.File
	if err := h.db.Where("id = ?", fileID).First(&fileRecord).Error; err != nil {
		response.Error(ctx, errorx.ErrNotFound.WithMsg("文件不存在"))
		return
	}

	// 生成新的下载URL
	downloadURL, err := h.storage.GetDownloadURL(ctx.Request.Context(), fileRecord.Path, fileRecord.IsPublic)
	if err != nil {
		logger.ErrorContext(ctx.Request.Context(), "生成下载URL失败", "error", err)
		response.Error(ctx, errorx.ErrFileStorage.WithMsg("生成下载URL失败"))
		return
	}

	response.Success(ctx, &dto.FileDownloadResponse{
		FileID:      fileRecord.ID,
		Filename:    fileRecord.OriginalName,
		DownloadURL: downloadURL,
		IsPublic:    fileRecord.IsPublic,
		Size:        fileRecord.Size,
		StorageType: fileRecord.StorageType,
	})
}

// UploadFile 统一文件上传接口（支持本地和MinIO）
func (h *FileHandler) UploadFile(ctx *gin.Context) {
	var req struct {
		Filename    string `form:"filename" binding:"required"`
		ContentType string `form:"content_type"`
		IsPublic    bool   `form:"is_public"`
		Usage       string `form:"usage" binding:"required"` // avatar, banner, document, etc.
	}

	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams.WithError(err))
		return
	}

	// 获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, errorx.ErrUnauthorized.WithMsg("用户未登录"))
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("获取上传文件失败"))
		return
	}

	// 如果没有提供content_type，从文件头推断
	if req.ContentType == "" {
		req.ContentType = file.Header.Get("Content-Type")
		if req.ContentType == "" {
			req.ContentType = "application/octet-stream"
		}
	}

	// 获取文件用途
	usage := h.pathManager.GetUsageFromString(req.Usage)

	// 验证文件
	if err := h.pathManager.ValidateFile(usage, req.Filename, file.Size); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	// 生成智能文件路径
	filePath, _, err := h.pathManager.GenerateFilePath(usage, req.Filename, cast.ToInt64(userID))
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("打开上传文件失败"))
		return
	}
	defer src.Close()

	// 创建文件记录
	ext := filepath.Ext(req.Filename)
	fileRecord := &model.File{
		Name:         filepath.Base(filePath),
		OriginalName: req.Filename,
		Path:         filePath,
		Size:         file.Size,
		MimeType:     req.ContentType,
		Extension:    ext,
		Usage:        req.Usage,
		StorageType:  h.storage.GetStorageType(),
		UploadedBy:   cast.ToInt64(userID),
		IsPublic:     req.IsPublic,
		Status:       1, // 已完成
		UploadedAt:   time.Now(),
	}

	// 上传文件到存储
	if err := h.storage.UploadFile(ctx.Request.Context(), filePath, src, req.IsPublic); err != nil {
		logger.ErrorContext(ctx.Request.Context(), "文件上传失败", "error", err)
		response.Error(ctx, errorx.ErrFileStorage.WithMsg("文件上传失败"))
		return
	}

	// 生成下载URL
	downloadURL, err := h.storage.GetDownloadURL(ctx.Request.Context(), filePath, req.IsPublic)
	if err != nil {
		logger.ErrorContext(ctx.Request.Context(), "生成下载URL失败", "error", err)
		// 不返回错误，继续保存记录
	} else {
		fileRecord.URL = downloadURL
	}

	// 保存文件记录到数据库
	if err := h.db.Create(fileRecord).Error; err != nil {
		logger.ErrorContext(ctx.Request.Context(), "创建文件记录失败", "error", err)
		response.Error(ctx, errorx.ErrDatabase.WithMsg("创建文件记录失败"))
		return
	}

	logger.InfoContext(ctx.Request.Context(), "文件上传成功",
		"file_id", fileRecord.ID,
		"filename", req.Filename,
		"storage_type", h.storage.GetStorageType(),
		"user_id", userID)

	response.Success(ctx, &dto.FileUploadCompleteResponse{
		FileID:      fileRecord.ID,
		Filename:    fileRecord.OriginalName,
		DownloadURL: downloadURL,
		Size:        fileRecord.Size,
		StorageType: fileRecord.StorageType,
		IsPublic:    fileRecord.IsPublic,
		Usage:       fileRecord.Usage,
	})
}

// DeleteFile 删除文件
func (h *FileHandler) DeleteFile(ctx *gin.Context) {
	fileID := ctx.Param("id")
	if fileID == "" {
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件ID不能为空"))
		return
	}

	var fileRecord model.File
	if err := h.db.Where("id = ?", fileID).First(&fileRecord).Error; err != nil {
		response.Error(ctx, errorx.ErrNotFound.WithMsg("文件不存在"))
		return
	}

	// 删除存储中的文件
	if err := h.storage.DeleteFile(ctx.Request.Context(), fileRecord.Path, fileRecord.IsPublic); err != nil {
		logger.ErrorContext(ctx.Request.Context(), "删除存储文件失败", "error", err)
		response.Error(ctx, errorx.ErrFileStorage.WithMsg("删除文件失败"))
		return
	}

	// 删除数据库记录
	if err := h.db.Delete(&fileRecord).Error; err != nil {
		logger.ErrorContext(ctx.Request.Context(), "删除文件记录失败", "error", err)
		response.Error(ctx, errorx.ErrDatabase.WithMsg("删除文件记录失败"))
		return
	}

	response.Success(ctx, &dto.DeleteResponse{Message: "删除成功"})
}
