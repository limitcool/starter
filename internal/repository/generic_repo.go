package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"gorm.io/gorm"
)

// 注意：Entity接口已移至base_interfaces.go

// GenericRepo 通用仓库实现
type GenericRepo[T Entity] struct {
	DB        *gorm.DB
	TableName string
	ErrorCode int // 用于NotFound错误
}

// NewGenericRepo 创建通用仓库
// 返回实现了Repository[T]接口的实例
func NewGenericRepo[T Entity](db *gorm.DB) Repository[T] {
	// 创建一个零值的T类型实例，用于获取表名
	var entity T
	var tableName string

	// 尝试获取表名
	if tabler, ok := any(entity).(interface{ TableName() string }); ok {
		tableName = tabler.TableName()
	} else {
		// 使用反射获取类型名称作为表名
		t := reflect.TypeOf(entity)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		tableName = t.Name()
	}

	// 确保GenericRepo实现了Repository接口
	// 这是一个编译时检查，如果GenericRepo没有实现Repository接口，编译会失败
	var _ Repository[T] = (*GenericRepo[T])(nil)

	return &GenericRepo[T]{
		DB:        db,
		TableName: tableName,
		ErrorCode: errorx.ErrorNotFoundCode, // 默认错误码
	}
}

// SetErrorCode 设置NotFound错误的错误码
func (r *GenericRepo[T]) SetErrorCode(code int) Repository[T] {
	r.ErrorCode = code
	return r
}

// Create 创建实体
func (r *GenericRepo[T]) Create(ctx context.Context, entity *T) error {
	err := r.DB.WithContext(ctx).Create(entity).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("创建%s失败", r.TableName))
	}
	return nil
}

// GetByID 根据ID获取实体
func (r *GenericRepo[T]) GetByID(ctx context.Context, id any) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建一个错误对象
		notFoundErr := errorx.NewAppError(r.ErrorCode, fmt.Sprintf("%s ID %v 不存在", r.TableName, id), http.StatusNotFound)
		return nil, errorx.WrapError(notFoundErr, "")
	}

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询%s失败: id=%v", r.TableName, id))
	}

	return &entity, nil
}

// Update 更新实体
func (r *GenericRepo[T]) Update(ctx context.Context, entity *T) error {
	err := r.DB.WithContext(ctx).Save(entity).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新%s失败", r.TableName))
	}
	return nil
}

// Delete 删除实体
func (r *GenericRepo[T]) Delete(ctx context.Context, id any) error {
	var entity T
	err := r.DB.WithContext(ctx).Delete(&entity, id).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("删除%s失败: id=%v", r.TableName, id))
	}
	return nil
}

// normalizePagination 标准化分页参数
func normalizePagination(page, pageSize int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	return page, pageSize, offset
}

// List 获取实体列表
func (r *GenericRepo[T]) List(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	_, normalizedPageSize, offset := normalizePagination(page, pageSize)

	var entities []T
	var total int64

	// 获取总数
	if err := r.DB.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, errorx.WrapError(err, fmt.Sprintf("获取%s总数失败", r.TableName))
	}

	// 获取列表
	if err := r.DB.WithContext(ctx).Offset(offset).Limit(normalizedPageSize).Find(&entities).Error; err != nil {
		return nil, 0, errorx.WrapError(err, fmt.Sprintf("获取%s列表失败", r.TableName))
	}

	return entities, total, nil
}

// Count 获取实体总数
func (r *GenericRepo[T]) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.DB.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return 0, errorx.WrapError(err, fmt.Sprintf("获取%s总数失败", r.TableName))
	}
	return total, nil
}

// UpdateFields 更新实体字段
func (r *GenericRepo[T]) UpdateFields(ctx context.Context, id any, fields map[string]any) error {
	err := r.DB.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新%s字段失败: id=%v", r.TableName, id))
	}
	return nil
}

// FindByField 根据字段查询实体
func (r *GenericRepo[T]) FindByField(ctx context.Context, field string, value any) (*T, error) {
	var entity T
	// 使用 map 构建查询条件，避免 SQL 注入
	err := r.DB.WithContext(ctx).Where(map[string]any{field: value}).First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建一个错误对象
		notFoundErr := errorx.NewAppError(r.ErrorCode, fmt.Sprintf("%s %s=%v 不存在", r.TableName, field, value), http.StatusNotFound)
		return nil, errorx.WrapError(notFoundErr, "")
	}

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询%s失败: %s=%v", r.TableName, field, value))
	}

	return &entity, nil
}

