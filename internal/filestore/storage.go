package filestore

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/types"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
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
	Config   Config
	bucket   *blob.Bucket
	s3Client *s3.Client // 用于生成预签名URL
}

// New 创建存储服务
func New(config Config) (*Storage, error) {
	var bucket *blob.Bucket
	var s3Client *s3.Client
	var err error

	ctx := context.Background()

	switch config.Type {
	case types.StorageTypeLocal:
		// 确保本地存储路径存在
		err = os.MkdirAll(config.Path, os.ModePerm)
		if err != nil {
			return nil, errorx.WrapError(err, fmt.Sprintf("创建存储路径失败: %s", config.Path))
		}
		bucket, err = fileblob.OpenBucket(config.Path, nil)
		if err != nil {
			return nil, errorx.WrapError(err, "创建本地存储失败")
		}
	case types.StorageTypeS3:
		// S3配置检查
		if config.AccessKey == "" || config.SecretKey == "" || config.Bucket == "" {
			return nil, errorx.Errorf(errorx.ErrFileStorage, "S3配置不完整")
		}

		// 设置环境变量供gocloud.dev使用
		_ = os.Setenv("AWS_ACCESS_KEY_ID", config.AccessKey)
		_ = os.Setenv("AWS_SECRET_ACCESS_KEY", config.SecretKey)
		_ = os.Setenv("AWS_REGION", config.Region)

		// 构建S3 URL
		s3URL := fmt.Sprintf("s3://%s?region=%s&awssdk=v1", config.Bucket, config.Region)

		// 如果有自定义端点，添加到URL中
		if config.Endpoint != "" {
			s3URL += fmt.Sprintf("&endpoint=%s&disableSSL=true&s3ForcePathStyle=true", config.Endpoint)
		}

		// 使用URL方式打开bucket
		bucket, err = blob.OpenBucket(ctx, s3URL)
		if err != nil {
			return nil, errorx.WrapError(err, "创建S3存储失败")
		}

		// 创建S3客户端用于预签名URL
		s3Client, err = createS3Client(config)
		if err != nil {
			return nil, errorx.WrapError(err, "创建S3客户端失败")
		}
	default:
		return nil, errorx.Errorf(errorx.ErrFileStorage, "不支持的存储类型: %s", config.Type)
	}

	return &Storage{
		Config:   config,
		bucket:   bucket,
		s3Client: s3Client,
	}, nil
}

// createS3Client 创建S3客户端
func createS3Client(config Config) (*s3.Client, error) {
	ctx := context.Background()

	// 创建AWS配置
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			config.AccessKey, config.SecretKey, "")),
		awsconfig.WithRegion(config.Region),
	)
	if err != nil {
		return nil, err
	}

	// 如果有自定义端点，设置它
	if config.Endpoint != "" {
		awsCfg.BaseEndpoint = aws.String(config.Endpoint)
	}

	// 创建S3客户端
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if config.Endpoint != "" {
			o.UsePathStyle = true // MinIO需要使用路径样式
		}
	})

	return client, nil
}

// GetPresignedURL 生成预签名URL（重载方法，兼容旧接口）
func (s *Storage) GetPresignedURL(ctx context.Context, path string, expiresInOrAction interface{}, args ...interface{}) (string, error) {
	// 处理不同的参数组合
	var expiresIn int

	switch v := expiresInOrAction.(type) {
	case int:
		// 新接口：GetPresignedURL(ctx, path, expiresIn)
		expiresIn = v
	case string:
		// 旧接口：GetPresignedURL(ctx, path, action, expiresIn)
		// action参数被忽略，因为我们统一使用GetObject
		if len(args) > 0 {
			if exp, ok := args[0].(int); ok {
				expiresIn = exp
			} else {
				expiresIn = 60 // 默认1小时
			}
		} else {
			expiresIn = 60
		}
	default:
		expiresIn = 60
	}

	if s.Config.Type != types.StorageTypeS3 {
		return "", errorx.Errorf(errorx.ErrFileStorage, "预签名URL仅支持S3存储")
	}

	if s.s3Client == nil {
		return "", errorx.Errorf(errorx.ErrFileStorage, "S3客户端未初始化")
	}

	// 设置过期时间
	duration := time.Duration(expiresIn) * time.Minute

	// 生成GetObject预签名URL
	presigner := s3.NewPresignClient(s.s3Client)
	req, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", errorx.WrapError(err, "生成预签名URL失败")
	}
	return req.URL, nil
}

