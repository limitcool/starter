package v1

import "time"

// OperationLogQuery 操作日志查询请求
type OperationLogQuery struct {
	PageRequest
	BaseQuery
	UserType string `json:"user_type" form:"user_type"` // 用户类型
	Username string `json:"username" form:"username"`   // 用户名
	Module   string `json:"module" form:"module"`       // 操作模块
	Action   string `json:"action" form:"action"`       // 操作类型
	IP       string `json:"ip" form:"ip"`               // IP地址
	// 这里的StartTime和EndTime是time.Time类型，用于精确时间查询
	// 覆盖了BaseQuery中的string类型字段
	StartTime *time.Time `json:"start_time" form:"start_time"` // 开始时间
	EndTime   *time.Time `json:"end_time" form:"end_time"`     // 结束时间
}

// OperationLogBatchDeleteRequest 批量删除操作日志请求
type OperationLogBatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"` // ID列表
}
