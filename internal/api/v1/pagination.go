package v1

// 默认分页大小
const DefaultPageSize = 20

// PageRequest 分页请求
type PageRequest struct {
	Page     int    `form:"page" json:"page"`           // 页码
	PageSize int    `form:"size" json:"size"`           // 每页大小
	SortBy   string `form:"sort_by" json:"sort_by"`     // 排序字段
	SortDesc bool   `form:"sort_desc" json:"sort_desc"` // 是否降序排序
}

// GetDefaultPageRequest 获取默认分页请求
func GetDefaultPageRequest() *PageRequest {
	return &PageRequest{
		Page:     1,
		PageSize: DefaultPageSize,
		SortBy:   "id",
		SortDesc: true,
	}
}

// Normalize 标准化分页参数
func (p *PageRequest) Normalize() *PageRequest {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	} else if p.PageSize > 100 {
		// 限制最大每页大小为100
		p.PageSize = 100
	}

	if p.SortBy == "" {
		p.SortBy = "id"
	}

	return p
}

// GetSortDirection 获取排序方向
func (p *PageRequest) GetSortDirection() string {
	if p.SortDesc {
		return "desc"
	}
	return "asc"
}

// GetOffset 获取偏移量
func (p *PageRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数
func (p *PageRequest) GetLimit() int {
	return p.PageSize
}
