package dto

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status string `json:"status"` // 状态
}

// DeleteResponse 删除响应
type DeleteResponse struct {
	Message string `json:"message"` // 消息
}

// MessageResponse 通用消息响应
type MessageResponse struct {
	Message string `json:"message"` // 消息
}