// GetUploadPresignedURL 生成上传预签名URL
func (s *Storage) GetUploadPresignedURL(ctx context.Context, path, contentType string, expiresIn int) (string, error) {
	if s.Config.Type != types.StorageTypeS3 {
		return "", errorx.Errorf(errorx.ErrFileStorage, "上传预签名URL仅支持S3存储")
	}

	if s.s3Client == nil {
		return "", errorx.Errorf(errorx.ErrFileStorage, "S3客户端未初始化")
	}

	// 设置过期时间
	duration := time.Duration(expiresIn) * time.Minute

	// 生成PutObject预签名URL
	presigner := s3.NewPresignClient(s.s3Client)
	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Config.Bucket),
		Key:         aws.String(path),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", errorx.WrapError(err, "生成上传预签名URL失败")
	}
	return req.URL, nil
}

// GetDirectURL 获取直接访问URL（用于公开文件的长期链接）
func (s *Storage) GetDirectURL(ctx context.Context, path string) (string, error) {
	switch s.Config.Type {
	case types.StorageTypeLocal:
		// 本地文件返回HTTP URL
		return fmt.Sprintf("%s/%s", strings.TrimRight(s.Config.URL, "/"), strings.TrimLeft(path, "/")), nil
	case types.StorageTypeS3:
		// S3直接URL（无签名，仅适用于公开文件）
		if s.Config.Endpoint != "" {
			// MinIO或自定义端点
			return fmt.Sprintf("%s/%s/%s",
				strings.TrimRight(s.Config.Endpoint, "/"),
				s.Config.Bucket,
				strings.TrimLeft(path, "/")), nil
		}
		// AWS S3
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
			s.Config.Bucket,
			s.Config.Region,
			strings.TrimLeft(path, "/")), nil
	default:
		return "", errorx.Errorf(errorx.ErrFileStorage, "不支持的存储类型: %s", s.Config.Type)
	}
}

// Put 上传文件
func (s *Storage) Put(ctx context.Context, path string, reader io.Reader) error {
	writer, err := s.bucket.NewWriter(ctx, path, nil)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("创建写入器失败: %s", path))
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("上传文件失败: %s", path))
	}

	return nil
}

// PutWithACL 上传文件并设置访问权限
func (s *Storage) PutWithACL(ctx context.Context, path string, reader io.Reader, isPublic bool, contentType string) error {
	switch s.Config.Type {
	case types.StorageTypeLocal:
		// 本地存储直接调用Put方法
		return s.Put(ctx, path, reader)
	case types.StorageTypeS3:
		// S3存储使用原生客户端上传以设置ACL
		if s.s3Client == nil {
			return errorx.Errorf(errorx.ErrFileStorage, "S3客户端未初始化")
		}

		// 读取文件内容到内存（对于大文件可以考虑流式处理）
		content, err := io.ReadAll(reader)
		if err != nil {
			return errorx.WrapError(err, "读取文件内容失败")
		}

		// 设置上传参数
		input := &s3.PutObjectInput{
			Bucket:      aws.String(s.Config.Bucket),
			Key:         aws.String(path),
			Body:        strings.NewReader(string(content)),
			ContentType: aws.String(contentType),
		}

		// 设置对象元数据标记是否公开（MinIO可能不支持ACL，我们使用元数据标记）
		if input.Metadata == nil {
			input.Metadata = make(map[string]string)
		}
		if isPublic {
			input.Metadata["x-amz-meta-public"] = "true"
		} else {
			input.Metadata["x-amz-meta-public"] = "false"
		}

		// 上传文件
		_, err = s.s3Client.PutObject(ctx, input)
		if err != nil {
			return errorx.WrapError(err, fmt.Sprintf("上传文件失败: %s", path))
		}

		return nil
	default:
		return errorx.Errorf(errorx.ErrFileStorage, "不支持的存储类型: %s", s.Config.Type)
	}
}

// Get 获取文件
func (s *Storage) Get(ctx context.Context, path string) (*os.File, error) {
	// 对于gocloud.dev，我们需要返回一个临时文件
	reader, err := s.bucket.NewReader(ctx, path, nil)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取文件失败: %s", path))
	}
	defer reader.Close()

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "storage_*")
	if err != nil {
		return nil, errorx.WrapError(err, "创建临时文件失败")
	}

	// 复制内容到临时文件
	_, err = io.Copy(tempFile, reader)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, errorx.WrapError(err, "复制文件内容失败")
	}

	// 重置文件指针到开始
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, errorx.WrapError(err, "重置文件指针失败")
	}

	return tempFile, nil
}

// GetStream 获取文件流
func (s *Storage) GetStream(ctx context.Context, path string) (io.ReadCloser, error) {
	reader, err := s.bucket.NewReader(ctx, path, nil)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取文件流失败: %s", path))
	}
	return reader, nil
}

