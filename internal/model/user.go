package model

import (
	"context"
	"time"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
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

	// 管理员字段
	IsAdmin bool `json:"is_admin" gorm:"default:false;comment:是否管理员"`
}

func (User) TableName() string {
	return "user"
}

func NewUser() *User {
	return &User{}
}

// UserRepo 用户仓库
type UserRepo struct {
	DB *gorm.DB
}

// NewUserRepo 创建用户仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

// GetByID 根据ID获取用户
func (r *UserRepo) GetByID(ctx context.Context, id uint) (*User, error) {
	var user User
	if err := r.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrUserNotFound
		}
		return nil, errorx.WrapError(err, "查询用户失败")
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrUserNotFound
		}
		return nil, errorx.WrapError(err, "查询用户失败")
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepo) Create(ctx context.Context, user *User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *UserRepo) Update(ctx context.Context, user *User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&User{}, id).Error
}

// IsExist 检查用户是否存在
func (r *UserRepo) IsExist(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, errorx.WrapError(err, "检查用户是否存在失败")
	}
	return count > 0, nil
}

// UpdateAvatar 更新用户头像
func (r *UserRepo) UpdateAvatar(ctx context.Context, userID uint, fileID uint) error {
	return r.DB.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("avatar_file_id", fileID).Error
}

// UpdatePassword 更新用户密码
func (r *UserRepo) UpdatePassword(ctx context.Context, userID uint, password string) error {
	return r.DB.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("password", password).Error
}

// UpdateLastLogin 更新最后登录信息
func (r *UserRepo) UpdateLastLogin(ctx context.Context, userID uint, ip string) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(map[string]any{
		"last_login": now,
		"last_ip":    ip,
	}).Error
}
