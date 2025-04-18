package services

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/repository"
)

// PageResult 分页结果
type PageResult struct {
	Total    int64                `json:"total"`     // 总记录数
	Page     int                  `json:"page"`      // 当前页码
	PageSize int                  `json:"page_size"` // 每页大小
	List     []model.OperationLog `json:"list"`      // 数据列表
}

// OperationLogService 操作日志服务
type OperationLogService struct {
	logRepo repository.OperationLogRepository
}

// NewOperationLogService 创建操作日志服务
func NewOperationLogService(logRepo repository.OperationLogRepository) *OperationLogService {
	return &OperationLogService{
		logRepo: logRepo,
	}
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

	return s.logRepo.Create(&operationLog)
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

	return s.logRepo.Create(&operationLog)
}

// GetOperationLogs 分页获取操作日志
func (s *OperationLogService) GetOperationLogs(query *v1.OperationLogQuery) (*PageResult, error) {
	// 标准化分页请求
	query.PageRequest.Normalize()

	// 创建仓库查询参数
	repoQuery := &repository.OperationLogQuery{
		Page:      query.Page,
		PageSize:  query.PageSize,
		UserType:  query.UserType,
		Username:  query.Username,
		Module:    query.Module,
		Action:    query.Action,
		IP:        query.IP,
		StartTime: query.StartTime,
		EndTime:   query.EndTime,
		SortBy:    query.SortBy,
		SortDesc:  query.SortDesc,
	}

	// 执行查询
	result, err := s.logRepo.GetLogs(repoQuery)
	if err != nil {
		return nil, err
	}

	// 构建响应
	return &PageResult{
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
		List:     result.List,
	}, nil
}

// DeleteOperationLog 删除操作日志
func (s *OperationLogService) DeleteOperationLog(id uint) error {
	return s.logRepo.Delete(id)
}

// BatchDeleteOperationLogs 批量删除操作日志
func (s *OperationLogService) BatchDeleteOperationLogs(ids []uint) error {
	return s.logRepo.BatchDelete(ids)
}
