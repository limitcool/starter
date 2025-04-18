package model

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/limitcool/starter/internal/storage/mongodb"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// 雪花ID节点，用于生成雪花ID
var snowflakeNode *snowflake.Node

// 初始化雪花ID节点
func init() {
	var err error
	snowflakeNode, err = snowflake.NewNode(1) // 使用节点ID 1
	if err != nil {
		panic(err)
	}
}

// GenerateSnowflakeID 生成雪花ID
func GenerateSnowflakeID() int64 {
	return snowflakeNode.Generate().Int64()
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// DB 获取SQL数据库连接
// 推荐直接使用sqldb.Instance().DB()获取数据库连接
func DB() *gorm.DB {
	// 优先使用组件实例
	if sqldb.Instance() != nil {
		return sqldb.Instance().DB()
	}
	// 记录警告日志
	log.Warn("使用了已弃用的数据库访问方式，请使用 sqldb.Instance().DB()", "caller", "model.DB()")
	return nil
}

// GetMongoDB 获取MongoDB数据库
func GetMongoDB() *mongo.Database {
	// 优先使用组件实例
	if mongodb.Instance() != nil {
		return mongodb.Instance().GetDB()
	}
	// 记录警告日志
	log.Warn("使用了已弃用的MongoDB访问方式，请使用 mongodb.Instance().GetDB()", "caller", "model.GetMongoDB()")
	return nil
}

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
		log.Info("生成雪花ID", "id", m.ID)
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
