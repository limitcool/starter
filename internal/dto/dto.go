package dto

const DefaultPageSize = 20

type Pagination struct {
	Page int64 `form:"page" json:"page"`
	Size int64 `form:"size" json:"size"`
}

// GetSkip 计算需要跳过的文档数
func (p *Pagination) GetSkip() int64 {
	return (p.Page - 1) * p.Size
}

// GetLimit 返回每页的大小
func (p *Pagination) GetLimit() int64 {
	if p.Size <= 0 || p.Size > DefaultPageSize {
		return DefaultPageSize
	}
	return p.Size
}

// total Resp
type ListResponse struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}
