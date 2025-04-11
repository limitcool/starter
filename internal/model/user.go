package model

import (
	"errors"
	"time"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"gorm.io/gorm"
)

// User 普通用户
type User struct {
	SnowflakeModel

	Username     string     `json:"username" gorm:"size:50;not null;unique;comment:用户名"`
	Password     string     `json:"-" gorm:"size:100;not null;comment:密码"`
	Nickname     string     `json:"nickname" gorm:"size:50;comment:昵称"`
	AvatarFileID uint       `json:"-" gorm:"size:255;comment:头像文件ID"`
	AvatarURL    string     `json:"avatar" gorm:"-"`                            // 头像URL，不存储到数据库
	AvatarFile   *File      `json:"avatar_file" gorm:"foreignKey:AvatarFileID"` // 关联的头像文件
	Email        string     `json:"email" gorm:"size:100;comment:邮箱"`
	Mobile       string     `json:"mobile" gorm:"size:20;comment:手机号"`
	Enabled      bool       `json:"enabled" gorm:"default:true;comment:是否启用"`
	Remark       string     `json:"remark" gorm:"size:500;comment:备注"`
	LastLogin    *time.Time `json:"last_login" gorm:"comment:最后登录时间"`
	LastIP       string     `json:"last_ip" gorm:"size:50;comment:最后登录IP"`

	// 普通用户特有字段
	Gender     string     `json:"gender" gorm:"size:10;comment:性别"`
	Birthday   *time.Time `json:"birthday" gorm:"comment:生日"`
	Address    string     `json:"address" gorm:"size:255;comment:地址"`
	RegisterIP string     `json:"register_ip" gorm:"size:50;comment:注册IP"`
}

func (User) TableName() string {
	return "user"
}

func NewUser() *User {
	return &User{}
}

// IsExist 判断用户是否存在
func (u *User) IsExist() (bool, error) {
	db := sqldb.Instance().DB()
	return db.Model(&User{}).Where("username = ?", u.Username).First(&User{}).RowsAffected > 0, nil
}

// Create 创建用户
func (u *User) Create() error {
	db := sqldb.Instance().DB()
	return db.Create(u).Error
}

// GetUserByUsername 根据用户名获取用户
func (u *User) GetUserByUsername(username string) (*User, error) {
	db := sqldb.Instance().DB()
	var user User
	err := db.Model(&User{}).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound.WithError(err)
	}
	return &user, err
}

func (u *User) GetUserByID(id int64) (*User, error) {
	db := sqldb.Instance().DB()
	var user User
	err := db.Model(&User{}).Where("id = ?", id).First(&user).Error
	return &user, err
}
