package repository

import (
	"context"

	"gorm.io/gorm"
)

// Aggregate 聚合查询
type Aggregate string

const (
	// Sum 求和
	Sum Aggregate = "SUM"
	// Avg 求平均
	Avg Aggregate = "AVG"
	// Max 求最大值
	Max Aggregate = "MAX"
	// Min 求最小值
	Min Aggregate = "MIN"
	// Count 求数量
	Count Aggregate = "COUNT"
)

// Entity 表示可以存储在数据库中的实体
// 它必须有一个 ID 字段，可以是 uint、int64 或 string 类型
type Entity interface {
	any
}

// CRUDRepository 基础CRUD操作接口
// 提供创建、读取、更新和删除实体的基本功能
type CRUDRepository[T Entity] interface {
	// Create 创建实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - entity: 要创建的实体指针
	// 返回:
	//   - error: 如果创建失败，返回错误
	Create(ctx context.Context, entity *T) error

	// GetByID 根据ID获取实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - id: 实体ID，可以是任何类型（通常是int64、uint或string）
	// 返回:
	//   - *T: 如果找到，返回实体指针
	//   - error: 如果查询失败或未找到实体，返回错误
	GetByID(ctx context.Context, id any) (*T, error)

	// Update 更新实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - entity: 要更新的实体指针
	// 返回:
	//   - error: 如果更新失败，返回错误
	Update(ctx context.Context, entity *T) error

	// Delete 删除实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - id: 要删除的实体ID
	// 返回:
	//   - error: 如果删除失败，返回错误
	Delete(ctx context.Context, id any) error

	// UpdateFields 更新实体字段
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - id: 要更新的实体ID
	//   - fields: 要更新的字段和值的映射
	// 返回:
	//   - error: 如果更新失败，返回错误
	UpdateFields(ctx context.Context, id any, fields map[string]any) error
}

// QueryRepository 查询操作接口
// 提供各种查询实体的方法
type QueryRepository[T Entity] interface {
	// List 获取实体列表（分页）
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - page: 页码，从1开始
	//   - pageSize: 每页大小
	// 返回:
	//   - []T: 实体列表
	//   - int64: 总记录数
	//   - error: 如果查询失败，返回错误
	List(ctx context.Context, page, pageSize int) ([]T, int64, error)

	// Count 获取实体总数
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	// 返回:
	//   - int64: 总记录数
	//   - error: 如果查询失败，返回错误
	Count(ctx context.Context) (int64, error)

	// FindByField 根据字段查询单个实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - field: 字段名
	//   - value: 字段值
	// 返回:
	//   - *T: 如果找到，返回实体指针
	//   - error: 如果查询失败或未找到实体，返回错误
	FindByField(ctx context.Context, field string, value any) (*T, error)

	// FindAllByField 根据字段查询多个实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - field: 字段名
	//   - value: 字段值
	// 返回:
	//   - []T: 实体列表
	//   - error: 如果查询失败，返回错误
	FindAllByField(ctx context.Context, field string, value any) ([]T, error)

	// FindByCondition 根据条件查询实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - condition: 查询条件（SQL WHERE子句）
	//   - args: 条件参数
	// 返回:
	//   - []T: 实体列表
	//   - error: 如果查询失败，返回错误
	FindByCondition(ctx context.Context, condition string, args ...any) ([]T, error)

	// FindOneByCondition 根据条件查询单个实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - condition: 查询条件（SQL WHERE子句）
	//   - args: 条件参数
	// 返回:
	//   - *T: 如果找到，返回实体指针
	//   - error: 如果查询失败或未找到实体，返回错误
	FindOneByCondition(ctx context.Context, condition string, args ...any) (*T, error)

	// GetPage 分页查询
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - page: 页码，从1开始
	//   - pageSize: 每页大小
	//   - condition: 查询条件（SQL WHERE子句）
	//   - args: 条件参数
	// 返回:
	//   - []T: 实体列表
	//   - int64: 总记录数
	//   - error: 如果查询失败，返回错误
	GetPage(ctx context.Context, page, pageSize int, condition string, args ...any) ([]T, int64, error)
}

