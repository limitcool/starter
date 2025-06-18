package dto

// SystemSettingsResponse 系统设置响应
type SystemSettingsResponse struct {
	AppName    string `json:"app_name"`    // 应用名称
	AppVersion string `json:"app_version"` // 应用版本
	AppMode    string `json:"app_mode"`    // 应用模式
}
