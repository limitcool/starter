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
	Mobile    string    `json:"mobile" gorm:"size:20;comment:手机号"`
	Enabled   bool      `json:"enabled" gorm:"default:true;comment:是否启用"`
	Remark    string    `json:"remark" gorm:"size:500;comment:备注"`
	LastLogin time.Time `json:"last_login" gorm:"comment:最后登录时间"`

	// 关联
	Roles     []*Role  `json:"roles" gorm:"many2many:sys_user_role;"` // 关联的角色
	RoleIDs   []uint   `json:"role_ids" gorm:"-"`                     // 角色ID列表，不映射到数据库
	RoleCodes []string `json:"role_codes" gorm:"-"`                   // 角色编码列表
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
