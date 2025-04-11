package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/storage"
	"github.com/limitcool/starter/internal/storage/sqldb"
)

// FileService 文件服务
type FileService struct {
	storage *storage.Storage
}

// NewFileService 创建文件服务
func NewFileService(storage *storage.Storage) *FileService {
	return &FileService{
		storage: storage,
	}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(c *gin.Context, fileHeader *multipart.FileHeader, fileType string, uploaderID int64, userType enum.UserType, isPublic bool) (*model.File, error) {
	if fileHeader == nil {
		return nil, errorx.ErrInvalidParams
	}

	// 获取文件基本信息
	originalName := fileHeader.Filename
	ext := strings.ToLower(filepath.Ext(originalName))
	mimeType := storage.GetMimeType(originalName)
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
	err = s.storage.Put(storagePath, file)
	if err != nil {
		return nil, errorx.ErrFileStroage.WithError(err)
	}

	// 获取文件访问URL
	fileURL, err := s.storage.GetURL(storagePath)
	if err != nil {
		return nil, errorx.ErrFileStroage.WithError(err)
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

	if err := sqldb.Instance().DB().Create(fileModel).Error; err != nil {
		return nil, errorx.ErrDatabase.WithError(err)
	}

	return fileModel, nil
}

// GetFile 获取文件信息
func (s *FileService) GetFile(id string, currentUserID int64, currentUserType enum.UserType) (*model.File, error) {
	var file model.File
	if err := sqldb.Instance().DB().First(&file, id).Error; err != nil {
		return nil, errorx.ErrNotFound.WithError(err)
	}

	// 特殊情况：未登录用户（ID=0）只能访问公开文件、图片和文档类型的文件
	if currentUserID == 0 {
		if file.IsPublic || file.Type == model.FileTypeImage || file.Type == model.FileTypeDocument {
			return &file, nil
		}
		return nil, errorx.ErrAccessDenied.WithMsg("您无权访问此文件")
	}

	// 系统用户可以访问所有文件
	if currentUserType == enum.UserTypeSysUser {
		return &file, nil
	}

	// 公开文件任何人都可以访问
	if file.IsPublic {
		return &file, nil
	}

	// 普通用户可以访问自己上传的文件
	if currentUserType == enum.UserTypeUser && uint8(enum.UserTypeUser) == file.UploadedByType && file.UploadedBy == currentUserID {
		return &file, nil
	}

	// 普通用户可以访问特定类型的公开文件（如图片和文档）
	if currentUserType == enum.UserTypeUser && (file.Type == model.FileTypeImage || file.Type == model.FileTypeDocument) {
		// 这里可以添加额外的权限逻辑
		return &file, nil
	}

	return nil, errorx.ErrAccessDenied.WithMsg("您无权访问此文件")
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(id string, currentUserID int64, currentUserType enum.UserType) error {
	var file model.File
	if err := sqldb.Instance().DB().First(&file, id).Error; err != nil {
		return errorx.ErrNotFound.WithError(err)
	}

	// 权限检查：系统用户可以删除所有文件，普通用户只能删除自己上传的文件
	if currentUserType == enum.UserTypeUser && (file.UploadedByType != uint8(enum.UserTypeUser) || file.UploadedBy != currentUserID) {
		return errorx.ErrAccessDenied.WithMsg("您无权删除此文件")
	}

	// 删除存储中的文件
	if err := s.storage.Delete(file.Path); err != nil {
		return errorx.ErrFileStroage.WithError(err)
	}

	// 从数据库中删除记录
	if err := sqldb.Instance().DB().Delete(&file).Error; err != nil {
		return errorx.ErrDatabase.WithError(err)
	}

	return nil
}

// 获取文件内容
func (s *FileService) GetFileContent(id string, currentUserID int64, currentUserType enum.UserType) (io.ReadCloser, string, error) {
	// 先验证权限
	file, err := s.GetFile(id, currentUserID, currentUserType)
	if err != nil {
		return nil, "", err
	}

	// 获取文件流
	fileStream, err := s.storage.GetStream(file.Path)
	if err != nil {
		return nil, "", errorx.ErrFileStroage.WithError(err)
	}

	return fileStream, file.MimeType, nil
}

// UpdateUserAvatar 更新用户头像
func (s *FileService) UpdateUserAvatar(userID int64, fileHeader *multipart.FileHeader) (string, error) {
	// 上传头像文件（图片类型，头像用途）
	file, err := s.UploadFile(nil, fileHeader, model.FileTypeImage, userID, enum.UserTypeUser, false)
	if err != nil {
		return "", err
	}

	// 设置文件用途为头像
	file.Usage = model.FileUsageAvatar
	if err := sqldb.Instance().DB().Save(file).Error; err != nil {
		return "", errorx.ErrDatabase.WithError(err)
	}

	// 查找用户
	user := model.User{}
	if err := sqldb.Instance().DB().First(&user, userID).Error; err != nil {
		return "", errorx.ErrNotFound.WithError(err)
	}

	// 更新用户头像
	user.AvatarFileID = file.ID
	if err := sqldb.Instance().DB().Save(&user).Error; err != nil {
		return "", errorx.ErrDatabase.WithError(err)
	}

	return file.URL, nil
}

// UpdateSysUserAvatar 更新系统用户头像
func (s *FileService) UpdateSysUserAvatar(userID int64, fileHeader *multipart.FileHeader) (string, error) {
	// 上传头像文件（图片类型，头像用途）
	file, err := s.UploadFile(nil, fileHeader, model.FileTypeImage, userID, enum.UserTypeSysUser, false)
	if err != nil {
		return "", err
	}

	// 设置文件用途为头像
	file.Usage = model.FileUsageAvatar
	if err := sqldb.Instance().DB().Save(file).Error; err != nil {
		return "", errorx.ErrDatabase.WithError(err)
	}

	// 查找系统用户
	sysUser := model.SysUser{}
	if err := sqldb.Instance().DB().First(&sysUser, userID).Error; err != nil {
		return "", errorx.ErrNotFound.WithError(err)
	}

	// 更新用户头像
	sysUser.AvatarFileID = file.ID
	if err := sqldb.Instance().DB().Save(&sysUser).Error; err != nil {
		return "", errorx.ErrDatabase.WithError(err)
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
