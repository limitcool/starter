package dto

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileID      string   `json:"file_id"`      // 文件ID
	UploadURL   string   `json:"upload_url"`   // 上传URL
	Method      string   `json:"method"`       // HTTP方法
	ExpiresIn   int      `json:"expires_in"`   // 过期时间（分钟）
	StorageType string   `json:"storage_type"` // 存储类型
	Usage       string   `json:"usage"`        // 文件用途
	PathInfo    PathInfo `json:"path_info"`    // 路径信息
}

// PathInfo 路径信息
type PathInfo struct {
	Category string `json:"category"` // 分类
	Path     string `json:"path"`     // 路径
}

// FileDownloadResponse 文件下载响应
type FileDownloadResponse struct {
	FileID      string `json:"file_id"`      // 文件ID
	Filename    string `json:"filename"`     // 文件名
	DownloadURL string `json:"download_url"` // 下载URL
	IsPublic    bool   `json:"is_public"`    // 是否公开
	Size        int64  `json:"size"`         // 文件大小
	StorageType string `json:"storage_type"` // 存储类型
}

// FileUploadCompleteResponse 文件上传完成响应
type FileUploadCompleteResponse struct {
	FileID      string `json:"file_id"`      // 文件ID
	Filename    string `json:"filename"`     // 文件名
	DownloadURL string `json:"download_url"` // 下载URL
	Size        int64  `json:"size"`         // 文件大小
	StorageType string `json:"storage_type"` // 存储类型
	IsPublic    bool   `json:"is_public"`    // 是否公开
	Usage       string `json:"usage"`        // 文件用途
}

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	IsPublic    bool   `json:"is_public"`
	Usage       string `json:"usage" binding:"required"` // avatar, banner, document, etc.
	Size        int64  `json:"size,omitempty"`           // 文件大小（可选，用于预验证）
}

// FileConfirmRequest 文件确认上传请求
type FileConfirmRequest struct {
	FileID string `json:"file_id" binding:"required"`
	Size   int64  `json:"size" binding:"required"`
}
