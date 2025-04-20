package services

import (
	"context"
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
func (s *APIService) Create(ctx context.Context, api *model.API) error {
	// 检查API是否已存在
	existingAPI, err := s.apiRepo.GetByPath(ctx, api.Path, api.Method)
	if err == nil && existingAPI.ID > 0 {
		return fmt.Errorf("API已存在: %s %s", api.Method, api.Path)
	}

	return s.apiRepo.Create(ctx, api)
}

// Update 更新API
func (s *APIService) Update(ctx context.Context, api *model.API) error {
	// 检查API是否存在
	existingAPI, err := s.apiRepo.GetByID(ctx, api.ID)
	if err != nil {
		return fmt.Errorf("API不存在: %w", err)
	}

	// 如果路径或方法发生变化，检查是否与其他API冲突
	if existingAPI.Path != api.Path || existingAPI.Method != api.Method {
		conflictAPI, err := s.apiRepo.GetByPath(ctx, api.Path, api.Method)
		if err == nil && conflictAPI.ID > 0 && conflictAPI.ID != api.ID {
			return fmt.Errorf("API已存在: %s %s", api.Method, api.Path)
		}
	}

	return s.apiRepo.Update(ctx, api)
}

// Delete 删除API
func (s *APIService) Delete(ctx context.Context, id uint) error {
	return s.apiRepo.Delete(ctx, id)
}

// GetByID 根据ID获取API
func (s *APIService) GetByID(ctx context.Context, id uint) (*model.API, error) {
	return s.apiRepo.GetByID(ctx, id)
}

// GetAll 获取所有API
func (s *APIService) GetAll(ctx context.Context) ([]*model.API, error) {
	return s.apiRepo.GetAll(ctx)
}

// GetByPath 根据路径获取API
func (s *APIService) GetByPath(ctx context.Context, path string, method string) (*model.API, error) {
	return s.apiRepo.GetByPath(ctx, path, method)
}

// GetByMenuID 获取菜单关联的API
func (s *APIService) GetByMenuID(ctx context.Context, menuID uint) ([]*model.API, error) {
	return s.apiRepo.GetByMenuID(ctx, menuID)
}
