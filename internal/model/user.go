package model

import (
	"context"
	"time"

	"github.com/limitcool/starter/internal/errorx"
	"gorm.io/gorm"
)

// User 用户模型
// 可以作为普通用户，也可以在合并模式下同时作为管理员用户
type User struct {
	SnowflakeModel

	Username     string     `json:"username" gorm:"size:50;not null;unique;comment:用户名"`
	Password     string     `json:"-" gorm:"size:100;not null;comment:密码"`
	Nickname     string     `json:"nickname" gorm:"size:50;comment:昵称"`
	AvatarFileID int64      `json:"-" gorm:"comment:头像文件ID"`
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
	*GenericRepo[User]
}

// NewUserRepo 创建用户仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
	genericRepo := NewGenericRepo[User](db)
	genericRepo.ErrorCode = errorx.ErrUserNotFound.Code()

	return &UserRepo{
		GenericRepo: genericRepo,
	}
}

// GetByID 根据ID获取用户
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*User, error) {
	user, err := r.Get(ctx, id, nil)
	if err != nil {
		return nil, errorx.ErrQueryUser.New(ctx, errorx.None).Wrap(err)
	}
	return user, nil
}

// GetUserWithAvatar 获取用户信息，包括头像
func (r *UserRepo) GetUserWithAvatar(ctx context.Context, id int64) (*User, error) {
	// 先查询用户基本信息
	user, err := r.GenericRepo.Get(ctx, id, nil)
	if err != nil {
		return nil, errorx.ErrQueryUser.New(ctx, errorx.None)
	}

	// 如果用户有头像，再预加载头像
	if user.AvatarFileID > 0 {
		user, err = r.Get(ctx, id, &QueryOptions{
			Preloads: []string{"AvatarFile"},
		})
		if err != nil {
			return nil, errorx.ErrQueryUserAvatar.New(ctx, errorx.None).Wrap(err)
		}

		// 设置头像URL
		if user.AvatarFile != nil {
			user.AvatarURL = user.AvatarFile.URL
		}
	}

	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	user, err := r.Get(ctx, nil, &QueryOptions{
		Condition: "username = ?",
		Args:      []any{username},
	})
	if err != nil {
		if errorx.ErrNotFound.Is(err) {
			return nil, errorx.ErrNotFound.New(ctx, errorx.None)
		}
		return nil, errorx.ErrQueryUser.New(ctx, errorx.None)
	}
	return user, nil
}

// IsExist 检查用户是否存在
func (r *UserRepo) IsExist(ctx context.Context, username string) (bool, error) {
	count, err := r.Count(ctx, &QueryOptions{
		Condition: "username = ?",
		Args:      []any{username},
	})
	if err != nil {
		return false, errorx.ErrCheckUserExist.New(ctx, errorx.None)
	}
	return count > 0, nil
}

// ListUsers 获取用户列表
func (r *UserRepo) ListUsers(ctx context.Context, page, pageSize int, keyword string) ([]User, int64, error) {
	var opts *QueryOptions

	// 如果有关键字，添加模糊查询条件
	if keyword != "" {
		opts = &QueryOptions{
			Condition: "username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			Args:      []any{"%" + keyword + "%", "%" + keyword + "%", "%" + keyword + "%"},
			Preloads:  []string{"AvatarFile"},
		}
	} else {
		opts = &QueryOptions{
			Preloads: []string{"AvatarFile"},
		}
	}

	// 获取用户列表
	users, err := r.List(ctx, page, pageSize, opts)
	if err != nil {
		return nil, 0, errorx.ErrQueryUserList.New(ctx, errorx.None).Wrap(err)
	}

	// 设置头像URL
	for i := range users {
		if users[i].AvatarFile != nil {
			users[i].AvatarURL = users[i].AvatarFile.URL
		}
	}

	// 获取总数
	total, err := r.Count(ctx, opts)
	if err != nil {
		return nil, 0, errorx.ErrQueryUserTotal.New(ctx, errorx.None).Wrap(err)
	}

	return users, total, nil
}

// UpdateAvatar 更新用户头像
func (r *UserRepo) UpdateAvatar(ctx context.Context, userID int64, fileID int64) error {
	return r.DB.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("avatar_file_id", fileID).Error
}

// UpdatePassword 更新用户密码
func (r *UserRepo) UpdatePassword(ctx context.Context, userID int64, password string) error {
	return r.DB.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("password", password).Error
}

// UpdateLastLogin 更新最后登录信息
func (r *UserRepo) UpdateLastLogin(ctx context.Context, userID int64, ip string) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(map[string]any{
		"last_login": now,
		"last_ip":    ip,
	}).Error
}
