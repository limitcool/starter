package vo

// PageResult 分页结果VO
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

// Response 通用响应VO
type Response[T any] struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 消息
	Data    T      `json:"data"`    // 数据
}

// Success 成功响应
func Success[T any](data T) *Response[T] {
	return &Response[T]{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, message string) *Response[any] {
	return &Response[any]{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
