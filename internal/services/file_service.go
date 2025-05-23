package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/limitcool/starter/internal/filestore"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/repository"
)

// FileService 文件服务
type FileService struct {
	storage       *filestore.Storage
	fileRepo      *repository.FileRepo
	userRepo      *repository.UserRepo
	adminUserRepo *repository.AdminUserRepo
}

// NewFileService 创建文件服务
func NewFileService(storage *filestore.Storage, fileRepo *repository.FileRepo, userRepo *repository.UserRepo, adminUserRepo *repository.AdminUserRepo) *FileService {
	return &FileService{
		storage:       storage,
		fileRepo:      fileRepo,
		userRepo:      userRepo,
		adminUserRepo: adminUserRepo,
	}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, fileType string, uploaderID int64, userType enum.UserType, isPublic bool) (*model.File, error) {
	if fileHeader == nil {
		return nil, errorx.ErrInvalidParams
	}

	// 获取文件基本信息
	originalName := fileHeader.Filename
	ext := strings.ToLower(filepath.Ext(originalName))
	mimeType := filestore.GetMimeType(originalName)
	size := fileHeader.Size

	// 验证文件类型
	if !isAllowedFileType(ext, fileType) {
		return nil, errorx.ErrInvalidParams
	}

	// 验证文件大小
	if !isAllowedFileSize(size, fileType) {
		return nil, errorx.ErrInvalidParams.WithMsg("文件大小超出限制")
	}

	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, errorx.ErrInvalidParams.WithError(err)
	}
	defer file.Close()

	// 生成文件名和存储路径
	fileName := generateFileName(originalName)
	storagePath := s.getStoragePath(fileType, fileName)

	// 上传文件到存储
	err = s.storage.Put(ctx, storagePath, file)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("上传文件到存储失败: %s", storagePath))
	}

	// 获取文件访问URL
	fileURL, err := s.storage.GetURL(ctx, storagePath)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取文件访问URL失败: %s", storagePath))
	}

	// 如果未指定用户类型，默认为系统用户
	if userType == 0 {
		userType = enum.UserTypeSysUser
	}

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
		StorageType:    string(s.storage.Config.Type),
		UploadedBy:     uploaderID,
		UploadedByType: uint8(userType),
		UploadedAt:     time.Now(),
		Status:         1,        // 状态正常
		IsPublic:       isPublic, // 设置是否公开
	}

	if err := s.fileRepo.Create(ctx, fileModel); err != nil {
		return nil, errorx.ErrDatabase.WithError(err)
	}

	return fileModel, nil
}

// GetFile 获取文件信息
func (s *FileService) GetFile(ctx context.Context, id string, currentUserID int64, currentUserType enum.UserType) (*model.File, error) {
	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 特殊情况：未登录用户（ID=0）只能访问公开文件、图片和文档类型的文件
	if currentUserID == 0 {
		if file.IsPublic || file.Type == model.FileTypeImage || file.Type == model.FileTypeDocument {
			return file, nil
		}
		return nil, errorx.Errorf(errorx.ErrAccessDenied, "您无权访问此文件")
	}

	// 系统用户和管理员用户可以访问所有文件
	if currentUserType == enum.UserTypeSysUser || currentUserType == enum.UserTypeAdminUser {
		return file, nil
	}

	// 公开文件任何人都可以访问
	if file.IsPublic {
		return file, nil
	}

	// 普通用户可以访问自己上传的文件
	if currentUserType == enum.UserTypeUser && uint8(enum.UserTypeUser) == file.UploadedByType && file.UploadedBy == currentUserID {
		return file, nil
	}

	// 普通用户可以访问特定类型的公开文件（如图片和文档）
	if currentUserType == enum.UserTypeUser && (file.Type == model.FileTypeImage || file.Type == model.FileTypeDocument) {
		// 这里可以添加额外的权限逻辑
		return file, nil
	}

	return nil, errorx.Errorf(errorx.ErrAccessDenied, "您无权访问此文件")
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(ctx context.Context, id string, currentUserID int64, currentUserType enum.UserType) error {
	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 权限检查：系统用户和管理员用户可以删除所有文件，普通用户只能删除自己上传的文件
	if currentUserType == enum.UserTypeUser && (file.UploadedByType != uint8(enum.UserTypeUser) || file.UploadedBy != currentUserID) {
		return errorx.Errorf(errorx.ErrAccessDenied, "您无权删除此文件")
	}

	// 删除存储中的文件
	if err := s.storage.Delete(ctx, file.Path); err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除存储中的文件失败: %s", file.Path))
	}

	// 从数据库中删除记录
	if err := s.fileRepo.Delete(ctx, id); err != nil {
		return errorx.WrapError(err, fmt.Sprintf("从数据库中删除文件记录失败: id=%s", id))
	}

	return nil
}

// GetFileContent 获取文件内容
// 注意：调用者必须负责关闭返回的io.ReadCloser，建议使用defer fileStream.Close()
// 示例：
//
//	fileStream, contentType, err := fileService.GetFileContent(ctx, id, userID, userType)
//	if err != nil {
//		// 处理错误
//	}
//	defer fileStream.Close() // 重要：确保文件流被关闭
//	// 使用fileStream...
func (s *FileService) GetFileContent(ctx context.Context, id string, currentUserID int64, currentUserType enum.UserType) (io.ReadCloser, string, error) {
	// 先验证权限
	file, err := s.GetFile(ctx, id, currentUserID, currentUserType)
	if err != nil {
		return nil, "", err
	}

	// 获取文件流
	fileStream, err := s.storage.GetStream(ctx, file.Path)
	if err != nil {
		return nil, "", errorx.WrapError(err, fmt.Sprintf("获取文件流失败: %s", file.Path))
	}

	return fileStream, file.MimeType, nil
}

// UpdateUserAvatar 更新用户头像
func (s *FileService) UpdateUserAvatar(ctx context.Context, userID int64, fileHeader *multipart.FileHeader) (string, error) {
	// 上传头像文件（图片类型，头像用途）
	file, err := s.UploadFile(ctx, fileHeader, model.FileTypeImage, userID, enum.UserTypeUser, false)
	if err != nil {
		return "", err
	}

	// 设置文件用途为头像
	if err := s.fileRepo.UpdateFileUsage(ctx, file, model.FileUsageAvatar); err != nil {
		return "", errorx.ErrDatabase.WithError(err)
	}

	// 更新用户头像
	if err := s.userRepo.UpdateAvatar(ctx, userID, file.ID); err != nil {
		return "", err
	}

	return file.URL, nil
}

// UpdateAdminUserAvatar 更新管理员用户头像
func (s *FileService) UpdateAdminUserAvatar(ctx context.Context, userID int64, fileHeader *multipart.FileHeader) (string, error) {
	// 上传头像文件（图片类型，头像用途）
	file, err := s.UploadFile(ctx, fileHeader, model.FileTypeImage, userID, enum.UserTypeAdminUser, false)
	if err != nil {
		return "", err
	}

	// 设置文件用途为头像
	if err := s.fileRepo.UpdateFileUsage(ctx, file, model.FileUsageAvatar); err != nil {
		return "", errorx.ErrDatabase.WithError(err)
	}

	// 更新管理员用户头像
	if err := s.adminUserRepo.UpdateAvatar(ctx, userID, file.ID); err != nil {
		return "", err
	}

	return file.URL, nil
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
func (s *FileService) getStoragePath(fileType, fileName string) string {
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
