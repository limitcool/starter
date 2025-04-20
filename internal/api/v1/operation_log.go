package v1

import "time"

// OperationLogQuery 操作日志查询参数
type OperationLogQuery struct {
	PageRequest
	Username  string     `form:"username" json:"username"`   // 用户名
	UserType  string     `form:"user_type" json:"user_type"` // 用户类型
	Module    string     `form:"module" json:"module"`       // 模块
	Action    string     `form:"action" json:"action"`       // 操作
	IP        string     `form:"ip" json:"ip"`               // IP地址
	StartTime *time.Time `form:"start_time" json:"start_time"`
	EndTime   *time.Time `form:"end_time" json:"end_time"`
}

// OperationLogBatchDeleteRequest 批量删除操作日志请求
type OperationLogBatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}