// FindAllByField 根据字段查询多个实体
func (r *GenericRepo[T]) FindAllByField(ctx context.Context, field string, value any) ([]T, error) {
	var entities []T
	// 使用 map 构建查询条件，避免 SQL 注入
	err := r.DB.WithContext(ctx).Where(map[string]any{field: value}).Find(&entities).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("查询%s列表失败: %s=%v", r.TableName, field, value))
	}

	return entities, nil
}

// BatchCreate 批量创建实体
func (r *GenericRepo[T]) BatchCreate(ctx context.Context, entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	err := r.DB.WithContext(ctx).Create(&entities).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("批量创建%s失败", r.TableName))
	}
	return nil
}

// BatchDelete 批量删除实体
func (r *GenericRepo[T]) BatchDelete(ctx context.Context, ids []any) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.DB.WithContext(ctx).Where("id IN ?", ids).Delete(new(T)).Error
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("批量删除%s失败", r.TableName))
	}
	return nil
}

// FindByCondition 根据条件查询实体
// 安全使用指南：
// 1. 始终使用参数化查询（占位符）传递用户输入或变量
// 2. 正确示例：repo.FindByCondition(ctx, "name = ?", userInput)
// 3. 错误示例：repo.FindByCondition(ctx, "name = '" + userInput + "'") // 极危险！
// 4. 如果可能，优先使用FindByField等更安全的方法
func (r *GenericRepo[T]) FindByCondition(ctx context.Context, condition string, args ...any) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Where(condition, args...).Find(&entities).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("条件查询%s失败", r.TableName))
	}

	return entities, nil
}

// FindByMap 使用map条件查询实体（更安全，防止SQL注入）
// 例如：repo.FindByMap(ctx, map[string]any{"status": 1, "type": "user"})
func (r *GenericRepo[T]) FindByMap(ctx context.Context, conditions map[string]any) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Where(conditions).Find(&entities).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("条件查询%s失败", r.TableName))
	}

	return entities, nil
}

// FindOneByCondition 根据条件查询单个实体
// 安全使用指南：
// 1. 始终使用参数化查询（占位符）传递用户输入或变量
// 2. 正确示例：repo.FindOneByCondition(ctx, "name = ?", userInput)
// 3. 错误示例：repo.FindOneByCondition(ctx, "name = '" + userInput + "'") // 极危险！
// 4. 如果可能，优先使用FindByField等更安全的方法
func (r *GenericRepo[T]) FindOneByCondition(ctx context.Context, condition string, args ...any) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where(condition, args...).First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建一个错误对象
		notFoundErr := errorx.NewAppError(r.ErrorCode, fmt.Sprintf("%s不存在", r.TableName), http.StatusNotFound)
		return nil, errorx.WrapError(notFoundErr, "")
	}

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("条件查询%s失败", r.TableName))
	}

	return &entity, nil
}

// FindOneByMap 使用map条件查询单个实体（更安全，防止SQL注入）
// 例如：repo.FindOneByMap(ctx, map[string]any{"username": username})
func (r *GenericRepo[T]) FindOneByMap(ctx context.Context, conditions map[string]any) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where(conditions).First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建一个错误对象
		notFoundErr := errorx.NewAppError(r.ErrorCode, fmt.Sprintf("%s不存在", r.TableName), http.StatusNotFound)
		return nil, errorx.WrapError(notFoundErr, "")
	}

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("条件查询%s失败", r.TableName))
	}

	return &entity, nil
}

