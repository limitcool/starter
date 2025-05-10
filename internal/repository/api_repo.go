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
	genericRepo := NewGenericRepo[model.API](db).SetErrorCode(errorx.ErrNotFoundCodeValue)

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
	// 使用泛型仓库的List方法
	apis, err := r.GenericRepo.List(ctx, 1, 1000, nil)
	if err != nil {
		return nil, errorx.WrapError(err, "查询所有API失败")
	}

	// 转换为指针切片
	apiPtrs := make([]*model.API, len(apis))
	for i := range apis {
		apiPtrs[i] = &apis[i]
	}

	return apiPtrs, nil
}

// GetByPath 根据路径获取API
func (r *APIRepo) GetByPath(ctx context.Context, path string, method string) (*model.API, error) {
	// 使用泛型仓库的Get方法
	opts := &QueryOptions{
		Condition: "path = ? AND method = ?",
		Args:      []any{path, method},
	}

	api, err := r.GenericRepo.Get(ctx, nil, opts)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回nil表示未找到
		}
		return nil, errorx.WrapError(err, fmt.Sprintf("查询API失败: path=%s, method=%s", path, method))
	}

	return api, nil
}

// GetByMenuID 获取菜单关联的API
func (r *APIRepo) GetByMenuID(ctx context.Context, menuID uint) ([]*model.API, error) {
	// 创建MenuAPI的泛型仓库
	menuAPIRepo := NewGenericRepo[model.MenuAPI](r.DB)

	// 使用查询选项
	menuAPIopts := &QueryOptions{
		Condition: "menu_id = ?",
		Args:      []any{menuID},
	}

	// 获取所有关联记录
	menuAPIs, err := menuAPIRepo.List(ctx, 1, 1000, menuAPIopts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单API关联失败: menuID=%d", menuID))
	}

	// 提取API ID
	var apiIDs []uint
	for _, menuAPI := range menuAPIs {
		apiIDs = append(apiIDs, menuAPI.APIID)
	}

	if len(apiIDs) == 0 {
		return []*model.API{}, nil
	}

	// 使用泛型仓库的List方法
	opts := &QueryOptions{
		Condition: "id IN ?",
		Args:      []any{apiIDs},
	}

	apis, err := r.GenericRepo.List(ctx, 1, 1000, opts)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询菜单关联API失败: menuID=%d", menuID))
	}

	// 转换为指针切片
	apiPtrs := make([]*model.API, len(apis))
	for i := range apis {
		apiPtrs[i] = &apis[i]
	}

	return apiPtrs, nil
}

// WithTx 使用事务
func (r *APIRepo) WithTx(tx *gorm.DB) *APIRepo {
	// 创建通用仓库并设置错误码
	genericRepo := NewGenericRepo[model.API](tx).SetErrorCode(errorx.ErrNotFoundCodeValue)

	return &APIRepo{
		DB:          tx,
		GenericRepo: genericRepo,
	}
}
