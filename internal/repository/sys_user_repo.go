package repository

import (
	"errors"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// SysUserRepo 系统用户仓库
type SysUserRepo struct {
	DB *gorm.DB
}

// NewSysUserRepo 创建系统用户仓库
func NewSysUserRepo(db *gorm.DB) *SysUserRepo {
	return &SysUserRepo{DB: db}
}

// GetByID 根据ID获取系统用户
func (r *SysUserRepo) GetByID(id int64) (*model.SysUser, error) {
	var user model.SysUser
	err := r.DB.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		notFoundErr := errorx.Errorf(errorx.ErrUserNotFound, "系统用户ID %d 不存在", id)
		return nil, errorx.WrapError(notFoundErr, "")
	}
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询系统用户失败: id=%d", id))
	}

	// 获取用户的角色
	if err := r.DB.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询系统用户角色失败: id=%d", id))
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return &user, nil
}

// GetByUsername 根据用户名获取系统用户
func (r *SysUserRepo) GetByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	err := r.DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		notFoundErr := errorx.Errorf(errorx.ErrUserNotFound, "系统用户名 %s 不存在", username)
		return nil, errorx.WrapError(notFoundErr, "")
	}
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询系统用户失败: username=%s", username))
	}

	// 获取用户的角色
	if err := r.DB.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询系统用户角色失败: username=%s", username))
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return &user, nil
}

// Create 创建系统用户
func (r *SysUserRepo) Create(user *model.SysUser) error {
	err := r.DB.Create(user).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("创建系统用户失败: username=%s", user.Username))
	}
	return nil
}

// Update 更新系统用户
func (r *SysUserRepo) Update(user *model.SysUser) error {
	err := r.DB.Save(user).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新系统用户失败: id=%d, username=%s", user.ID, user.Username))
	}
	return nil
}

// UpdateFields 更新系统用户字段
func (r *SysUserRepo) UpdateFields(id int64, fields map[string]any) error {
	err := r.DB.Model(&model.SysUser{}).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新系统用户字段失败: id=%d", id))
	}
	return nil
}

// Delete 删除系统用户
func (r *SysUserRepo) Delete(id int64) error {
	err := r.DB.Delete(&model.SysUser{}, id).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除系统用户失败: id=%d", id))
	}
	return nil
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
		return nil, 0, errorx.WrapError(err, "获取系统用户总数失败")
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
		return nil, 0, errorx.WrapError(err, "查询系统用户列表失败")
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFoundErr := errorx.Errorf(errorx.ErrUserNotFound, "系统用户ID %d 不存在", userID)
			return errorx.WrapError(notFoundErr, "")
		}
		return errorx.WrapError(err, fmt.Sprintf("查询系统用户失败: id=%d", userID))
	}

	// 更新用户头像
	user.AvatarFileID = fileID
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return errorx.WrapError(err, fmt.Sprintf("更新系统用户头像失败: id=%d, fileID=%d", userID, fileID))
	}

	if err := tx.Commit().Error; err != nil {
		return errorx.WrapError(err, fmt.Sprintf("提交事务失败: 更新系统用户头像, userID=%d, fileID=%d", userID, fileID))
	}
	return nil
}
