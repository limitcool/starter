package repository

import (
	"github.com/limitcool/starter/internal/model"
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
	return r.DB.Create(api).Error
}

// Update 更新API
func (r *APIRepo) Update(api *model.API) error {
	return r.DB.Save(api).Error
}

// Delete 删除API
func (r *APIRepo) Delete(id uint) error {
	return r.DB.Delete(&model.API{}, id).Error
}

// GetByID 根据ID获取API
func (r *APIRepo) GetByID(id uint) (*model.API, error) {
	var api model.API
	err := r.DB.First(&api, id).Error
	return &api, err
}

// GetAll 获取所有API
func (r *APIRepo) GetAll() ([]*model.API, error) {
	var apis []*model.API
	err := r.DB.Find(&apis).Error
	return apis, err
}

// GetByPath 根据路径获取API
func (r *APIRepo) GetByPath(path string, method string) (*model.API, error) {
	var api model.API
	err := r.DB.Where("path = ? AND method = ?", path, method).First(&api).Error
	return &api, err
}

// GetByMenuID 获取菜单关联的API
func (r *APIRepo) GetByMenuID(menuID uint) ([]*model.API, error) {
	var apis []*model.API
	err := r.DB.Joins("JOIN sys_menu_api ON sys_menu_api.api_id = sys_api.id").
		Where("sys_menu_api.menu_id = ?", menuID).
		Find(&apis).Error
	return apis, err
}

// WithTx 使用事务
func (r *APIRepo) WithTx(tx *gorm.DB) *APIRepo {
	return &APIRepo{DB: tx}
}
