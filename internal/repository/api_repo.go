package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// APIRepo API仓库
type APIRepo struct {
	DB          *gorm.DB
	GenericRepo Repository[model.API] // 使用接口而非具体实现
}

// NewAPIRepo 创建API仓库
func NewAPIRepo(db *gorm.DB) *APIRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.API](db).SetErrorCode(errorx.ErrorNotFoundCode)

	return &APIRepo{
		DB:          db,
		GenericRepo: genericRepo,
	}
}

// Create 创建API
func (r *APIRepo) Create(ctx context.Context, api *model.API) error {
	// 使用仓库接口
	return r.GenericRepo.Create(ctx, api)
}

// Update 更新API
func (r *APIRepo) Update(ctx context.Context, api *model.API) error {
	// 使用仓库接口
	return r.GenericRepo.Update(ctx, api)
}

// Delete 删除API
func (r *APIRepo) Delete(ctx context.Context, id uint) error {
	// 使用仓库接口
	return r.GenericRepo.Delete(ctx, id)
}

// GetByID 根据ID获取API
func (r *APIRepo) GetByID(ctx context.Context, id uint) (*model.API, error) {
	// 使用仓库接口
	return r.GenericRepo.GetByID(ctx, id)
}

// GetAll 获取所有API
func (r *APIRepo) GetAll(ctx context.Context) ([]*model.API, error) {
	var apis []*model.API
	err := r.DB.WithContext(ctx).Find(&apis).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询所有API失败")
	}
	return apis, nil
}

// GetByPath 根据路径获取API
func (r *APIRepo) GetByPath(ctx context.Context, path string, method string) (*model.API, error) {
	// 使用数据库直接查询
	var api model.API
	err := r.DB.WithContext(ctx).Where("path = ? AND method = ?", path, method).First(&api).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 返回nil表示未找到
	}
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询API失败: path=%s, method=%s", path, method))
	}
	return &api, nil
}

// GetByMenuID 获取菜单关联的API
func (r *APIRepo) GetByMenuID(ctx context.Context, menuID uint) ([]*model.API, error) {
	var apis []*model.API
	err := r.DB.WithContext(ctx).Joins("JOIN sys_menu_api ON sys_menu_api.api_id = sys_api.id").
		Where("sys_menu_api.menu_id = ?", menuID).
		Find(&apis).Error
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单关联API失败: menuID=%d", menuID))
	}
	return apis, nil
}

// WithTx 使用事务
func (r *APIRepo) WithTx(tx *gorm.DB) *APIRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.API](tx).SetErrorCode(errorx.ErrorNotFoundCode)

	return &APIRepo{
		DB:          tx,
		GenericRepo: genericRepo,
	}
}
