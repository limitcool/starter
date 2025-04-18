package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/options"
	"github.com/limitcool/starter/internal/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 使用辅助方法获取集合
func getOperationLogCollection() *mongo.Collection {
	if mongodb.Instance() != nil {
		return mongodb.Instance().GetCollection("operation_log")
	}
	return nil
}

// OperationLog 操作记录模型
type OperationLog struct {
	BaseModel

	// 操作信息
	Module      string    `json:"module" gorm:"size:50;not null;comment:操作模块"`
	Action      string    `json:"action" gorm:"size:50;not null;comment:操作类型"`
	Description string    `json:"description" gorm:"size:255;comment:操作描述"`
	IP          string    `json:"ip" gorm:"size:50;comment:IP地址"`
	UserAgent   string    `json:"user_agent" gorm:"size:500;comment:用户代理"`
	RequestURL  string    `json:"request_url" gorm:"size:500;comment:请求URL"`
	Method      string    `json:"method" gorm:"size:10;comment:请求方法"`
	Params      string    `json:"params" gorm:"type:text;comment:请求参数"`
	Response    string    `json:"response" gorm:"type:text;comment:返回结果"`
	Status      int       `json:"status" gorm:"comment:状态码"`
	ExecuteTime int64     `json:"execute_time" gorm:"comment:执行时间(ms)"`
	OperateAt   time.Time `json:"operate_at" gorm:"comment:操作时间"`

	// 操作人信息 - 支持系统用户和普通用户
	UserType string   `json:"user_type" gorm:"size:20;comment:用户类型(sys_user/user)"`
	UserID   int64    `json:"user_id" gorm:"type:bigint;comment:用户ID"`
	Username string   `json:"username" gorm:"size:50;comment:用户名"`
	SysUser  *SysUser `json:"sys_user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	User     *User    `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (OperationLog) TableName() string {
	return "sys_operation_log"
}

// 创建操作日志
func (o *OperationLog) Create() error {
	return DB().Create(o).Error
}

// 创建系统用户操作日志
func (o *OperationLog) CreateSysUserLog(c *gin.Context, userID int64, username string, module, action, description string, startTime time.Time) error {
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
	operationLog := OperationLog{
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

	return operationLog.Create()
}

// 创建普通用户操作日志
func (o *OperationLog) CreateUserLog(c *gin.Context, userID int64, username string, module, action, description string, startTime time.Time) error {
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
	operationLog := OperationLog{
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

	return operationLog.Create()
}

// 分页获取操作日志
func (o *OperationLog) GetPageList(query any, opts ...options.Option) ([]OperationLog, int64, error) {
	// 构建查询
	tx := DB().Model(&OperationLog{})

	// 获取总数
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用所有选项
	tx = options.Apply(tx, opts...)

	// 执行查询
	var logs []OperationLog
	if err := tx.Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// 删除操作日志
func (o *OperationLog) Delete(id uint) error {
	return DB().Delete(&OperationLog{}, id).Error
}

// 批量删除操作日志
func (o *OperationLog) BatchDelete(ids []uint) error {
	return DB().Where("id IN ?", ids).Delete(&OperationLog{}).Error
}

func (OperationLog) Registry() {
	var ctx = context.Background()
	coll := getOperationLogCollection()
	if coll == nil {
		return
	}
	// 创建索引等操作
	coll.FindOne(ctx, bson.M{"module": "system"})
}
