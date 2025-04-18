package repository

import (
	"time"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/options"
	"gorm.io/gorm"
)

// PageResult 分页结果
type PageResult struct {
	Total    int64                `json:"total"`     // 总记录数
	Page     int                  `json:"page"`      // 当前页码
	PageSize int                  `json:"page_size"` // 每页大小
	List     []model.OperationLog `json:"list"`      // 数据列表
}

// OperationLogQuery 操作日志查询参数
type OperationLogQuery struct {
	Page      int        `json:"page" form:"page"`           // 页码
	PageSize  int        `json:"page_size" form:"page_size"` // 每页大小
	UserType  string     `json:"user_type" form:"user_type"` // 用户类型
	Username  string     `json:"username" form:"username"`   // 用户名
	Module    string     `json:"module" form:"module"`       // 模块
	Action    string     `json:"action" form:"action"`       // 操作
	IP        string     `json:"ip" form:"ip"`               // IP地址
	StartTime *time.Time `json:"start_time" form:"start_time"`
	EndTime   *time.Time `json:"end_time" form:"end_time"`
	SortBy    string     `json:"sort_by" form:"sort_by"`     // 排序字段
	SortDesc  bool       `json:"sort_desc" form:"sort_desc"` // 是否降序
}

// GetSortDirection 获取排序方向
func (q *OperationLogQuery) GetSortDirection() string {
	if q.SortDesc {
		return "DESC"
	}
	return "ASC"
}

// Normalize 标准化分页请求
func (q *OperationLogQuery) Normalize() {
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

// OperationLogRepo 操作日志仓库
type OperationLogRepo struct {
	DB *gorm.DB
}

// NewOperationLogRepo 创建操作日志仓库
func NewOperationLogRepo(db *gorm.DB) OperationLogRepository {
	return &OperationLogRepo{DB: db}
}

// Create 创建操作日志
func (r *OperationLogRepo) Create(log *model.OperationLog) error {
	return r.DB.Create(log).Error
}

// GetLogs 获取操作日志列表
func (r *OperationLogRepo) GetLogs(query *OperationLogQuery) (*PageResult, error) {
	// 标准化分页请求
	query.Normalize()

	// 构建查询选项
	var opts []options.Option

	// 添加分页选项
	opts = append(opts, options.WithPage(query.Page, query.PageSize))

	// 添加排序选项
	opts = append(opts, options.WithOrder(query.SortBy, query.GetSortDirection()))

	// 添加条件过滤选项
	if query.UserType != "" {
		opts = append(opts, options.WithExactMatch("user_type", query.UserType))
	}

	if query.Username != "" {
		opts = append(opts, options.WithLike("username", query.Username))
	}

	if query.Module != "" {
		opts = append(opts, options.WithExactMatch("module", query.Module))
	}

	if query.Action != "" {
		opts = append(opts, options.WithExactMatch("action", query.Action))
	}

	if query.IP != "" {
		opts = append(opts, options.WithLike("ip", query.IP))
	}

	// 添加时间范围选项
	if query.StartTime != nil || query.EndTime != nil {
		opts = append(opts, options.WithTimeRange("operate_at", query.StartTime, query.EndTime))
	}

	// 构建查询
	tx := r.DB.Model(&model.OperationLog{})

	// 获取总数
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}

	// 应用所有选项
	tx = options.Apply(tx, opts...)

	// 执行查询
	var logs []model.OperationLog
	if err := tx.Find(&logs).Error; err != nil {
		return nil, err
	}

	// 构建响应
	return &PageResult{
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
		List:     logs,
	}, nil
}

// Delete 删除操作日志
func (r *OperationLogRepo) Delete(id uint) error {
	return r.DB.Delete(&model.OperationLog{}, id).Error
}

// BatchDelete 批量删除操作日志
func (r *OperationLogRepo) BatchDelete(ids []uint) error {
	return r.DB.Where("id IN ?", ids).Delete(&model.OperationLog{}).Error
}
