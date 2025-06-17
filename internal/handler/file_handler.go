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
	*BaseHandler
	storage *filestore.Storage
}

// NewFileHandler 创建文件处理器
func NewFileHandler(db *gorm.DB, config *configs.Config, storage *filestore.Storage) *FileHandler {
	handler := &FileHandler{
		BaseHandler: NewBaseHandler(db, config),
		storage:     storage,
	}

	handler.LogInit("FileHandler")
	return handler
}

// GetUploadURL 获取文件上传预签名URL（前端直接上传到MinIO）
func (h *FileHandler) GetUploadURL(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	// 获取当前用户信息
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.WarnContext(reqCtx, "GetUploadURL 未找到用户ID")
		response.Error(ctx, errorx.ErrUnauthorized.WithMsg("用户未登录"))
		return
	}
	id := cast.ToInt64(userID)

	// 获取用户类型
	userType := uint8(1) // 默认为普通用户
	if isAdmin, exists := ctx.Get("is_admin"); exists && cast.ToBool(isAdmin) {
		userType = uint8(2) // 管理员
	}

	// 解析请求参数
	var req struct {
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type" binding:"required"`
		FileType    string `json:"file_type" binding:"required"` // image, document, video, audio
		IsPublic    bool   `json:"is_public"`
		Usage       string `json:"usage"` // avatar, document, etc.
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.WarnContext(reqCtx, "GetUploadURL 请求参数无效", "error", err)
		response.Error(ctx, errorx.ErrInvalidParams.WithError(err))
		return
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(req.Filename))
	if !h.FileUtil.IsAllowedFileType(ext, req.FileType) {
		logger.WarnContext(reqCtx, "GetUploadURL 不支持的文件类型",
			"user_id", id,
			"file_type", req.FileType,
			"extension", ext)
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("不支持的文件类型"))
		return
	}

	// 生成唯一的文件名和存储路径
	fileName := h.FileUtil.GenerateFileName(req.Filename)
	storagePath := h.FileUtil.GetStoragePath(req.FileType, fileName, req.IsPublic)

	// 生成上传预签名URL
	uploadURL, err := h.storage.GetUploadPresignedURL(reqCtx, storagePath, req.ContentType, 15) // 15分钟有效期
	if err != nil {
		logger.ErrorContext(reqCtx, "GetUploadURL 生成上传预签名URL失败",
			"error", err,
			"user_id", id,
			"storage_path", storagePath)
		response.Error(ctx, errorx.WrapError(err, "生成上传链接失败"))
		return
	}

	// 创建文件记录（状态为pending，等待上传完成确认）
	fileRepo := model.NewFileRepo(h.DB)
	fileModel := &model.File{
		Name:           fileName,
		OriginalName:   req.Filename,
		Path:           storagePath,
		URL:            "", // 上传完成后再设置
		Type:           req.FileType,
		Usage:          req.Usage,
		Size:           0, // 上传完成后再设置
		MimeType:       req.ContentType,
		Extension:      ext,
		StorageType:    string(h.storage.Config.Type),
		UploadedBy:     id,
		UploadedByType: uint8(userType),
		UploadedAt:     time.Now(),
		Status:         0, // 0=pending, 1=completed
		IsPublic:       req.IsPublic,
	}

	if err := fileRepo.Create(reqCtx, fileModel); err != nil {
		logger.ErrorContext(reqCtx, "GetUploadURL 创建文件记录失败",
			"error", err,
			"user_id", id)
		response.Error(ctx, err)
		return
	}

	// 构建响应数据
	result := map[string]interface{}{
		"file_id":      fileModel.ID,
		"upload_url":   uploadURL,
		"storage_path": storagePath,
		"expires_in":   15, // 分钟
		"expires_at":   time.Now().Add(15 * time.Minute).Unix(),
		"method":       "PUT", // 上传方法
		"headers": map[string]string{
			"Content-Type": req.ContentType,
		},
	}

	logger.InfoContext(reqCtx, "GetUploadURL 生成上传URL成功",
		"user_id", id,
		"file_id", fileModel.ID,
		"storage_path", storagePath,
		"is_public", req.IsPublic)

	response.Success(ctx, result)
}

