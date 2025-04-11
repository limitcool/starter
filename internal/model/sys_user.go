package model

import (
	"errors"
	"time"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"gorm.io/gorm"
)

// SysUser 系统用户
type SysUser struct {
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

	// 关联
	Roles     []*Role  `json:"roles" gorm:"many2many:sys_user_role;"` // 关联的角色
	RoleIDs   []int64  `json:"role_ids" gorm:"-"`                     // 角色ID列表，不映射到数据库
	RoleCodes []string `json:"role_codes" gorm:"-"`                   // 角色编码列表
}

func (SysUser) TableName() string {
	return "sys_user"
}

func NewSysUser() *SysUser {
	return &SysUser{}
}

// GetUserByUsername 根据用户名获取用户
func (s *SysUser) GetUserByUsername(username string) (*SysUser, error) {
	var user SysUser
	db := sqldb.Instance().DB()
	err := db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	if err != nil {
		// 直接返回错误，错误会自动捕获堆栈
		return nil, err
	}

	// // 获取用户的角色
	// if err := sqldb.Instance().DB().Model(&user).Association("Roles").Find(&user.Roles); err != nil {
	// 	// 直接返回错误，错误会自动捕获堆栈
	// 	return nil, err
	// }

	// // 提取角色编码
	// for _, role := range user.Roles {
	// 	user.RoleCodes = append(user.RoleCodes, role.Code)
	// }

	return &user, nil
}

// GetUserByID 根据ID获取用户
func (s *SysUser) GetUserByID(id int64) (*SysUser, error) {
	var user SysUser
	db := sqldb.Instance().DB()
	err := db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	if err != nil {
		// 直接返回错误，错误会自动捕获堆栈
		return nil, err
	}

	// 获取用户的角色
	if err := db.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		// 直接返回错误，错误会自动捕获堆栈
		return nil, err
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return &user, nil
}
