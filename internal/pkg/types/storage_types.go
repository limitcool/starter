package types

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal StorageType = "local" // 本地文件系统
	StorageTypeS3    StorageType = "s3"    // AWS S3
	StorageTypeOSS   StorageType = "oss"   // 阿里云OSS
)
