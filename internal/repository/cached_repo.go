package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/limitcool/starter/internal/pkg/cache"
	"gorm.io/gorm"
)

// CachedRepo 缓存仓库
type CachedRepo[T Entity] struct {
	repo  Repository[T]
	cache cache.Cache
	// 缓存前缀，用于区分不同实体的缓存
	prefix string
	// 缓存过期时间
	expiration time.Duration
}

// NewCachedRepo 创建缓存仓库
func NewCachedRepo[T Entity](repo Repository[T], cache cache.Cache, prefix string, expiration time.Duration) Repository[T] {
	return &CachedRepo[T]{
		repo:       repo,
		cache:      cache,
		prefix:     prefix,
		expiration: expiration,
	}
}

// generateKey 生成缓存键
func (r *CachedRepo[T]) generateKey(key string) string {
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

// Create 创建实体
func (r *CachedRepo[T]) Create(ctx context.Context, entity *T) error {
	// 调用底层仓库创建实体
	if err := r.repo.Create(ctx, entity); err != nil {
		return err
	}

	// 缓存实体
	return r.cacheEntity(ctx, entity)
}

// GetByID 根据ID获取实体
func (r *CachedRepo[T]) GetByID(ctx context.Context, id any) (*T, error) {
	// 尝试从缓存获取
	key := r.generateKey(fmt.Sprintf("id:%v", id))
	data, err := r.cache.Get(ctx, key)
	if err == nil {
		// 缓存命中，解析实体
		var entity T
		if err := json.Unmarshal(data, &entity); err != nil {
			return nil, err
		}
		return &entity, nil
	}

	// 缓存未命中，从底层仓库获取
	entity, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 缓存实体
	if err := r.cacheEntity(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

// Update 更新实体
func (r *CachedRepo[T]) Update(ctx context.Context, entity *T) error {
	// 调用底层仓库更新实体
	if err := r.repo.Update(ctx, entity); err != nil {
		return err
	}

	// 缓存实体
	return r.cacheEntity(ctx, entity)
}

// Delete 删除实体
func (r *CachedRepo[T]) Delete(ctx context.Context, id any) error {
	// 调用底层仓库删除实体
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 删除缓存
	key := r.generateKey(fmt.Sprintf("id:%v", id))
	return r.cache.Delete(ctx, key)
}

// List 获取实体列表
func (r *CachedRepo[T]) List(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	// 直接调用底层仓库，不缓存列表
	return r.repo.List(ctx, page, pageSize)
}

// Count 获取实体总数
func (r *CachedRepo[T]) Count(ctx context.Context) (int64, error) {
	// 尝试从缓存获取
	key := r.generateKey("count")
	data, err := r.cache.Get(ctx, key)
	if err == nil {
		// 缓存命中，解析总数
		var count int64
		if err := json.Unmarshal(data, &count); err != nil {
			return 0, err
		}
		return count, nil
	}

	// 缓存未命中，从底层仓库获取
	count, err := r.repo.Count(ctx)
	if err != nil {
		return 0, err
	}

	// 缓存总数
	countData, err := json.Marshal(count)
	if err != nil {
		return 0, err
	}
	if err := r.cache.Set(ctx, key, countData, r.expiration); err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateFields 更新实体字段
func (r *CachedRepo[T]) UpdateFields(ctx context.Context, id any, fields map[string]any) error {
	// 调用底层仓库更新实体字段
	if err := r.repo.UpdateFields(ctx, id, fields); err != nil {
		return err
	}

	// 删除缓存，下次获取时会重新缓存
	key := r.generateKey(fmt.Sprintf("id:%v", id))
	return r.cache.Delete(ctx, key)
}

// BatchCreate 批量创建实体
func (r *CachedRepo[T]) BatchCreate(ctx context.Context, entities []T) error {
	// 调用底层仓库批量创建实体
	if err := r.repo.BatchCreate(ctx, entities); err != nil {
		return err
	}

	// 批量缓存实体
	for i := range entities {
		if err := r.cacheEntity(ctx, &entities[i]); err != nil {
			return err
		}
	}

	return nil
}

// BatchDelete 批量删除实体
func (r *CachedRepo[T]) BatchDelete(ctx context.Context, ids []any) error {
	// 调用底层仓库批量删除实体
	if err := r.repo.BatchDelete(ctx, ids); err != nil {
		return err
	}

	// 批量删除缓存
	for _, id := range ids {
		key := r.generateKey(fmt.Sprintf("id:%v", id))
		if err := r.cache.Delete(ctx, key); err != nil {
			return err
		}
	}

	return nil
}

// FindByField 根据字段查询实体
func (r *CachedRepo[T]) FindByField(ctx context.Context, field string, value any) (*T, error) {
	// 尝试从缓存获取
	key := r.generateKey(fmt.Sprintf("%s:%v", field, value))
	data, err := r.cache.Get(ctx, key)
	if err == nil {
		// 缓存命中，解析实体
		var entity T
		if err := json.Unmarshal(data, &entity); err != nil {
			return nil, err
		}
		return &entity, nil
	}

	// 缓存未命中，从底层仓库获取
	entity, err := r.repo.FindByField(ctx, field, value)
	if err != nil {
		return nil, err
	}

	// 缓存实体
	if err := r.cacheEntity(ctx, entity); err != nil {
		return nil, err
	}

	// 缓存字段索引
	fieldKey := r.generateKey(fmt.Sprintf("%s:%v", field, value))
	idKey := r.generateKey(fmt.Sprintf("id:%v", getEntityID(entity)))
	if err := r.cache.Set(ctx, fieldKey, []byte(idKey), r.expiration); err != nil {
		return nil, err
	}

	return entity, nil
}

// FindAllByField 根据字段查询多个实体
func (r *CachedRepo[T]) FindAllByField(ctx context.Context, field string, value any) ([]T, error) {
	// 直接调用底层仓库，不缓存列表
	return r.repo.FindAllByField(ctx, field, value)
}

// FindByCondition 根据条件查询实体
func (r *CachedRepo[T]) FindByCondition(ctx context.Context, condition string, args ...any) ([]T, error) {
	// 直接调用底层仓库，不缓存列表
	return r.repo.FindByCondition(ctx, condition, args...)
}

// FindOneByCondition 根据条件查询单个实体
func (r *CachedRepo[T]) FindOneByCondition(ctx context.Context, condition string, args ...any) (*T, error) {
	// 直接调用底层仓库，不缓存
	return r.repo.FindOneByCondition(ctx, condition, args...)
}

// GetPage 分页查询
func (r *CachedRepo[T]) GetPage(ctx context.Context, page, pageSize int, condition string, args ...any) ([]T, int64, error) {
	// 直接调用底层仓库，不缓存列表
	return r.repo.GetPage(ctx, page, pageSize, condition, args...)
}

// FindWithLike 使用LIKE查询实体
func (r *CachedRepo[T]) FindWithLike(ctx context.Context, field string, value string) ([]T, error) {
	// 直接调用底层仓库，不缓存列表
	return r.repo.FindWithLike(ctx, field, value)
}

// FindWithIn 使用IN查询实体
func (r *CachedRepo[T]) FindWithIn(ctx context.Context, field string, values []any) ([]T, error) {
	// 直接调用底层仓库，不缓存列表
	return r.repo.FindWithIn(ctx, field, values)
}

// FindWithBetween 使用BETWEEN查询实体
func (r *CachedRepo[T]) FindWithBetween(ctx context.Context, field string, min, max any) ([]T, error) {
	// 直接调用底层仓库，不缓存列表
	return r.repo.FindWithBetween(ctx, field, min, max)
}

// CountWithCondition 根据条件统计数量
func (r *CachedRepo[T]) CountWithCondition(ctx context.Context, condition string, args ...any) (int64, error) {
	// 直接调用底层仓库，不缓存
	return r.repo.CountWithCondition(ctx, condition, args...)
}

// AggregateField 聚合查询字段
func (r *CachedRepo[T]) AggregateField(ctx context.Context, aggregate Aggregate, field string, condition string, args ...any) (float64, error) {
	// 直接调用底层仓库，不缓存
	return r.repo.AggregateField(ctx, aggregate, field, condition, args...)
}

// GroupBy 分组查询
func (r *CachedRepo[T]) GroupBy(ctx context.Context, groupFields []string, selectFields []string, condition string, args ...any) ([]map[string]any, error) {
	// 直接调用底层仓库，不缓存
	return r.repo.GroupBy(ctx, groupFields, selectFields, condition, args...)
}

// Join 连接查询
func (r *CachedRepo[T]) Join(ctx context.Context, joinType string, table string, on string, selectFields []string, condition string, args ...any) ([]map[string]any, error) {
	// 直接调用底层仓库，不缓存
	return r.repo.Join(ctx, joinType, table, on, selectFields, condition, args...)
}

// Exists 检查是否存在
func (r *CachedRepo[T]) Exists(ctx context.Context, condition string, args ...any) (bool, error) {
	// 直接调用底层仓库，不缓存
	return r.repo.Exists(ctx, condition, args...)
}

// 注意：Raw方法已被移除，因为它存在潜在的SQL注入风险
// 请使用更安全的方法如FindByField等进行查询

// Transaction 在事务中执行函数
func (r *CachedRepo[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	// 调用底层仓库的事务
	return r.repo.Transaction(ctx, fn)
}

// WithTx 使用事务
func (r *CachedRepo[T]) WithTx(tx *gorm.DB) Repository[T] {
	// 创建一个新的缓存仓库，使用底层仓库的事务
	return &CachedRepo[T]{
		repo:       r.repo.WithTx(tx),
		cache:      r.cache,
		prefix:     r.prefix,
		expiration: r.expiration,
	}
}

// SetErrorCode 设置NotFound错误的错误码
// 实现Repository接口
func (r *CachedRepo[T]) SetErrorCode(code int) Repository[T] {
	// 将错误码设置传递给底层仓库
	r.repo.SetErrorCode(code)
	return r
}

// cacheEntity 缓存实体
func (r *CachedRepo[T]) cacheEntity(ctx context.Context, entity *T) error {
	// 序列化实体
	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	// 缓存实体
	key := r.generateKey(fmt.Sprintf("id:%v", getEntityID(entity)))
	return r.cache.Set(ctx, key, data, r.expiration)
}

// getEntityID 获取实体ID
func getEntityID[T Entity](entity *T) any {
	// 使用反射获取ID字段
	v := reflect.ValueOf(entity).Elem()
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return nil
	}
	return idField.Interface()
}

// WithCache 为仓库添加缓存支持
func WithCache[T Entity](repo Repository[T], cacheName string, prefix string, expiration time.Duration) (Repository[T], error) {
	return WithCacheContext(context.Background(), repo, cacheName, prefix, expiration)
}

// WithCacheContext 使用上下文为仓库添加缓存支持
func WithCacheContext[T Entity](ctx context.Context, repo Repository[T], cacheName string, prefix string, expiration time.Duration) (Repository[T], error) {
	// 获取缓存
	cacheInstance, err := cache.GetFactory().GetWithContext(ctx, cacheName)
	if err != nil {
		// 如果缓存不存在，创建一个新的内存缓存
		cacheInstance, err = cache.GetFactory().CreateWithContext(ctx, cacheName, cache.Memory,
			cache.WithExpiration(expiration),
			cache.WithMaxEntries(10000),
		)
		if err != nil {
			return nil, err
		}
	}

	// 创建缓存仓库
	return NewCachedRepo(repo, cacheInstance, prefix, expiration), nil
}
