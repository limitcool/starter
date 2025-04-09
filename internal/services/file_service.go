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
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/storage"
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
func (s *FileService) UploadFile(c *gin.Context, fileHeader *multipart.FileHeader, fileType, uploader string) (*model.File, error) {
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

	// 记录到数据库
	fileModel := &model.File{
		Name:         fileName,
		OriginalName: originalName,
		Path:         storagePath,
		URL:          fileURL,
		Type:         fileType,
		Size:         size,
		MimeType:     mimeType,
		Extension:    ext,
		StorageType:  string(s.storage.Config.Type),
		UploadedAt:   time.Now(),
		Status:       1, // 状态正常
	}

	// 如果有上传者ID，转换并设置
	if uploader != "" {
		var uploaderID uint
		fmt.Sscanf(uploader, "%d", &uploaderID)
		if uploaderID > 0 {
			fileModel.UploadedBy = uploaderID
		}
	}

	if err := db.Create(fileModel).Error; err != nil {
		return nil, errorx.ErrDatabase.WithError(err)
	}

	return fileModel, nil
}

// GetFile 获取文件信息
func (s *FileService) GetFile(id string) (*model.File, error) {
	var file model.File
	if err := db.First(&file, id).Error; err != nil {
		return nil, errorx.ErrNotFound.WithError(err)
	}
	return &file, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(id string) error {
	var file model.File
	if err := db.First(&file, id).Error; err != nil {

		return errorx.ErrNotFound.WithError(err)

	}

	// 删除存储中的文件
	if err := s.storage.Delete(file.Path); err != nil {
		return errorx.ErrFileStroage.WithError(err)
	}

	// 从数据库中删除记录
	if err := db.Delete(&file).Error; err != nil {
		return errorx.ErrDatabase.WithError(err)
	}

	return nil
}

// 获取文件内容
func (s *FileService) GetFileContent(id string) (io.ReadCloser, string, error) {
	var file model.File
	if err := db.First(&file, id).Error; err != nil {
		return nil, "", errorx.ErrNotFound.WithError(err)
	}

	// 获取文件流
	fileStream, err := s.storage.GetStream(file.Path)
	if err != nil {
		return nil, "", errorx.ErrFileStroage.WithError(err)
	}

	return fileStream, file.MimeType, nil
}

// UpdateUserAvatar 更新用户头像
func (s *FileService) UpdateUserAvatar(userID string, fileHeader *multipart.FileHeader) (string, error) {
	// 上传头像文件
	file, err := s.UploadFile(nil, fileHeader, model.FileTypeAvatar, userID)
	if err != nil {
		return "", err
	}

	// 查找用户
	user := model.User{}
	if err := db.First(&user, userID).Error; err != nil {
		return "", errorx.ErrNotFound.WithError(err)
	}

	// 更新用户头像
	user.Avatar = fmt.Sprintf("%d", file.ID)
	if err := db.Save(&user).Error; err != nil {
		return "", errorx.ErrDatabase.WithError(err)
	}

	return file.URL, nil
}

// UpdateSysUserAvatar 更新系统用户头像
func (s *FileService) UpdateSysUserAvatar(userID string, fileHeader *multipart.FileHeader) (string, error) {
	// 上传头像文件
	file, err := s.UploadFile(nil, fileHeader, model.FileTypeAvatar, userID)
	if err != nil {
		return "", err
	}

	// 查找系统用户
	sysUser := model.SysUser{}
	if err := db.First(&sysUser, userID).Error; err != nil {
		return "", errorx.ErrNotFound.WithError(err)
	}

	// 更新用户头像
	sysUser.Avatar = fmt.Sprintf("%d", file.ID)
	if err := db.Save(&sysUser).Error; err != nil {
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
	case model.FileTypeAvatar:
		baseDir = "avatars"
	case model.FileTypeDocument:
		baseDir = "documents"
	case model.FileTypeImage:
		baseDir = "images"
	case model.FileTypeVideo:
		baseDir = "videos"
	case model.FileTypeAudio:
		baseDir = "audios"
	default:
		baseDir = "others"
	}

	// 添加日期子目录
	dateDir := time.Now().Format("2006/01/02")
	return filepath.Join(baseDir, dateDir, fileName)
}

// 检查文件类型是否允许
func isAllowedFileType(ext string, fileType string) bool {
	ext = strings.ToLower(ext)
	switch fileType {
	case model.FileTypeAvatar, model.FileTypeImage:
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
	case model.FileTypeAvatar:
		return size <= 2*MB // 头像最大2MB
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
