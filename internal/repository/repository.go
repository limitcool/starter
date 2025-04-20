package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository 通用仓库接口
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
	List(ctx context.Context, page, pageSize int) ([]T, int64, error)

	// Count 获取实体总数
	Count(ctx context.Context) (int64, error)

	// UpdateFields 更新实体字段
	UpdateFields(ctx context.Context, id any, fields map[string]any) error

	// BatchCreate 批量创建实体
	BatchCreate(ctx context.Context, entities []T) error

	// BatchDelete 批量删除实体
	BatchDelete(ctx context.Context, ids []any) error

	// FindByField 根据字段查询实体
	FindByField(ctx context.Context, field string, value any) (*T, error)

	// FindAllByField 根据字段查询多个实体
	FindAllByField(ctx context.Context, field string, value any) ([]T, error)

	// FindByCondition 根据条件查询实体
	FindByCondition(ctx context.Context, condition string, args ...any) ([]T, error)

	// FindOneByCondition 根据条件查询单个实体
	FindOneByCondition(ctx context.Context, condition string, args ...any) (*T, error)

	// GetPage 分页查询
	GetPage(ctx context.Context, page, pageSize int, condition string, args ...any) ([]T, int64, error)

	// FindWithLike 使用LIKE查询实体
	FindWithLike(ctx context.Context, field string, value string) ([]T, error)

	// FindWithIn 使用IN查询实体
	FindWithIn(ctx context.Context, field string, values []any) ([]T, error)

	// FindWithBetween 使用BETWEEN查询实体
	FindWithBetween(ctx context.Context, field string, min, max any) ([]T, error)

	// CountWithCondition 根据条件统计数量
	CountWithCondition(ctx context.Context, condition string, args ...any) (int64, error)

	// AggregateField 聚合查询字段
	AggregateField(ctx context.Context, aggregate Aggregate, field string, condition string, args ...any) (float64, error)

	// GroupBy 分组查询
	GroupBy(ctx context.Context, groupFields []string, selectFields []string, condition string, args ...any) ([]map[string]any, error)

	// Join 连接查询
	Join(ctx context.Context, joinType string, table string, on string, selectFields []string, condition string, args ...any) ([]map[string]any, error)

	// Exists 检查是否存在
	Exists(ctx context.Context, condition string, args ...any) (bool, error)

	// Raw 原生查询
	Raw(ctx context.Context, sql string, values ...any) ([]map[string]any, error)

	// Transaction 在事务中执行函数
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	// WithTx 使用事务
	WithTx(tx *gorm.DB) Repository[T]
}
