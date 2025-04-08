package model

import (
	"context"
	"time"

	"github.com/limitcool/starter/internal/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 使用辅助方法获取集合
func getUserCollection() *mongo.Collection {
	return mongodb.Collection("user")
}

// User 系统用户
type User struct {
	BaseModel

	Username  string    `json:"username" gorm:"size:50;not null;unique;comment:用户名"`
	Password  string    `json:"-" gorm:"size:100;not null;comment:密码"`
	Nickname  string    `json:"nickname" gorm:"size:50;comment:昵称"`
	Avatar    string    `json:"avatar" gorm:"size:255;comment:头像"`
	Email     string    `json:"email" gorm:"size:100;comment:邮箱"`
	Phone     string    `json:"phone" gorm:"size:20;comment:手机号"`
	Status    int8      `json:"status" gorm:"default:1;comment:状态(0:禁用,1:正常)"`
	LastLogin time.Time `json:"last_login" gorm:"comment:最后登录时间"`

	// 非数据库字段
	RoleIDs   []uint   `json:"role_ids" gorm:"-"`   // 角色ID列表
	RoleCodes []string `json:"role_codes" gorm:"-"` // 角色编码列表
}

func (User) TableName() string {
	return "sys_user"
}

func (User) Registry() {
	var ctx = context.Background()
	coll := getUserCollection()
	if coll == nil {
		return
	}
	coll.FindOne(ctx, bson.M{"name": "cool"})
	coll.InsertOne(ctx, bson.M{"name": "cool"})
}
