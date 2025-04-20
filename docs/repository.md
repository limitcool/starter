# 仓库层设计文档

本文档描述了项目中仓库层的设计和使用方法。

## 1. 概述

仓库层是应用程序与数据库之间的抽象层，负责处理数据的持久化和检索。我们使用泛型和接口来实现一个灵活、可扩展的仓库层设计。

主要特点：

- 使用 Go 1.18+ 泛型特性，减少重复代码
- 提供统一的 CRUD 操作接口
- 支持高级查询功能
- 支持缓存机制
- 支持事务处理

## 2. 核心组件

### 2.1 Repository 接口

`Repository` 接口定义了所有仓库应该实现的方法：

```go
type Repository[T Entity] interface {
    // 基本 CRUD 操作
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id any) (*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id any) error
    
    // 列表和分页
    List(ctx context.Context, page, pageSize int) ([]T, int64, error)
    GetPage(ctx context.Context, page, pageSize int, condition string, args ...any) ([]T, int64, error)
    
    // 字段操作
    UpdateFields(ctx context.Context, id any, fields map[string]any) error
    
    // 批量操作
    BatchCreate(ctx context.Context, entities []T) error
    BatchDelete(ctx context.Context, ids []any) error
    
    // 条件查询
    FindByField(ctx context.Context, field string, value any) (*T, error)
    FindAllByField(ctx context.Context, field string, value any) ([]T, error)
    FindByCondition(ctx context.Context, condition string, args ...any) ([]T, error)
    FindOneByCondition(ctx context.Context, condition string, args ...any) (*T, error)
    
    // 高级查询
    FindWithLike(ctx context.Context, field string, value string) ([]T, error)
    FindWithIn(ctx context.Context, field string, values []any) ([]T, error)
    FindWithBetween(ctx context.Context, field string, min, max any) ([]T, error)
    
    // 聚合查询
    CountWithCondition(ctx context.Context, condition string, args ...any) (int64, error)
    AggregateField(ctx context.Context, aggregate Aggregate, field string, condition string, args ...any) (float64, error)
    GroupBy(ctx context.Context, groupFields []string, selectFields []string, condition string, args ...any) ([]map[string]any, error)
    
    // 连接查询
    Join(ctx context.Context, joinType string, table string, on string, selectFields []string, condition string, args ...any) ([]map[string]any, error)
    
    // 其他查询
    Exists(ctx context.Context, condition string, args ...any) (bool, error)
    Raw(ctx context.Context, sql string, values ...any) ([]map[string]any, error)
    
    // 事务支持
    Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
    WithTx(tx *gorm.DB) Repository[T]
}
```

### 2.2 GenericRepo 实现

`GenericRepo` 是 `Repository` 接口的通用实现，使用泛型来处理不同类型的实体：

```go
type GenericRepo[T Entity] struct {
    DB        *gorm.DB
    TableName string
    ErrorCode int // 用于NotFound错误
}
```

### 2.3 CachedRepo 缓存实现

`CachedRepo` 是带缓存功能的仓库实现，它包装了一个基础仓库，并添加了缓存层：

```go
type CachedRepo[T Entity] struct {
    repo       Repository[T]
    cache      cache.Cache
    prefix     string
    expiration time.Duration
}
```

## 3. 使用方法

### 3.1 创建基础仓库

```go
// 创建用户仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
    genericRepo := NewGenericRepo[model.User](db)
    genericRepo.SetErrorCode(errorx.ErrorUserNotFoundCode) // 设置错误码
    
    return &UserRepo{
        DB:          db,
        genericRepo: genericRepo,
    }
}
```

### 3.2 使用泛型仓库

```go
// GetByID 根据ID获取用户
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
    // 使用泛型仓库
    return r.genericRepo.GetByID(ctx, id)
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
    // 使用泛型仓库的高级查询
    return r.genericRepo.FindByField(ctx, "username", username)
}
```

### 3.3 添加缓存支持

```go
// 创建缓存工厂
cacheFactory := cache.GetFactory()

// 创建内存缓存
userCache, err := cacheFactory.Create("user_cache", cache.Memory,
    cache.WithExpiration(5*time.Minute),
    cache.WithMaxEntries(1000),
)
if err != nil {
    panic(err)
}

// 创建基础用户仓库
baseUserRepo := repository.NewUserRepo(db)

// 创建缓存用户仓库
cachedUserRepo := repository.NewCachedRepo(baseUserRepo, userCache, "user", 5*time.Minute)

// 或者使用辅助函数
cachedUserRepo, err := repository.WithCache(baseUserRepo, "user_cache", "user", 5*time.Minute)
if err != nil {
    panic(err)
}
```

### 3.4 使用高级查询

```go
// 模糊查询
users, err := userRepo.FindWithLike(ctx, "username", "admin")

// IN 查询
users, err := userRepo.FindWithIn(ctx, "status", []any{1, 2, 3})

// BETWEEN 查询
users, err := userRepo.FindWithBetween(ctx, "created_at", startTime, endTime)

// 聚合查询
count, err := userRepo.AggregateField(ctx, repository.Count, "id", "status = ?", 1)
sum, err := userRepo.AggregateField(ctx, repository.Sum, "amount", "user_id = ?", userID)

// 分组查询
results, err := userRepo.GroupBy(ctx, 
    []string{"status"}, 
    []string{"status", "COUNT(*) as count"}, 
    "created_at > ?", 
    startTime,
)

// 连接查询
results, err := userRepo.Join(ctx, 
    "LEFT", 
    "orders", 
    "users.id = orders.user_id", 
    []string{"users.id", "users.username", "COUNT(orders.id) as order_count"}, 
    "users.status = ?", 
    1,
)
```