// GetPage 分页查询
// 安全使用指南：
// 1. 始终使用参数化查询（占位符）传递用户输入或变量
// 2. 正确示例：repo.GetPage(ctx, 1, 10, "status = ?", userInput)
// 3. 错误示例：repo.GetPage(ctx, 1, 10, "status = '" + userInput + "'") // 极危险！
func (r *GenericRepo[T]) GetPage(ctx context.Context, page, pageSize int, condition string, args ...any) ([]T, int64, error) {
	// 标准化分页参数
	_, normalizedPageSize, offset := normalizePagination(page, pageSize)

	// 构建查询
	tx := r.DB.WithContext(ctx).Model(new(T))

	// 添加条件
	if condition != "" {
		tx = tx.Where(condition, args...)
	}

	// 获取总数
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, errorx.WrapError(err, fmt.Sprintf("获取%s总数失败", r.TableName))
	}

	// 执行查询
	var entities []T
	if err := tx.Offset(offset).Limit(normalizedPageSize).Find(&entities).Error; err != nil {
		return nil, 0, errorx.WrapError(err, fmt.Sprintf("分页查询%s失败", r.TableName))
	}

	return entities, total, nil
}

// GetPageByMap 使用map条件进行分页查询（更安全，防止SQL注入）
// 例如：repo.GetPageByMap(ctx, 1, 10, map[string]any{"status": 1, "type": "user"})
func (r *GenericRepo[T]) GetPageByMap(ctx context.Context, page, pageSize int, conditions map[string]any) ([]T, int64, error) {
	// 标准化分页参数
	_, normalizedPageSize, offset := normalizePagination(page, pageSize)

	// 构建查询
	tx := r.DB.WithContext(ctx).Model(new(T))

	// 添加条件
	if len(conditions) > 0 {
		tx = tx.Where(conditions)
	}

	// 获取总数
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, errorx.WrapError(err, fmt.Sprintf("获取%s总数失败", r.TableName))
	}

	// 执行查询
	var entities []T
	if err := tx.Offset(offset).Limit(normalizedPageSize).Find(&entities).Error; err != nil {
		return nil, 0, errorx.WrapError(err, fmt.Sprintf("分页查询%s失败", r.TableName))
	}

	return entities, total, nil
}

// Transaction 在事务中执行函数
func (r *GenericRepo[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.DB.WithContext(ctx).Transaction(fn)
}

// FindWithLike 使用LIKE查询实体
func (r *GenericRepo[T]) FindWithLike(ctx context.Context, field string, value string) ([]T, error) {
	var entities []T
	// 使用条件表达式和参数化查询，避免 SQL 注入
	// 注意：对于 LIKE 查询，我们使用条件表达式和参数化查询，而不是 map
	// 因为 map 不支持 LIKE 操作符
	condition := map[string]any{}
	condition[field] = gorm.Expr("LIKE ?", "%"+value+"%")
	err := r.DB.WithContext(ctx).Where(condition).Find(&entities).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("模糊查询%s失败: %s=%s", r.TableName, field, value))
	}

	return entities, nil
}

// FindWithIn 使用IN查询实体
func (r *GenericRepo[T]) FindWithIn(ctx context.Context, field string, values []any) ([]T, error) {
	if len(values) == 0 {
		return []T{}, nil
	}

	var entities []T
	// 使用条件表达式和参数化查询，避免 SQL 注入
	condition := map[string]any{}
	condition[field] = values
	err := r.DB.WithContext(ctx).Where(condition).Find(&entities).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("使用IN查询%s失败: %s", r.TableName, field))
	}

	return entities, nil
}

// FindWithBetween 使用BETWEEN查询实体
func (r *GenericRepo[T]) FindWithBetween(ctx context.Context, field string, min, max any) ([]T, error) {
	var entities []T
	// 使用条件表达式和参数化查询，避免 SQL 注入
	condition := map[string]any{}
	condition[field] = gorm.Expr("BETWEEN ? AND ?", min, max)
	err := r.DB.WithContext(ctx).Where(condition).Find(&entities).Error

	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("使用BETWEEN查询%s失败: %s", r.TableName, field))
	}

	return entities, nil
}

// Count 根据条件统计数量
func (r *GenericRepo[T]) CountWithCondition(ctx context.Context, condition string, args ...any) (int64, error) {
	var count int64
	tx := r.DB.WithContext(ctx).Model(new(T))

	if condition != "" {
		tx = tx.Where(condition, args...)
	}

	err := tx.Count(&count).Error
	if err != nil {
		return 0, errorx.WrapError(err, fmt.Sprintf("统计%s数量失败", r.TableName))
	}

	return count, nil
}