// Delete 删除文件
func (s *Storage) Delete(ctx context.Context, path string) error {
	err := s.bucket.Delete(ctx, path)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除文件失败: %s", path))
	}
	return nil
}

// Exists 检查文件是否存在
func (s *Storage) Exists(ctx context.Context, path string) (bool, error) {
	exists, err := s.bucket.Exists(ctx, path)
	if err != nil {
		return false, errorx.WrapError(err, fmt.Sprintf("检查文件是否存在失败: %s", path))
	}
	return exists, nil
}

// GetURL 获取文件访问URL
func (s *Storage) GetURL(ctx context.Context, path string) (string, error) {
	switch s.Config.Type {
	case types.StorageTypeLocal:
		// 本地文件返回HTTP URL
		return fmt.Sprintf("%s/%s", strings.TrimRight(s.Config.URL, "/"), strings.TrimLeft(path, "/")), nil
	case types.StorageTypeS3:
		// 对于S3，根据路径判断是否为公开文件
		if strings.HasPrefix(path, "public/") {
			// 公开文件返回直接URL
			return s.GetDirectURL(ctx, path)
		} else {
			// 私有文件返回预签名URL（默认1小时）
			return s.GetPresignedURL(ctx, path, 60)
		}
	default:
		return "", errorx.Errorf(errorx.ErrFileStorage, "不支持的存储类型: %s", s.Config.Type)
	}
}

// Close 关闭存储服务
func (s *Storage) Close() error {
	if s.bucket != nil {
		return s.bucket.Close()
	}
	return nil
}

// GetEndpoint 获取存储端点
func (s *Storage) GetEndpoint() string {
	return s.Config.Endpoint
}

// GetEndpointWithContext 使用上下文获取存储端点
func (s *Storage) GetEndpointWithContext(ctx context.Context) string {
	return s.Config.Endpoint
}

// GeneratePath 生成存储路径
func GeneratePath(ctx context.Context, baseDir, fileName string) string {
	// 清理文件名，移除危险字符
	fileName = filepath.Base(fileName)

	// 文件路径格式：baseDir/fileName
	return filepath.Join(baseDir, fileName)
}

// MimeTypes 常见文件类型的MIME类型映射
var MimeTypes = map[string]string{
	// 图片类型
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".ico":  "image/x-icon",

	// 文档类型
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".txt":  "text/plain",
	".rtf":  "application/rtf",

	// 音频类型
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".flac": "audio/flac",
	".aac":  "audio/aac",
	".ogg":  "audio/ogg",

	// 视频类型
	".mp4":  "video/mp4",
	".avi":  "video/x-msvideo",
	".mov":  "video/quicktime",
	".wmv":  "video/x-ms-wmv",
	".flv":  "video/x-flv",
	".webm": "video/webm",
	".mkv":  "video/x-matroska",

	// 压缩文件
	".zip": "application/zip",
	".rar": "application/vnd.rar",
	".7z":  "application/x-7z-compressed",
	".tar": "application/x-tar",
	".gz":  "application/gzip",

	// 其他常见类型
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

// NewWithConfig 使用配置创建存储服务
func NewWithConfig(cfg configs.Config) (*Storage, error) {
	if !cfg.Storage.Enabled {
		return nil, nil
	}

	storageConfig := Config{
		Type: cfg.Storage.Type,
	}

	switch cfg.Storage.Type {
	case types.StorageTypeLocal:
		storageConfig.Path = cfg.Storage.Local.Path
		storageConfig.URL = cfg.Storage.Local.URL
	case types.StorageTypeS3:
		storageConfig.AccessKey = cfg.Storage.S3.AccessKey
		storageConfig.SecretKey = cfg.Storage.S3.SecretKey
		storageConfig.Region = cfg.Storage.S3.Region
		storageConfig.Bucket = cfg.Storage.S3.Bucket
		storageConfig.Endpoint = cfg.Storage.S3.Endpoint
	case types.StorageTypeOSS:
		storageConfig.AccessKey = cfg.Storage.OSS.AccessKey
		storageConfig.SecretKey = cfg.Storage.OSS.SecretKey
		storageConfig.Region = cfg.Storage.OSS.Region
		storageConfig.Bucket = cfg.Storage.OSS.Bucket
		storageConfig.Endpoint = cfg.Storage.OSS.Endpoint
	default:
		return nil, errorx.Errorf(errorx.ErrFileStorage, "不支持的存储类型: %s", cfg.Storage.Type)
	}

	return New(storageConfig)
}