### 3.5 使用事务

```go
err := userRepo.Transaction(ctx, func(tx *gorm.DB) error {
    // 在事务中使用仓库
    txRepo := userRepo.WithTx(tx)
    
    // 创建用户
    if err := txRepo.Create(ctx, user); err != nil {
        return err
    }
    
    // 创建用户配置
    if err := configRepo.WithTx(tx).Create(ctx, config); err != nil {
        return err
    }
    
    return nil
})
```

## 4. 最佳实践

### 4.1 仓库层设计原则

1. **单一职责原则**：每个仓库只负责一种实体的数据访问
2. **接口隔离原则**：使用接口定义仓库行为，便于测试和替换实现
3. **依赖注入**：通过构造函数注入数据库连接
4. **错误处理**：使用 errorx 包统一处理错误，包含详细的上下文信息
5. **事务支持**：提供事务支持，确保数据一致性

### 4.2 性能优化

1. **使用缓存**：对频繁访问的数据使用缓存
2. **批量操作**：使用批量操作减少数据库交互
3. **索引优化**：确保查询条件有适当的索引
4. **延迟加载**：使用延迟加载避免加载不必要的数据
5. **连接池管理**：合理配置连接池参数

### 4.3 测试策略

1. **单元测试**：使用模拟对象测试仓库逻辑
2. **集成测试**：使用测试数据库测试实际数据库交互
3. **性能测试**：测试仓库方法的性能和并发性
4. **事务测试**：测试事务的正确性和回滚机制

## 5. 示例

### 5.1 基本用法

```go
// 创建用户仓库
userRepo := repository.NewUserRepo(db)

// 创建用户
user := &model.User{
    Username: "admin",
    Password: "password",
    Email:    "admin@example.com",
}
if err := userRepo.Create(ctx, user); err != nil {
    // 处理错误
}

// 获取用户
user, err := userRepo.GetByID(ctx, 1)
if err != nil {
    // 处理错误
}

// 更新用户
user.Email = "new_email@example.com"
if err := userRepo.Update(ctx, user); err != nil {
    // 处理错误
}

// 删除用户
if err := userRepo.Delete(ctx, 1); err != nil {
    // 处理错误
}
```

### 5.2 高级查询

```go
// 分页查询
users, total, err := userRepo.GetPage(ctx, 1, 10, "status = ?", 1)
if err != nil {
    // 处理错误
}
fmt.Printf("总数: %d, 当前页数据: %d\n", total, len(users))

// 聚合查询
activeCount, err := userRepo.AggregateField(ctx, repository.Count, "id", "status = ?", 1)
if err != nil {
    // 处理错误
}
fmt.Printf("活跃用户数: %f\n", activeCount)

// 分组查询
results, err := userRepo.GroupBy(ctx, 
    []string{"status"}, 
    []string{"status", "COUNT(*) as count"}, 
    "", 
)
if err != nil {
    // 处理错误
}
for _, result := range results {
    fmt.Printf("状态: %v, 数量: %v\n", result["status"], result["count"])
}
```

### 5.3 缓存用法

```go
// 创建缓存用户仓库
cachedUserRepo, err := repository.WithCache(userRepo, "user_cache", "user", 5*time.Minute)
if err != nil {
    // 处理错误
}

// 获取用户（首次会从数据库获取并缓存）
user, err := cachedUserRepo.GetByID(ctx, 1)
if err != nil {
    // 处理错误
}

// 再次获取用户（从缓存获取）
user, err = cachedUserRepo.GetByID(ctx, 1)
if err != nil {
    // 处理错误
}

// 更新用户（会自动更新缓存）
user.Email = "updated@example.com"
if err := cachedUserRepo.Update(ctx, user); err != nil {
    // 处理错误
}

// 删除用户（会自动删除缓存）
if err := cachedUserRepo.Delete(ctx, 1); err != nil {
    // 处理错误
}
```

## 6. 常见问题

### 6.1 如何选择合适的缓存策略？

- **读多写少的数据**：使用较长的缓存过期时间
- **频繁变化的数据**：使用较短的缓存过期时间或不缓存
- **大型数据集**：考虑分片缓存或只缓存关键字段
- **关键业务数据**：确保缓存与数据库的一致性

### 6.2 如何处理缓存一致性问题？

- **写入时更新缓存**：在更新数据库时同时更新缓存
- **删除而非更新**：在数据变更时删除缓存，而不是更新缓存
- **设置合理的过期时间**：避免缓存长时间不一致
- **使用事务**：在事务中同时更新数据库和缓存

### 6.3 如何优化批量操作？

- **使用事务**：在事务中执行批量操作
- **分批处理**：将大批量操作分成小批次处理
- **使用原生SQL**：对于复杂的批量操作，考虑使用原生SQL
- **并发处理**：对于独立的批量操作，考虑并发处理

### 6.4 如何处理复杂查询？

- **使用原生SQL**：对于非常复杂的查询，使用 Raw 方法
- **组合查询**：将复杂查询拆分为多个简单查询
- **使用视图**：考虑在数据库中创建视图
- **预计算**：对于频繁使用的复杂统计，考虑预计算和缓存结果