// WithTx 使用事务
func (r *GenericRepo[T]) WithTx(tx *gorm.DB) Repository[T] {
	repo := &GenericRepo[T]{
		DB:        tx,
		TableName: r.TableName,
		ErrorCode: r.ErrorCode,
	}

	// 确保GenericRepo实现了Repository接口
	// 这是一个编译时检查，如果GenericRepo没有实现Repository接口，编译会失败
	var _ Repository[T] = (*GenericRepo[T])(nil)
	var _ CRUDRepository[T] = (*GenericRepo[T])(nil)
	var _ QueryRepository[T] = (*GenericRepo[T])(nil)
	var _ AdvancedQueryRepository[T] = (*GenericRepo[T])(nil)
	var _ AggregateRepository[T] = (*GenericRepo[T])(nil)
	var _ BatchRepository[T] = (*GenericRepo[T])(nil)
	var _ TransactionRepository[T] = (*GenericRepo[T])(nil)

	return repo
}

// 注意：Aggregate类型已移至base_interfaces.go

// AggregateField 聚合查询字段
func (r *GenericRepo[T]) AggregateField(ctx context.Context, aggregate Aggregate, field string, condition string, args ...any) (float64, error) {
	// 构建查询
	tx := r.DB.WithContext(ctx).Model(new(T))

	// 添加条件
	if condition != "" {
		tx = tx.Where(condition, args...)
	}

	// 执行聚合查询
	var result float64
	query := fmt.Sprintf("%s(%s)", aggregate, field)
	if err := tx.Select(query).Scan(&result).Error; err != nil {
		return 0, errorx.WrapError(err, fmt.Sprintf("聚合查询%s失败: %s(%s)", r.TableName, aggregate, field))
	}

	return result, nil
}

// GroupBy 分组查询
func (r *GenericRepo[T]) GroupBy(ctx context.Context, groupFields []string, selectFields []string, condition string, args ...any) ([]map[string]any, error) {
	// 构建查询
	tx := r.DB.WithContext(ctx).Model(new(T))

	// 添加条件
	if condition != "" {
		tx = tx.Where(condition, args...)
	}

	// 添加分组
	tx = tx.Group(strings.Join(groupFields, ", "))

	// 添加选择字段
	if len(selectFields) > 0 {
		tx = tx.Select(selectFields)
	}

	// 执行查询
	var results []map[string]any
	if err := tx.Find(&results).Error; err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("分组查询%s失败", r.TableName))
	}

	return results, nil
}

// Join 连接查询
func (r *GenericRepo[T]) Join(ctx context.Context, joinType string, table string, on string, selectFields []string, condition string, args ...any) ([]map[string]any, error) {
	// 构建查询
	tx := r.DB.WithContext(ctx).Model(new(T))

	// 添加连接
	joinSQL := fmt.Sprintf("%s JOIN %s ON %s", joinType, table, on)
	tx = tx.Joins(joinSQL)

	// 添加条件
	if condition != "" {
		tx = tx.Where(condition, args...)
	}

	// 添加选择字段
	if len(selectFields) > 0 {
		tx = tx.Select(selectFields)
	}

	// 执行查询
	var results []map[string]any
	if err := tx.Find(&results).Error; err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("连接查询%s失败", r.TableName))
	}

	return results, nil
}

// Exists 检查是否存在
func (r *GenericRepo[T]) Exists(ctx context.Context, condition string, args ...any) (bool, error) {
	// 构建查询
	tx := r.DB.WithContext(ctx).Model(new(T))

	// 添加条件
	if condition != "" {
		tx = tx.Where(condition, args...)
	}

	// 执行查询
	var count int64
	if err := tx.Limit(1).Count(&count).Error; err != nil {
		return false, errorx.WrapError(err, fmt.Sprintf("检查%s是否存在失败", r.TableName))
	}

	return count > 0, nil
}

// 注意：Raw方法已被移除，因为它存在潜在的SQL注入风险
// 请使用更安全的方法如FindByField等进行查询
// 如果需要执行复杂查询，请考虑以下替代方案：
// 1. 使用GORM的高级查询功能，如Joins、Preload等
// 2. 创建自定义的查询方法，使用参数化查询
// 3. 使用视图或存储过程
