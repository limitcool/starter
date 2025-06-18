package dto

// PageRequest 分页请求基础结构
type PageRequest struct {
	Page     int    `form:"page" json:"page"`           // 页码
	PageSize int    `form:"page_size" json:"page_size"` // 每页大小
	SortBy   string `form:"sort_by" json:"sort_by"`     // 排序字段
	SortDesc bool   `form:"sort_desc" json:"sort_desc"` // 是否降序
}

// Normalize 标准化分页请求
func (p *PageRequest) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// GetSortDirection 获取排序方向
func (p *PageRequest) GetSortDirection() string {
	if p.SortDesc {
		return "DESC"
	}
	return "ASC"
}
