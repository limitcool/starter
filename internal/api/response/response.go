package response

// PageResult 分页结果
type PageResult[T any] struct {
	Total    int64 `json:"total"`     // 总记录数
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 每页大小
	List     T     `json:"list"`      // 数据列表
}

// NewPageResult 创建分页结果
func NewPageResult[T any](list T, total int64, page, pageSize int) *PageResult[T] {
	return &PageResult[T]{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}
}
