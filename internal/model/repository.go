package model

import (
	"context"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// Entity 实体接口
// 所有可以被仓库管理的实体都应该实现这个接口
type Entity any

// Repository 数据库操作接口
// 提供基本的CRUD操作
type Repository[T Entity] interface {
	// Create 创建实体
	Create(ctx context.Context, entity *T) error

	// GetByID 根据ID获取实体
	GetByID(ctx context.Context, id any) (*T, error)

	// Update 更新实体
	Update(ctx context.Context, entity *T) error

	// Delete 删除实体
	Delete(ctx context.Context, id any) error

	// List 获取实体列表
	List(ctx context.Context, page, pageSize int) ([]T, error)

	// Count 获取实体总数
	Count(ctx context.Context) (int64, error)

	// Transaction 在事务中执行函数
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	// WithTx 使用事务
	WithTx(tx *gorm.DB) Repository[T]
}

// GenericRepo 通用仓库实现
type GenericRepo[T Entity] struct {
	DB        *gorm.DB
	TableName string
	ErrorCode int // 用于NotFound错误
}

// NewGenericRepo 创建通用仓库
func NewGenericRepo[T Entity](db *gorm.DB) *GenericRepo[T] {
	return &GenericRepo[T]{
		DB:        db,
		ErrorCode: errorx.ErrorNotFoundCode, // 默认错误码
	}
}

// Create 创建实体
func (r *GenericRepo[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

// GetByID 根据ID获取实体
func (r *GenericRepo[T]) GetByID(ctx context.Context, id any) (*T, error) {
	var entity T
	if err := r.DB.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrNotFound.WithMsg("记录不存在")
		}
		return nil, err
	}
	return &entity, nil
}

// Update 更新实体
func (r *GenericRepo[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

// Delete 删除实体
func (r *GenericRepo[T]) Delete(ctx context.Context, id any) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, id).Error
}

// List 获取实体列表
func (r *GenericRepo[T]) List(ctx context.Context, page, pageSize int) ([]T, error) {
	var entities []T
	offset := (page - 1) * pageSize
	if err := r.DB.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Count 获取实体总数
func (r *GenericRepo[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	var entity T
	if err := r.DB.WithContext(ctx).Model(&entity).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Transaction 在事务中执行函数
func (r *GenericRepo[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.DB.WithContext(ctx).Transaction(fn)
}

// WithTx 使用事务
func (r *GenericRepo[T]) WithTx(tx *gorm.DB) Repository[T] {
	return &GenericRepo[T]{
		DB:        tx,
		TableName: r.TableName,
		ErrorCode: r.ErrorCode,
	}
}
