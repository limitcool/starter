package model

import (
	"context"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/options"
	"gorm.io/gorm"
)

// Entity 实体接口
// 所有可以被仓库管理的实体都应该实现这个接口
type Entity any

// QueryOptions 查询选项
type QueryOptions struct {
	// 查询条件
	Condition string
	// 查询参数
	Args []any
	// 查询选项
	Opts []options.Option
	// 预加载关联
	Preloads []string
}

// Repository 数据库操作接口
// 提供基本的CRUD操作
type Repository[T Entity] interface {
	// Create 创建实体
	Create(ctx context.Context, entity *T) error

	// Get 根据ID或条件获取单个实体
	// id: 实体ID，如果为nil，则使用condition和args
	// opts: 查询选项，可以为nil
	Get(ctx context.Context, id any, opts *QueryOptions) (*T, error)

	// Update 更新实体
	Update(ctx context.Context, entity *T) error

	// Delete 删除实体
	Delete(ctx context.Context, id any) error

	// List 获取实体列表
	// page, pageSize: 分页参数
	// opts: 查询选项，可以为nil
	List(ctx context.Context, page, pageSize int, opts *QueryOptions) ([]T, error)

	// Count 获取实体总数
	// opts: 查询选项，可以为nil
	Count(ctx context.Context, opts *QueryOptions) (int64, error)

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

// Get 根据ID或条件获取单个实体
func (r *GenericRepo[T]) Get(ctx context.Context, id any, opts *QueryOptions) (*T, error) {
	var entity T
	
	// 创建查询
	query := r.DB.WithContext(ctx)
	
	// 应用预加载
	if opts != nil && opts.Preloads != nil {
		for _, preload := range opts.Preloads {
			query = query.Preload(preload)
		}
	}
	
	// 应用查询选项
	if opts != nil && opts.Opts != nil && len(opts.Opts) > 0 {
		query = options.Apply(query, opts.Opts...)
	}
	
	// 执行查询
	var err error
	if id != nil {
		// 根据ID查询
		err = query.First(&entity, id).Error
	} else if opts != nil && opts.Condition != "" {
		// 根据条件查询
		err = query.Where(opts.Condition, opts.Args...).First(&entity).Error
	} else {
		// 没有ID和条件，返回错误
		return nil, errorx.ErrInvalidParams.WithMsg("查询参数不能为空")
	}
	
	if err != nil {
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
func (r *GenericRepo[T]) List(ctx context.Context, page, pageSize int, opts *QueryOptions) ([]T, error) {
	var entities []T
	
	// 创建查询
	query := r.DB.WithContext(ctx)
	
	// 应用分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	
	// 应用预加载
	if opts != nil && opts.Preloads != nil {
		for _, preload := range opts.Preloads {
			query = query.Preload(preload)
		}
	}
	
	// 应用查询选项
	if opts != nil && opts.Opts != nil && len(opts.Opts) > 0 {
		query = options.Apply(query, opts.Opts...)
	}
	
	// 应用条件
	if opts != nil && opts.Condition != "" {
		query = query.Where(opts.Condition, opts.Args...)
	}
	
	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	
	return entities, nil
}

// Count 获取实体总数
func (r *GenericRepo[T]) Count(ctx context.Context, opts *QueryOptions) (int64, error) {
	var count int64
	var entity T
	
	// 创建查询
	query := r.DB.WithContext(ctx).Model(&entity)
	
	// 应用查询选项
	if opts != nil && opts.Opts != nil && len(opts.Opts) > 0 {
		query = options.Apply(query, opts.Opts...)
	}
	
	// 应用条件
	if opts != nil && opts.Condition != "" {
		query = query.Where(opts.Condition, opts.Args...)
	}
	
	// 执行查询
	if err := query.Count(&count).Error; err != nil {
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
