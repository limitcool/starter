package repository

import (
	"context"

	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

// MenuButtonRepo 菜单按钮仓库
type MenuButtonRepo struct {
	DB *gorm.DB
}

// NewMenuButtonRepo 创建菜单按钮仓库
func NewMenuButtonRepo(db *gorm.DB) *MenuButtonRepo {
	return &MenuButtonRepo{
		DB: db,
	}
}

// Create 创建菜单按钮
func (r *MenuButtonRepo) Create(ctx context.Context, button *model.MenuButton) error {
	return r.DB.WithContext(ctx).Create(button).Error
}

// Update 更新菜单按钮
func (r *MenuButtonRepo) Update(ctx context.Context, button *model.MenuButton) error {
	return r.DB.WithContext(ctx).Save(button).Error
}

// Delete 删除菜单按钮
func (r *MenuButtonRepo) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&model.MenuButton{}, id).Error
}

// GetByID 根据ID获取菜单按钮
func (r *MenuButtonRepo) GetByID(ctx context.Context, id uint) (*model.MenuButton, error) {
	var button model.MenuButton
	err := r.DB.WithContext(ctx).First(&button, id).Error
	return &button, err
}

// GetAll 获取所有菜单按钮
func (r *MenuButtonRepo) GetAll(ctx context.Context) ([]*model.MenuButton, error) {
	var buttons []*model.MenuButton
	err := r.DB.WithContext(ctx).Find(&buttons).Error
	return buttons, err
}

// GetByMenuID 获取菜单的所有按钮
func (r *MenuButtonRepo) GetByMenuID(ctx context.Context, menuID uint) ([]*model.MenuButton, error) {
	var buttons []*model.MenuButton
	err := r.DB.WithContext(ctx).Where("menu_id = ?", menuID).Find(&buttons).Error
	return buttons, err
}

// GetByPermission 根据权限标识获取按钮
func (r *MenuButtonRepo) GetByPermission(ctx context.Context, permission string) (*model.MenuButton, error) {
	var button model.MenuButton
	err := r.DB.WithContext(ctx).Where("permission = ?", permission).First(&button).Error
	return &button, err
}

// WithTx 使用事务
func (r *MenuButtonRepo) WithTx(tx *gorm.DB) *MenuButtonRepo {
	return &MenuButtonRepo{DB: tx}
}
