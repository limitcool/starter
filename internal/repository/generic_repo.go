package repository

import (
	"context"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/options"
	"gorm.io/gorm"
)

// Entity 实体接口
// 所有可以被仓库管理的实体都应该实现这个接口
type Entity interface {
	TableName() string
}

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
	// 基础方法
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

	// 便捷方法
	// GetByID 根据ID获取实体
	GetByID(ctx context.Context, id any) (*T, error)

	// FindByField 根据字段查询实体
	FindByField(ctx context.Context, field string, value any) (*T, error)

	// UpdateFields 更新实体字段
	UpdateFields(ctx context.Context, id any, fields map[string]any) error

	// GetPage 分页查询
	GetPage(ctx context.Context, page, pageSize int, condition string, args ...any) ([]T, int64, error)

	// BatchDelete 批量删除实体
	BatchDelete(ctx context.Context, ids []any) error
}

// GenericRepo 通用仓库实现
type GenericRepo[T Entity] struct {
	DB        *gorm.DB
	ErrorCode int // 用于NotFound错误
}

// NewGenericRepo 创建通用仓库
func NewGenericRepo[T Entity](db *gorm.DB) *GenericRepo[T] {
	return &GenericRepo[T]{
		DB:        db,
		ErrorCode: errorx.ErrNotFoundCodeValue, // 默认错误码
	}
}

// Create 创建实体
func (r *GenericRepo[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

// applyQueryOptions 应用查询选项
func (r *GenericRepo[T]) applyQueryOptions(query *gorm.DB, opts *QueryOptions) *gorm.DB {
	if opts == nil {
		return query
	}

	// 应用预加载
	if opts.Preloads != nil {
		for _, preload := range opts.Preloads {
			query = query.Preload(preload)
		}
	}

	// 应用查询选项
	if len(opts.Opts) > 0 {
		query = options.Apply(query, opts.Opts...)
	}

	// 应用条件
	if opts.Condition != "" {
		query = query.Where(opts.Condition, opts.Args...)
	}

	return query
}

// Get 根据ID或条件获取单个实体
func (r *GenericRepo[T]) Get(ctx context.Context, id any, opts *QueryOptions) (*T, error) {
	var entity T

	// 创建查询并应用选项
	query := r.applyQueryOptions(r.DB.WithContext(ctx), opts)

	// 执行查询
	var err error
	if id != nil {
		// 根据ID查询
		err = query.First(&entity, id).Error
	} else if opts != nil && opts.Condition != "" {
		// 根据条件查询
		err = query.First(&entity).Error
	} else {
		// 没有ID和条件，返回错误
		return nil, errorx.ErrInvalidParams.WithMsg("查询参数不能为空")
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 使用ErrorCode创建特定的错误
			return nil, errorx.GetError(r.ErrorCode).WithMsg("记录不存在")
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

	// 应用查询选项
	query = r.applyQueryOptions(query, opts)

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
	query = r.applyQueryOptions(query, opts)

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
		ErrorCode: r.ErrorCode,
	}
}

// SetErrorCode 设置错误码
func (r *GenericRepo[T]) SetErrorCode(code int) *GenericRepo[T] {
	r.ErrorCode = code
	return r
}

// UpdateFields 更新实体字段
func (r *GenericRepo[T]) UpdateFields(ctx context.Context, id any, fields map[string]any) error {
	var entity T
	return r.DB.WithContext(ctx).Model(&entity).Where("id = ?", id).Updates(fields).Error
}

// FindByField 根据字段查询实体
func (r *GenericRepo[T]) FindByField(ctx context.Context, field string, value any) (*T, error) {
	opts := &QueryOptions{
		Condition: field + " = ?",
		Args:      []any{value},
	}
	return r.Get(ctx, nil, opts)
}

// GetByID 根据ID获取实体（便捷方法）
func (r *GenericRepo[T]) GetByID(ctx context.Context, id any) (*T, error) {
	return r.Get(ctx, id, nil)
}

// GetPage 分页查询（便捷方法）
func (r *GenericRepo[T]) GetPage(ctx context.Context, page, pageSize int, condition string, args ...any) ([]T, int64, error) {
	// 创建查询选项
	opts := &QueryOptions{
		Condition: condition,
		Args:      args,
	}

	// 获取列表
	entities, err := r.List(ctx, page, pageSize, opts)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := r.Count(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// BatchDelete 批量删除实体
func (r *GenericRepo[T]) BatchDelete(ctx context.Context, ids []any) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, ids).Error
}