// AdvancedQueryRepository 高级查询操作接口
// 提供更复杂的查询功能
type AdvancedQueryRepository[T Entity] interface {
	// FindWithLike 使用LIKE查询实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - field: 字段名
	//   - value: 模糊匹配值
	// 返回:
	//   - []T: 实体列表
	//   - error: 如果查询失败，返回错误
	FindWithLike(ctx context.Context, field string, value string) ([]T, error)

	// FindWithIn 使用IN查询实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - field: 字段名
	//   - values: 值列表
	// 返回:
	//   - []T: 实体列表
	//   - error: 如果查询失败，返回错误
	FindWithIn(ctx context.Context, field string, values []any) ([]T, error)

	// FindWithBetween 使用BETWEEN查询实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - field: 字段名
	//   - min: 最小值
	//   - max: 最大值
	// 返回:
	//   - []T: 实体列表
	//   - error: 如果查询失败，返回错误
	FindWithBetween(ctx context.Context, field string, min, max any) ([]T, error)

	// CountWithCondition 根据条件统计数量
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - condition: 查询条件（SQL WHERE子句）
	//   - args: 条件参数
	// 返回:
	//   - int64: 记录数
	//   - error: 如果查询失败，返回错误
	CountWithCondition(ctx context.Context, condition string, args ...any) (int64, error)

	// Exists 检查是否存在
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - condition: 查询条件（SQL WHERE子句）
	//   - args: 条件参数
	// 返回:
	//   - bool: 如果存在返回true，否则返回false
	//   - error: 如果查询失败，返回错误
	Exists(ctx context.Context, condition string, args ...any) (bool, error)
}

// AggregateRepository 聚合查询接口
// 提供数据聚合和分析功能
type AggregateRepository[T Entity] interface {
	// AggregateField 聚合查询字段
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - aggregate: 聚合函数（如COUNT、SUM、AVG等）
	//   - field: 字段名
	//   - condition: 查询条件（SQL WHERE子句）
	//   - args: 条件参数
	// 返回:
	//   - float64: 聚合结果
	//   - error: 如果查询失败，返回错误
	AggregateField(ctx context.Context, aggregate Aggregate, field string, condition string, args ...any) (float64, error)

	// GroupBy 分组查询
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - groupFields: 分组字段列表
	//   - selectFields: 选择字段列表
	//   - condition: 查询条件（SQL WHERE子句）
	//   - args: 条件参数
	// 返回:
	//   - []map[string]any: 查询结果
	//   - error: 如果查询失败，返回错误
	GroupBy(ctx context.Context, groupFields []string, selectFields []string, condition string, args ...any) ([]map[string]any, error)

	// Join 连接查询
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - joinType: 连接类型（如LEFT、RIGHT、INNER等）
	//   - table: 连接表名
	//   - on: 连接条件
	//   - selectFields: 选择字段列表
	//   - condition: 查询条件（SQL WHERE子句）
	//   - args: 条件参数
	// 返回:
	//   - []map[string]any: 查询结果
	//   - error: 如果查询失败，返回错误
	Join(ctx context.Context, joinType string, table string, on string, selectFields []string, condition string, args ...any) ([]map[string]any, error)
}

// BatchRepository 批量操作接口
// 提供批量处理实体的功能
type BatchRepository[T Entity] interface {
	// BatchCreate 批量创建实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - entities: 要创建的实体列表
	// 返回:
	//   - error: 如果创建失败，返回错误
	BatchCreate(ctx context.Context, entities []T) error

	// BatchDelete 批量删除实体
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - ids: 要删除的实体ID列表
	// 返回:
	//   - error: 如果删除失败，返回错误
	BatchDelete(ctx context.Context, ids []any) error
}

// TransactionRepository 事务操作接口
// 提供事务支持
type TransactionRepository[T Entity] interface {
	// Transaction 在事务中执行函数
	// 参数:
	//   - ctx: 上下文，用于取消操作和传递请求范围的值
	//   - fn: 在事务中执行的函数
	// 返回:
	//   - error: 如果事务执行失败，返回错误
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	// WithTx 使用事务
	// 参数:
	//   - tx: GORM事务对象
	// 返回:
	//   - Repository[T]: 使用事务的仓库实例
	WithTx(tx *gorm.DB) Repository[T]
}

// Repository 完整的仓库接口
// 组合所有子接口，提供完整的数据访问功能
type Repository[T Entity] interface {
	CRUDRepository[T]
	QueryRepository[T]
	AdvancedQueryRepository[T]
	AggregateRepository[T]
	BatchRepository[T]
	TransactionRepository[T]

	// SetErrorCode 设置NotFound错误的错误码
	// 参数:
	//   - code: 错误码
	// 返回:
	//   - Repository[T]: 仓库实例，用于链式调用
	SetErrorCode(code int) Repository[T]
}
