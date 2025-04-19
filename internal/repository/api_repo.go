package repository

import (
	"errors"
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// APIRepo API仓库
type APIRepo struct {
	DB *gorm.DB
}

// NewAPIRepo 创建API仓库
func NewAPIRepo(db *gorm.DB) *APIRepo {
	return &APIRepo{
		DB: db,
	}
}

// Create 创建API
func (r *APIRepo) Create(api *model.API) error {
	err := r.DB.Create(api).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("创建API失败: path=%s, method=%s", api.Path, api.Method))
	}
	return nil
}

// Update 更新API
func (r *APIRepo) Update(api *model.API) error {
	err := r.DB.Save(api).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新API失败: id=%d, path=%s, method=%s", api.ID, api.Path, api.Method))
	}
	return nil
}

// Delete 删除API
func (r *APIRepo) Delete(id uint) error {
	err := r.DB.Delete(&model.API{}, id).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除API失败: id=%d", id))
	}
	return nil
}

// GetByID 根据ID获取API
func (r *APIRepo) GetByID(id uint) (*model.API, error) {
	var api model.API
	err := r.DB.First(&api, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		notFoundErr := errorx.Errorf(errorx.ErrNotFound, "API ID %d 不存在", id)
		return nil, errorx.WrapError(notFoundErr, "")
	}
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询API失败: id=%d", id))
	}
	return &api, nil
}

// GetAll 获取所有API
func (r *APIRepo) GetAll() ([]*model.API, error) {
	var apis []*model.API
	err := r.DB.Find(&apis).Error
	if err != nil {
		return nil, errorx.WrapError(err, "查询所有API失败")
	}
	return apis, nil
}

// GetByPath 根据路径获取API
func (r *APIRepo) GetByPath(path string, method string) (*model.API, error) {
	var api model.API
	err := r.DB.Where("path = ? AND method = ?", path, method).First(&api).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 返回nil表示未找到
	}
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询API失败: path=%s, method=%s", path, method))
	}
	return &api, nil
}

// GetByMenuID 获取菜单关联的API
func (r *APIRepo) GetByMenuID(menuID uint) ([]*model.API, error) {
	var apis []*model.API
	err := r.DB.Joins("JOIN sys_menu_api ON sys_menu_api.api_id = sys_api.id").
		Where("sys_menu_api.menu_id = ?", menuID).
		Find(&apis).Error
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单关联API失败: menuID=%d", menuID))
	}
	return apis, nil
}

// WithTx 使用事务
func (r *APIRepo) WithTx(tx *gorm.DB) *APIRepo {
	return &APIRepo{DB: tx}
}
