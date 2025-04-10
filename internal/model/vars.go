package model

import (
	"time"

	"github.com/limitcool/starter/internal/storage/mongodb"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// DB 获取SQL数据库连接
// 推荐直接使用sqldb.Instance().DB()获取数据库连接
func DB() *gorm.DB {
	if sqldb.Instance() != nil {
		return sqldb.Instance().DB()
	}
	// 回退到全局变量
	if sqldb.DB != nil {
		return sqldb.DB
	}
	return nil
}

// GetMongoDB 获取MongoDB数据库
func GetMongoDB() *mongo.Database {
	// 使用直接访问方式保持兼容性
	if mongodb.Mongo == nil {
		return nil
	}
	dbName := mongodb.GetDatabaseName()
	if dbName == "" {
		return nil
	}
	return mongodb.Mongo.Database(dbName)
}

// 全局错误定义
var ErrNotFound = gorm.ErrRecordNotFound

// BaseModel 基础模型结构
type BaseModel struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
