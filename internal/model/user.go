package model

import (
	"time"
)

// User 用户模型
// 可以作为普通用户，也可以在合并模式下同时作为管理员用户
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

	// 合并模式下的管理员字段
	IsAdmin bool `json:"is_admin" gorm:"default:false;comment:是否管理员"`
	// 关联
	Roles     []*Role  `json:"roles" gorm:"many2many:user_role;"` // 关联的角色
	RoleIDs   []int64  `json:"role_ids" gorm:"-"`                 // 角色ID列表，不映射到数据库
	RoleCodes []string `json:"role_codes" gorm:"-"`               // 角色编码列表
}

func (User) TableName() string {
	return "user"
}

func NewUser() *User {
	return &User{}
}

// 以下方法已移动到 repository/user_repository.go
// IsExist
// Create
// GetUserByUsername
// GetUserByID
