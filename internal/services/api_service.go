package services

import (
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/repository"
)

// APIService API服务
type APIService struct {
	apiRepo *repository.APIRepo
}

// NewAPIService 创建API服务
func NewAPIService(apiRepo *repository.APIRepo) *APIService {
	return &APIService{
		apiRepo: apiRepo,
	}
}

// Create 创建API
func (s *APIService) Create(api *model.API) error {
	// 检查API是否已存在
	existingAPI, err := s.apiRepo.GetByPath(api.Path, api.Method)
	if err == nil && existingAPI.ID > 0 {
		return fmt.Errorf("API已存在: %s %s", api.Method, api.Path)
	}

	return s.apiRepo.Create(api)
}

// Update 更新API
func (s *APIService) Update(api *model.API) error {
	// 检查API是否存在
	existingAPI, err := s.apiRepo.GetByID(api.ID)
	if err != nil {
		return fmt.Errorf("API不存在: %w", err)
	}

	// 如果路径或方法发生变化，检查是否与其他API冲突
	if existingAPI.Path != api.Path || existingAPI.Method != api.Method {
		conflictAPI, err := s.apiRepo.GetByPath(api.Path, api.Method)
		if err == nil && conflictAPI.ID > 0 && conflictAPI.ID != api.ID {
			return fmt.Errorf("API已存在: %s %s", api.Method, api.Path)
		}
	}

	return s.apiRepo.Update(api)
}

// Delete 删除API
func (s *APIService) Delete(id uint) error {
	return s.apiRepo.Delete(id)
}

// GetByID 根据ID获取API
func (s *APIService) GetByID(id uint) (*model.API, error) {
	return s.apiRepo.GetByID(id)
}

// GetAll 获取所有API
func (s *APIService) GetAll() ([]*model.API, error) {
	return s.apiRepo.GetAll()
}

// GetByPath 根据路径获取API
func (s *APIService) GetByPath(path string, method string) (*model.API, error) {
	return s.apiRepo.GetByPath(path, method)
}

// GetByMenuID 获取菜单关联的API
func (s *APIService) GetByMenuID(menuID uint) ([]*model.API, error) {
	return s.apiRepo.GetByMenuID(menuID)
}
