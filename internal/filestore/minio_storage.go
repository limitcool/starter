package filestore

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// MinIOStorage MinIO存储实现
type MinIOStorage struct {
	client   *s3.Client
	bucket   string
	endpoint string
	region   string
}

// MinIOConfig MinIO配置
type MinIOConfig struct {
	Endpoint  string
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
}

// NewMinIOStorage 创建MinIO存储实例
func NewMinIOStorage(config MinIOConfig) (*MinIOStorage, error) {
	ctx := context.Background()

	// 创建AWS配置
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			config.AccessKey, config.SecretKey, "")),
		awsconfig.WithRegion(config.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("创建AWS配置失败: %w", err)
	}

	// 设置自定义端点
	if config.Endpoint != "" {
		awsCfg.BaseEndpoint = aws.String(config.Endpoint)
	}

	// 创建S3客户端
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if config.Endpoint != "" {
			o.UsePathStyle = true // MinIO需要使用路径样式
		}
	})

	return &MinIOStorage{
		client:   client,
		bucket:   config.Bucket,
		endpoint: config.Endpoint,
		region:   config.Region,
	}, nil
}

// GetUploadURL 获取上传预签名URL
func (m *MinIOStorage) GetUploadURL(ctx context.Context, filePath string, contentType string, isPublic bool) (string, string, error) {
	fullPath := m.BuildFullPath(filePath, isPublic)

	// 生成PutObject预签名URL
	presigner := s3.NewPresignClient(m.client)
	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(m.bucket),
		Key:         aws.String(fullPath),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute // 15分钟有效期
	})
	if err != nil {
		return "", "", fmt.Errorf("生成上传预签名URL失败: %w", err)
	}

	return req.URL, "PUT", nil
}

// GetDownloadURL 获取下载URL
func (m *MinIOStorage) GetDownloadURL(ctx context.Context, filePath string, isPublic bool) (string, error) {
	fullPath := m.BuildFullPath(filePath, isPublic)

	if isPublic {
		// 公开文件返回直接URL（无需签名）
		return fmt.Sprintf("%s/%s/%s",
			strings.TrimRight(m.endpoint, "/"),
			m.bucket,
			strings.TrimLeft(fullPath, "/")), nil
	} else {
		// 私有文件返回预签名URL
		presigner := s3.NewPresignClient(m.client)
		req, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(m.bucket),
			Key:    aws.String(fullPath),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 1 * time.Hour // 1小时有效期
		})
		if err != nil {
			return "", fmt.Errorf("生成下载预签名URL失败: %w", err)
		}
		return req.URL, nil
	}
}

// UploadFile 直接上传文件（MinIO不推荐使用，应该使用预签名URL）
func (m *MinIOStorage) UploadFile(ctx context.Context, filePath string, reader io.Reader, isPublic bool) error {
	fullPath := m.BuildFullPath(filePath, isPublic)

	// 读取文件内容
	content, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("读取文件内容失败: %w", err)
	}

	// 上传到MinIO
	_, err = m.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(fullPath),
		Body:   strings.NewReader(string(content)),
	})
	if err != nil {
		return fmt.Errorf("上传文件到MinIO失败: %w", err)
	}

	return nil
}

// FileExists 检查文件是否存在
func (m *MinIOStorage) FileExists(ctx context.Context, filePath string, isPublic bool) (bool, error) {
	fullPath := m.BuildFullPath(filePath, isPublic)

	_, err := m.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(fullPath),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// DeleteFile 删除文件
func (m *MinIOStorage) DeleteFile(ctx context.Context, filePath string, isPublic bool) error {
	fullPath := m.BuildFullPath(filePath, isPublic)

	_, err := m.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(fullPath),
	})
	return err
}

// GetStorageType 获取存储类型
func (m *MinIOStorage) GetStorageType() string {
	return "minio"
}

// BuildFullPath 构建完整路径（包含public/private前缀）
func (m *MinIOStorage) BuildFullPath(filePath string, isPublic bool) string {
	if isPublic {
		return fmt.Sprintf("public/%s", filePath)
	}
	return fmt.Sprintf("private/%s", filePath)
}