// ConfirmUpload 确认文件上传完成
func (h *FileHandler) ConfirmUpload(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	// 获取当前用户信息
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.WarnContext(reqCtx, "ConfirmUpload 未找到用户ID")
		response.Error(ctx, errorx.ErrUnauthorized.WithMsg("用户未登录"))
		return
	}
	id := cast.ToInt64(userID)

	// 解析请求参数
	var req struct {
		FileID uint  `json:"file_id" binding:"required"`
		Size   int64 `json:"size" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.WarnContext(reqCtx, "ConfirmUpload 请求参数无效", "error", err)
		response.Error(ctx, errorx.ErrInvalidParams.WithError(err))
		return
	}

	// 创建文件仓库
	fileRepo := model.NewFileRepo(h.DB)

	// 获取文件记录
	fileModel, err := fileRepo.GetByID(reqCtx, req.FileID)
	if err != nil {
		logger.ErrorContext(reqCtx, "ConfirmUpload 获取文件记录失败",
			"error", err,
			"file_id", req.FileID)
		response.Error(ctx, err)
		return
	}

	// 检查文件所有权
	if fileModel.UploadedBy != id {
		// 检查是否是管理员
		isAdmin, exists := ctx.Get("is_admin")
		if !exists || !cast.ToBool(isAdmin) {
			logger.WarnContext(reqCtx, "ConfirmUpload 无权限操作文件",
				"file_id", req.FileID,
				"user_id", id,
				"uploaded_by", fileModel.UploadedBy)
			response.Error(ctx, errorx.ErrForbidden.WithMsg("无权限操作此文件"))
			return
		}
	}

	// 检查文件状态
	if fileModel.Status != 0 {
		logger.WarnContext(reqCtx, "ConfirmUpload 文件已确认",
			"file_id", req.FileID,
			"status", fileModel.Status)
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件已确认，无需重复操作"))
		return
	}

	// 验证文件是否真的存在于存储中
	exists, err = h.storage.Exists(reqCtx, fileModel.Path)
	if err != nil {
		logger.ErrorContext(reqCtx, "ConfirmUpload 检查文件存在性失败",
			"error", err,
			"file_id", req.FileID,
			"path", fileModel.Path)
		response.Error(ctx, errorx.WrapError(err, "验证文件失败"))
		return
	}

	if !exists {
		logger.WarnContext(reqCtx, "ConfirmUpload 文件不存在于存储中",
			"file_id", req.FileID,
			"path", fileModel.Path)
		response.Error(ctx, errorx.ErrNotFound.WithMsg("文件上传未完成或已被删除"))
		return
	}

	// 生成文件访问URL
	fileURL, err := h.storage.GetURL(reqCtx, fileModel.Path)
	if err != nil {
		logger.ErrorContext(reqCtx, "ConfirmUpload 获取文件访问URL失败",
			"error", err,
			"file_id", req.FileID,
			"path", fileModel.Path)
		response.Error(ctx, errorx.WrapError(err, "生成文件访问URL失败"))
		return
	}

	// 更新文件记录
	fileModel.Size = req.Size
	fileModel.URL = fileURL
	fileModel.Status = 1 // 完成状态
	fileModel.UploadedAt = time.Now()

	if err := fileRepo.Update(reqCtx, fileModel); err != nil {
		logger.ErrorContext(reqCtx, "ConfirmUpload 更新文件记录失败",
			"error", err,
			"file_id", req.FileID)
		response.Error(ctx, err)
		return
	}

	logger.InfoContext(reqCtx, "ConfirmUpload 确认上传成功",
		"user_id", id,
		"file_id", req.FileID,
		"size", req.Size)

	response.Success(ctx, fileModel)
}

// UploadFile 上传文件（保留旧接口，但标记为废弃）
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

	// 获取用户类型 - 根据is_admin字段推断
	isAdminValue, exists := ctx.Get("is_admin")
	if !exists {
		logger.WarnContext(reqCtx, "UploadFile 未找到用户管理员标识")
		response.Error(ctx, errorx.ErrUserNoLogin)
		return
	}

	// 根据is_admin推断用户类型（1=管理员，2=普通用户）
	var userType uint8 = 2 // 默认为普通用户
	if isAdmin, ok := isAdminValue.(bool); ok && isAdmin {
		userType = 1 // 管理员
	}

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
	if !h.FileUtil.IsAllowedFileType(ext, fileType) {
		logger.WarnContext(reqCtx, "UploadFile 文件类型不允许",
			"user_id", id,
			"file_type", fileType,
			"ext", ext)
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件类型不允许"))
		return
	}

	// 验证文件大小
	if !h.FileUtil.IsAllowedFileSize(size, fileType) {
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
	fileName := h.FileUtil.GenerateFileName(originalName)
	storagePath := h.FileUtil.GetStoragePath(fileType, fileName, isPublic)

	// 上传文件到存储（权限由路径和Bucket Policy控制）
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
	fileRepo := model.NewFileRepo(h.DB)

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

// GetFileURL 获取文件访问URL
func (h *FileHandler) GetFileURL(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	// 获取文件ID
	fileIDStr := ctx.Param("id")
	fileID := cast.ToUint(fileIDStr)
	if fileID == 0 {
		logger.WarnContext(reqCtx, "GetFileURL 文件ID无效", "file_id", fileIDStr)
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件ID无效"))
		return
	}

	// 获取URL类型：temporary（临时，默认）或 permanent（长期）
	urlType := ctx.DefaultQuery("type", "temporary")
	if urlType != "temporary" && urlType != "permanent" {
		logger.WarnContext(reqCtx, "GetFileURL URL类型无效", "type", urlType)
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("URL类型无效，只支持 temporary 或 permanent"))
		return
	}

	// 获取过期时间（分钟），仅对临时URL有效
	expiresIn := cast.ToInt(ctx.DefaultQuery("expires_in", "60")) // 默认1小时
	if expiresIn < 1 || expiresIn > 10080 {                       // 最大7天
		expiresIn = 60
	}

	// 创建文件仓库
	fileRepo := model.NewFileRepo(h.DB)

	// 获取文件信息
	fileModel, err := fileRepo.GetByID(reqCtx, fileID)
	if err != nil {
		logger.ErrorContext(reqCtx, "GetFileURL 获取文件信息失败",
			"error", err,
			"file_id", fileID)
		response.Error(ctx, err)
		return
	}

	// 检查文件权限
	if !fileModel.IsPublic {
		// 获取当前用户信息
		userID, exists := ctx.Get("user_id")
		if !exists {
			logger.WarnContext(reqCtx, "GetFileURL 未授权访问私有文件", "file_id", fileID)
			response.Error(ctx, errorx.ErrUnauthorized.WithMsg("需要登录才能访问此文件"))
			return
		}

		// 检查是否是文件上传者或管理员
		if cast.ToInt64(userID) != fileModel.UploadedBy {
			// 检查是否是管理员
			isAdmin, exists := ctx.Get("is_admin")
			if !exists || !cast.ToBool(isAdmin) {
				logger.WarnContext(reqCtx, "GetFileURL 无权限访问文件",
					"file_id", fileID,
					"user_id", userID,
					"uploaded_by", fileModel.UploadedBy)
				response.Error(ctx, errorx.ErrForbidden.WithMsg("无权限访问此文件"))
				return
			}
		}
	}

	var fileURL string
	var expiresAt *int64

	switch urlType {
	case "permanent":
		// 长期URL：仅支持公开文件，直接返回MinIO的公开URL
		if !fileModel.IsPublic {
			logger.WarnContext(reqCtx, "GetFileURL 私有文件不支持长期URL", "file_id", fileID)
			response.Error(ctx, errorx.ErrInvalidParams.WithMsg("私有文件不支持长期URL"))
			return
		}
		// 检查文件是否在public路径下
		if !strings.HasPrefix(fileModel.Path, "public/") {
			logger.WarnContext(reqCtx, "GetFileURL 文件路径不在public目录下", "file_id", fileID, "path", fileModel.Path)
			response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件不在公开目录中"))
			return
		}
		// 返回MinIO的直接访问URL（无需签名，长期有效）
		directURL, err := h.storage.GetDirectURL(reqCtx, fileModel.Path)
		if err != nil {
			logger.ErrorContext(reqCtx, "GetFileURL 获取直接URL失败",
				"error", err,
				"file_id", fileID,
				"path", fileModel.Path)
			response.Error(ctx, errorx.WrapError(err, "生成文件访问链接失败"))
			return
		}
		fileURL = directURL
		// 长期URL无过期时间
	case "temporary":
		// 临时URL：使用预签名URL
		presignedURL, err := h.storage.GetPresignedURL(reqCtx, fileModel.Path, expiresIn)
		if err != nil {
			logger.ErrorContext(reqCtx, "GetFileURL 生成预签名URL失败",
				"error", err,
				"file_id", fileID,
				"path", fileModel.Path)
			response.Error(ctx, errorx.WrapError(err, "生成文件访问链接失败"))
			return
		}
		fileURL = presignedURL
		expTime := time.Now().Add(time.Duration(expiresIn) * time.Minute).Unix()
		expiresAt = &expTime
	}

	// 构建响应数据
	result := map[string]interface{}{
		"file_id":       fileID,
		"file_name":     fileModel.Name,
		"original_name": fileModel.OriginalName,
		"url":           fileURL,
		"type":          urlType,
		"mime_type":     fileModel.MimeType,
		"size":          fileModel.Size,
		"is_public":     fileModel.IsPublic,
	}

	// 添加过期信息（临时和长期URL都有过期时间）
	if expiresAt != nil {
		result["expires_at"] = *expiresAt
		if urlType == "temporary" {
			result["expires_in"] = expiresIn
		} else {
			result["expires_in"] = 7 * 24 * 60 // 7天
		}
	}

	logger.InfoContext(reqCtx, "GetFileURL 获取文件URL成功",
		"file_id", fileID,
		"type", urlType,
		"expires_in", expiresIn)

	response.Success(ctx, result)
}

// getBaseURL 获取基础URL
func (h *FileHandler) getBaseURL(ctx *gin.Context) string {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, ctx.Request.Host)
}

// ServePublicFile 提供公开文件访问（长期URL的实现）
func (h *FileHandler) ServePublicFile(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	// 获取文件ID
	fileIDStr := ctx.Param("id")
	fileID := cast.ToUint(fileIDStr)
	if fileID == 0 {
		logger.WarnContext(reqCtx, "ServePublicFile 文件ID无效", "file_id", fileIDStr)
		ctx.JSON(404, gin.H{"error": "文件不存在"})
		return
	}

	// 创建文件仓库
	fileRepo := model.NewFileRepo(h.DB)

	// 获取文件信息
	fileModel, err := fileRepo.GetByID(reqCtx, fileID)
	if err != nil {
		logger.ErrorContext(reqCtx, "ServePublicFile 获取文件信息失败",
			"error", err,
			"file_id", fileID)
		ctx.JSON(404, gin.H{"error": "文件不存在"})
		return
	}

	// 检查是否是公开文件
	if !fileModel.IsPublic {
		logger.WarnContext(reqCtx, "ServePublicFile 尝试访问私有文件", "file_id", fileID)
		ctx.JSON(403, gin.H{"error": "文件不可公开访问"})
		return
	}

	// 检查文件是否在public路径下
	if !strings.HasPrefix(fileModel.Path, "public/") {
		logger.WarnContext(reqCtx, "ServePublicFile 文件不在public目录下", "file_id", fileID, "path", fileModel.Path)
		ctx.JSON(403, gin.H{"error": "文件不在公开目录中"})
		return
	}

	// 获取MinIO直接访问URL并重定向（避免服务器代理）
	directURL, err := h.storage.GetDirectURL(reqCtx, fileModel.Path)
	if err != nil {
		logger.ErrorContext(reqCtx, "ServePublicFile 获取直接URL失败",
			"error", err,
			"file_id", fileID)
		ctx.JSON(500, gin.H{"error": "文件访问失败"})
		return
	}

	// 302重定向到MinIO直接URL
	ctx.Redirect(302, directURL)

	logger.InfoContext(reqCtx, "ServePublicFile 重定向到直接URL成功",
		"file_id", fileID,
		"direct_url", directURL)
}

// ListFiles 获取文件列表
func (h *FileHandler) ListFiles(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	// 获取查询参数
	page := cast.ToInt(ctx.DefaultQuery("page", "1"))
	pageSize := cast.ToInt(ctx.DefaultQuery("page_size", "20"))
	fileType := ctx.Query("type")
	usage := ctx.Query("usage")

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 创建文件仓库
	fileRepo := model.NewFileRepo(h.DB)

	// 获取文件列表
	files, total, err := fileRepo.ListFiles(reqCtx, page, pageSize, fileType, usage, nil)
	if err != nil {
		logger.ErrorContext(reqCtx, "ListFiles 获取文件列表失败",
			"error", err,
			"page", page,
			"page_size", pageSize)
		response.Error(ctx, err)
		return
	}

	// 构建响应数据
	result := map[string]interface{}{
		"files": files,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	logger.InfoContext(reqCtx, "ListFiles 获取文件列表成功",
		"page", page,
		"page_size", pageSize,
		"total", total)

	response.Success(ctx, result)
}

// GetFileInfo 获取文件信息
func (h *FileHandler) GetFileInfo(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	// 获取文件ID
	fileIDStr := ctx.Param("id")
	fileID := cast.ToUint(fileIDStr)
	if fileID == 0 {
		logger.WarnContext(reqCtx, "GetFileInfo 文件ID无效", "file_id", fileIDStr)
		response.Error(ctx, errorx.ErrInvalidParams.WithMsg("文件ID无效"))
		return
	}

	// 创建文件仓库
	fileRepo := model.NewFileRepo(h.DB)

	// 获取文件信息
	fileModel, err := fileRepo.GetByID(reqCtx, fileID)
	if err != nil {
		logger.ErrorContext(reqCtx, "GetFileInfo 获取文件信息失败",
			"error", err,
			"file_id", fileID)
		response.Error(ctx, err)
		return
	}

	logger.InfoContext(reqCtx, "GetFileInfo 获取文件信息成功",
		"file_id", fileID,
		"file_name", fileModel.Name)

	response.Success(ctx, fileModel)
}
