package model

import (
	"time"
)

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

// 以下方法已移动到 repository/operation_log_repo.go
// Create
// CreateSysUserLog
// CreateUserLog
// GetPageList
// Delete
// BatchDelete
