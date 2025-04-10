package services

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/options"
	"github.com/limitcool/starter/internal/storage/sqldb"
)

// OperationLogService 操作日志服务
type OperationLogService struct {
}

// NewOperationLogService 创建操作日志服务
func NewOperationLogService() *OperationLogService {
	return &OperationLogService{}
}

// CreateSysUserLog 创建系统用户操作日志
func (s *OperationLogService) CreateSysUserLog(c *gin.Context, userID int64, username string, module, action, description string, startTime time.Time) error {
	// 计算执行时间
	executeTime := time.Since(startTime).Milliseconds()

	// 获取请求相关信息
	method := c.Request.Method
	requestURL := c.Request.URL.String()
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 获取请求参数
	var params string
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		// 获取请求体
		bodyData, exists := c.Get("requestBody")
		if exists {
			paramsBytes, _ := json.Marshal(bodyData)
			params = string(paramsBytes)
		}
	} else {
		// 获取查询参数
		queryParams := c.Request.URL.Query()
		paramsBytes, _ := json.Marshal(queryParams)
		params = string(paramsBytes)
	}

	// 创建操作日志
	operationLog := model.OperationLog{
		Module:      module,
		Action:      action,
		Description: description,
		IP:          ip,
		UserAgent:   userAgent,
		RequestURL:  requestURL,
		Method:      method,
		Params:      params,
		Status:      c.Writer.Status(),
		ExecuteTime: executeTime,
		OperateAt:   startTime,
		UserType:    "sys_user",
		UserID:      userID,
		Username:    username,
	}

	return sqldb.Instance().DB().Create(&operationLog).Error
}

// CreateUserLog 创建普通用户操作日志
func (s *OperationLogService) CreateUserLog(c *gin.Context, userID int64, username string, module, action, description string, startTime time.Time) error {
	// 计算执行时间
	executeTime := time.Since(startTime).Milliseconds()

	// 获取请求相关信息
	method := c.Request.Method
	requestURL := c.Request.URL.String()
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 获取请求参数
	var params string
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		// 获取请求体
		bodyData, exists := c.Get("requestBody")
		if exists {
			paramsBytes, _ := json.Marshal(bodyData)
			params = string(paramsBytes)
		}
	} else {
		// 获取查询参数
		queryParams := c.Request.URL.Query()
		paramsBytes, _ := json.Marshal(queryParams)
		params = string(paramsBytes)
	}

	// 创建操作日志
	operationLog := model.OperationLog{
		Module:      module,
		Action:      action,
		Description: description,
		IP:          ip,
		UserAgent:   userAgent,
		RequestURL:  requestURL,
		Method:      method,
		Params:      params,
		Status:      c.Writer.Status(),
		ExecuteTime: executeTime,
		OperateAt:   startTime,
		UserType:    "user",
		UserID:      userID,
		Username:    username,
	}

	return sqldb.Instance().DB().Create(&operationLog).Error
}

// GetOperationLogs 分页获取操作日志
func (s *OperationLogService) GetOperationLogs(query *v1.OperationLogQuery) (*response.PageResult[[]model.OperationLog], error) {
	// 标准化分页请求
	query.PageRequest.Normalize()

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
	tx := sqldb.Instance().DB().Model(&model.OperationLog{})

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
	return response.NewPageResult(logs, total, query.Page, query.PageSize), nil
}

// DeleteOperationLog 删除操作日志
func (s *OperationLogService) DeleteOperationLog(id uint) error {
	return sqldb.Instance().DB().Delete(&model.OperationLog{}, id).Error
}

// BatchDeleteOperationLogs 批量删除操作日志
func (s *OperationLogService) BatchDeleteOperationLogs(ids []uint) error {
	return sqldb.Instance().DB().Where("id IN ?", ids).Delete(&model.OperationLog{}).Error
}
