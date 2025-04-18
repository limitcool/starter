package repository

import (
	"errors"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// SysUserRepo 系统用户仓库
type SysUserRepo struct {
	DB *gorm.DB
}

// NewSysUserRepo 创建系统用户仓库
func NewSysUserRepo(db *gorm.DB) SysUserRepository {
	return &SysUserRepo{DB: db}
}

// GetByID 根据ID获取系统用户
func (r *SysUserRepo) GetByID(id int64) (*model.SysUser, error) {
	var user model.SysUser
	err := r.DB.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	return &user, err
}

// GetByUsername 根据用户名获取系统用户
func (r *SysUserRepo) GetByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	err := r.DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	return &user, err
}

// Create 创建系统用户
func (r *SysUserRepo) Create(user *model.SysUser) error {
	return r.DB.Create(user).Error
}

// Update 更新系统用户
func (r *SysUserRepo) Update(user *model.SysUser) error {
	return r.DB.Save(user).Error
}

// UpdateFields 更新系统用户字段
func (r *SysUserRepo) UpdateFields(id int64, fields map[string]interface{}) error {
	return r.DB.Model(&model.SysUser{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除系统用户
func (r *SysUserRepo) Delete(id int64) error {
	return r.DB.Delete(&model.SysUser{}, id).Error
}

// List 获取系统用户列表
func (r *SysUserRepo) List(query *model.SysUserQuery) ([]model.SysUser, int64, error) {
	db := r.DB.Model(&model.SysUser{})

	// 应用查询条件
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+query.Nickname+"%")
	}
	if query.Email != "" {
		db = db.Where("email LIKE ?", "%"+query.Email+"%")
	}
	if query.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+query.Phone+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(int(offset)).Limit(int(query.PageSize))
	}

	// 排序
	if query.OrderBy != "" {
		direction := "ASC"
		if query.OrderDesc {
			direction = "DESC"
		}
		db = db.Order(query.OrderBy + " " + direction)
	} else {
		db = db.Order("id DESC")
	}

	// 执行查询
	var users []model.SysUser
	if err := db.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateAvatar 更新系统用户头像
func (r *SysUserRepo) UpdateAvatar(userID int64, fileID uint) error {
	// 开始事务
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 查找用户
	user := model.SysUser{}
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return errorx.ErrNotFound.WithError(err)
	}

	// 更新用户头像
	user.AvatarFileID = fileID
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
