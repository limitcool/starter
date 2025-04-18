package model

// SysUserQuery 系统用户查询参数
type SysUserQuery struct {
	Page      int64  `json:"page" form:"page"`           // 页码
	PageSize  int64  `json:"page_size" form:"page_size"` // 每页大小
	Username  string `json:"username" form:"username"`   // 用户名
	Nickname  string `json:"nickname" form:"nickname"`   // 昵称
	Email     string `json:"email" form:"email"`         // 邮箱
	Phone     string `json:"phone" form:"phone"`         // 手机号
	Status    *int8  `json:"status" form:"status"`       // 状态
	OrderBy   string `json:"order_by" form:"order_by"`   // 排序字段
	OrderDesc bool   `json:"order_desc" form:"order_desc"` // 是否降序
}

// Normalize 标准化查询参数
func (q *SysUserQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
}
