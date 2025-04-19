package model

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// 初始化雪花ID节点
func init() {
	// 修改雪花ID的纪元为当前时间，这样可以减少时间戳部分的位数
	// 将纪元设置为2025年1月1日，这样可以减少时间戳部分的位数
	snowflake.Epoch = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
}

// GenerateSnowflakeID 生成雪花ID
func GenerateSnowflakeID() int64 {
	// 创建雪花ID节点
	node, err := snowflake.NewNode(1) // 使用节点ID 1
	if err != nil {
		logger.Error("Failed to create snowflake node", "error", err)
		return time.Now().UnixNano() // 备用方案：使用时间戳
	}

	// 生成雪花ID
	id := node.Generate().Int64()

	// 确保生成的ID小于2^53-1 (JavaScript安全整数最大值)
	// 如果ID大于2^53-1，则取模使其小于2^53-1
	const maxSafeInt = 9007199254740991 // 2^53-1
	if id > maxSafeInt {
		id = id % maxSafeInt
		logger.Info("Snowflake ID exceeds JavaScript safe integer range, taking modulo", "original_id", id, "new_id", id)
	}

	return id
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// 全局数据库访问函数已移除
// 请使用依赖注入或组件实例获取数据库连接
// 例如：sqldb.Instance().DB() 或 mongodb.Instance().GetDB()

// 全局错误定义
var ErrNotFound = gorm.ErrRecordNotFound

// BaseModel 基础模型结构 - 普通自增ID
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SnowflakeModel 使用雪花ID的基础模型结构 - 适用于需要分布式ID的模型（如User和SysUser）
type SnowflakeModel struct {
	ID        int64          `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UUIDModel 使用UUID的基础模型结构 - 适用于需要全局唯一ID且无序的场景
type UUIDModel struct {
	ID        string         `gorm:"primarykey;type:varchar(36)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 在创建记录前自动生成雪花ID
func (m *SnowflakeModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == 0 {
		m.ID = GenerateSnowflakeID()
		logger.Info("Generated snowflake ID", "id", m.ID)
	}
	return nil
}

// BeforeCreate 在创建记录前自动生成UUID
func (m *UUIDModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = GenerateUUID()
	}
	return nil
}
