package dto

// total Resp
type ListResponse[T any] struct {
	Total int64 `json:"total"`
	Items T     `json:"items"`
}

// BaseDTO 基础数据传输对象
type BaseDTO struct {
	ID uint `json:"id" form:"id"` // ID
}

// BaseQuery 基础查询对象
type BaseQuery struct {
	Keyword   string `json:"keyword" form:"keyword"`       // 关键字
	Status    *int   `json:"status" form:"status"`         // 状态
	StartTime string `json:"start_time" form:"start_time"` // 开始时间
	EndTime   string `json:"end_time" form:"end_time"`     // 结束时间
	CreateBy  string `json:"create_by" form:"create_by"`   // 创建人
}
