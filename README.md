# Starter

[![Go Reference](https://pkg.go.dev/badge/github.com/limitcool/starter.svg)](https://pkg.go.dev/github.com/limitcool/starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/limitcool/starter)](https://goreportcard.com/report/github.com/limitcool/starter)

## 特征
- 提供 gin 框架项目模版
- 集成 GORM 进行 ORM 映射和数据库操作
  - 支持 PostgreSQL (使用 pgx 驱动)
  - 支持 MySQL
  - 支持 SQLite
  - 提供丰富的查询选项工具函数
- 集成 Viper 进行配置管理
- 提供常用 gin 中间件和工具
  - 跨域中间件:处理 API 跨域请求,实现 CORS 支持
  - jwt 解析中间件:从请求中解析并验证 JWT Token,用于 API 身份认证
- 使用 Cobra 命令行框架，提供清晰的子命令结构
- 支持数据库迁移与服务器启动分离，提高启动速度
- 完善的数据库迁移系统，支持版本控制和回滚
- 内置用户、角色、权限和菜单管理系统

## 快速开始

```bash
go install github.com/go-eagle/eagle/cmd/eagle@latest
eagle new <project name> -r https://github.com/limitcool/starter -b main
```

## 使用方法

应用使用Cobra命令行框架，提供了更清晰的子命令结构。

### 基本命令

```bash
# 查看帮助信息
./<app-name> --help

# 查看版本信息
./<app-name> version

# 启动服务器
./<app-name> server

# 执行数据库迁移
./<app-name> migrate
```

### 服务器命令

服务器命令用于启动HTTP服务：

```bash
# 使用默认配置启动服务器
./<app-name> server

# 指定端口启动服务器
./<app-name> server --port 9000

# 使用指定配置文件启动服务器
./<app-name> server --config custom.yaml
```

### 数据库迁移命令

数据库迁移命令用于初始化或更新数据库结构：

```bash
# 执行数据库迁移
./<app-name> migrate

# 使用指定配置文件执行迁移
./<app-name> migrate --config prod.yaml

# 在迁移前清空数据库（危险操作）
./<app-name> migrate --fresh

# 回滚上一批数据库迁移
./<app-name> migrate rollback

# 显示数据库迁移状态
./<app-name> migrate status

# 重置所有数据库迁移
./<app-name> migrate reset
```

## 数据库迁移系统

本项目实现了一个完整的数据库迁移系统，用于管理数据库表结构的创建、更新和回滚。

### 迁移系统特点

- 支持按版本号顺序执行迁移
- 跟踪已执行的迁移记录
- 支持事务性迁移，确保数据一致性
- 提供向上和向下迁移功能
- 支持批次回滚和完全重置

### 迁移文件结构

迁移定义在 `internal/migration/migrations.go` 文件中，遵循以下结构：

```go
migrator.Register(&MigrationEntry{
    Version: "202504080001",        // 版本号格式：年月日序号
    Name:    "create_users_table",  // 迁移名称
    Up: func(tx *gorm.DB) error {   // 向上迁移函数
        return tx.AutoMigrate(&model.SysUser{})
    },
    Down: func(tx *gorm.DB) error { // 向下迁移函数
        return tx.Migrator().DropTable("sys_user")
    },
})
```

### 预定义迁移

系统已预定义了基础的迁移项：

1. 用户表 (`sys_user`)
2. 角色相关表 (`sys_role`, `sys_user_role`, `sys_role_menu`)
3. 权限相关表 (`sys_permission`, `sys_role_permission`)
4. 菜单表 (`sys_menu`)

### 添加新迁移

要添加新的迁移，在 `internal/migration/migrations.go` 文件中：

1. 创建新的注册函数或在已有函数中添加
2. 确保版本号遵循时间戳顺序
3. 在 `RegisterAllMigrations` 函数中注册

```go
// 示例：添加新的业务表迁移
func RegisterBusinessMigrations(migrator *Migrator) {
    migrator.Register(&MigrationEntry{
        Version: "202504080010",
        Name:    "create_products_table",
        Up: func(tx *gorm.DB) error {
            return tx.AutoMigrate(&model.Product{})
        },
        Down: func(tx *gorm.DB) error {
            return tx.Migrator().DropTable("products")
        },
    })
}

// 在RegisterAllMigrations中添加
func RegisterAllMigrations(migrator *Migrator) {
    // 已有迁移...
    RegisterBusinessMigrations(migrator)
}
```

### 迁移记录表

系统通过 `sys_migrations` 表跟踪迁移的执行状态，包含以下字段：

- `id`：自增主键
- `version`：迁移版本号（唯一索引）
- `name`：迁移名称
- `created_at`：执行时间
- `batch`：批次号（用于回滚）

## 环境配置

通过环境变量 `APP_ENV` 来指定运行环境，或通过 `--config` 标志直接指定配置文件：

- `APP_ENV=dev` 或 `APP_ENV=development` - 开发环境（默认）
- `APP_ENV=test` 或 `APP_ENV=testing` - 测试环境
- `APP_ENV=prod` 或 `APP_ENV=production` - 生产环境

示例：
```bash
# 开发环境运行服务器
APP_ENV=dev ./<app-name> server

# 生产环境执行数据库迁移
APP_ENV=prod ./<app-name> migrate
```

## 配置文件

配置文件根据运行环境自动加载对应的配置文件：

- `dev.yaml` - 开发环境配置
- `test.yaml` - 测试环境配置
- `prod.yaml` - 生产环境配置
- `example.yaml` - 示例配置（用于版本控制）

配置文件可以放置在以下位置（按查找顺序）：
1. 当前工作目录（项目根目录）
2. `configs/` 目录

首次使用时，请复制示例配置并根据环境重命名：

```bash
# 开发环境（放在根目录）
cp example.yaml ./dev.yaml

# 或放在configs目录
cp example.yaml configs/dev.yaml

# 生产环境
cp example.yaml configs/prod.yaml
```

应用程序会根据环境变量 `APP_ENV` 自动寻找并加载对应的配置文件。例如，当 `APP_ENV=dev` 时，将按以下顺序查找配置文件：
1. `./dev.yaml`（当前目录）
2. `./configs/dev.yaml`（configs目录）

如果找不到对应的配置文件，应用程序将无法启动。

## 日志配置

项目使用 [charmbracelet/log](https://github.com/charmbracelet/log) 作为日志库，支持控制台彩色输出和文件输出。

### 配置示例

```yaml
Log:
  Level: info                 # 日志级别: debug, info, warn, error
  Output: [console, file]     # 输出方式: console, file
  Format: text                # 日志格式: text, json
  FileConfig:
    Path: ./logs/app.log      # 日志文件路径
    MaxSize: 100              # 每个日志文件的最大大小（MB）
    MaxAge: 7                 # 日志文件保留天数
    MaxBackups: 10            # 保留的旧日志文件最大数量
    Compress: true            # 是否压缩旧日志文件
```

### 日志级别

- `debug`: 调试信息，包含详细的开发调试信息
- `info`: 一般信息，默认级别
- `warn`: 警告信息，需要注意的信息
- `error`: 错误信息，影响程序正常运行的错误

### 日志格式

- `text`: 普通文本格式，适合人类阅读（默认）
- `json`: JSON结构化格式，适合机器解析和日志系统收集

### 输出方式

- `console`: 输出到控制台，支持彩色输出
- `file`: 输出到文件，支持按大小自动分割、自动清理和压缩

可以同时配置多个输出方式，日志会同时输出到所有配置的目标。如果不配置 output，默认只输出到控制台。

### 文件输出配置

- `Path`: 日志文件路径
- `MaxSize`: 单个日志文件的最大大小（MB），超过后会自动分割
- `MaxAge`: 日志文件保留天数，超过后会自动删除
- `MaxBackups`: 保留的旧日志文件数量
- `Compress`: 是否压缩旧的日志文件

## 权限系统

项目集成了Casbin RBAC权限系统和动态菜单系统，实现了以下功能：

1. RBAC (基于角色的访问控制)权限模型
   - 用户 -> 角色 -> 权限
   - 支持资源级别和操作级别的权限控制

2. 动态菜单系统
   - 根据用户角色动态生成菜单
   - 菜单项与权限关联
   - 支持多级菜单树结构

3. 权限验证中间件
   - CasbinMiddleware：基于路径和HTTP方法的权限控制
   - PermissionMiddleware：基于菜单权限标识的权限控制

4. 数据表结构
   - sys_user - 用户表
   - sys_role - 角色表
   - sys_menu - 菜单表
   - sys_role_menu - 角色菜单关联表
   - sys_user_role - 用户角色关联表
   - casbin_rule - Casbin规则表(自动创建)

5. API接口
   - 菜单管理：创建、更新、删除、查询
   - 角色管理：创建、更新、删除、查询
   - 角色菜单分配
   - 角色权限分配
   - 用户角色分配

### 使用方法

1. 角色与菜单关联:
   ```
   POST /api/v1/admin/role/menu
   {
     "role_id": 1,
     "menu_ids": [1, 2, 3]
   }
   ```

2. 角色与权限关联:
   ```
   POST /api/v1/admin/role/permission
   {
     "role_code": "admin",
     "object": "/api/v1/admin/user",
     "action": "GET"
   }
   ```

3. 获取用户菜单:
   ```
   GET /api/v1/user/menus
   ```

4. 获取用户权限:
   ```
   GET /api/v1/user/perms
   ```

## 数据库查询选项系统

本项目实现了一个完整的数据库查询选项系统，用于简化GORM查询构建过程，提高代码复用性和可读性。

### 查询选项特点

- 采用函数式选项模式设计
- 支持链式组合多个查询条件
- 提供统一的接口方式处理各种查询场景
- 易于扩展和自定义新的查询条件

### 基本使用方法

```go
// 导入查询选项包
import "your-project/internal/pkg/options"

// 创建查询实例
query := options.Apply(
    DB, // *gorm.DB实例
    options.WithPage(1, 10),
    options.WithOrder("created_at", "desc"),
    options.WithLike("name", keyword),
)

// 执行查询
var results []YourModel
query.Find(&results)
```

### 内置查询选项

系统提供了以下内置查询选项：

#### 分页与排序
- `WithPage(page, pageSize)` - 分页查询，自动限制最大页面大小
- `WithOrder(field, direction)` - 排序查询，direction支持"asc"或"desc"

#### 关联查询
- `WithPreload(relation, args...)` - 预加载关联关系
- `WithJoin(query, args...)` - 连接查询
- `WithSelect(query, args...)` - 指定查询字段
- `WithGroup(query)` - 分组查询
- `WithHaving(query, args...)` - HAVING条件查询

#### 条件过滤
- `WithWhere(query, args...)` - WHERE条件
- `WithOrWhere(query, args...)` - OR WHERE条件
- `WithLike(field, value)` - LIKE模糊查询
- `WithExactMatch(field, value)` - 精确匹配查询
- `WithTimeRange(field, start, end)` - 时间范围查询
- `WithKeyword(keyword, fields...)` - 关键字搜索（多字段OR条件）

#### 组合查询
- `WithBaseQuery(tableName, status, keyword, keywordFields, createBy, startTime, endTime)` - 应用基础查询条件，组合多个常用过滤条件

### 自定义查询选项

可以轻松扩展自定义的查询选项：

```go
// 自定义查询选项示例
func WithCustomCondition(param string) options.Option {
    return func(db *gorm.DB) *gorm.DB {
        if param == "" {
            return db
        }
        return db.Where("custom_field = ?", param)
    }
}

// 使用自定义查询选项
query := options.Apply(
    DB,
    options.WithPage(1, 10),
    WithCustomCondition("value"),
)
```

### 与DTO结合使用

可以结合DTO对象灵活构建查询条件：

```go
// 基于BaseQuery 构建查询条件
func BuildQueryOptions(q *request.BaseQuery, tableName string) []options.Option {
    var opts []options.Option

    // 添加基础查询条件
    opts = append(opts, options.WithBaseQuery(
        tableName,
        q.Status,
        q.Keyword,
        []string{"name", "description"}, // 关键字搜索字段
        q.CreateBy,
        q.StartTime,
        q.EndTime,
    ))

    return opts
}

// 在服务中使用
func (s *Service) List(query *request.YourQuery) ([]YourModel, int64, error) {
    opts := BuildQueryOptions(&query.BaseQuery, "your_table")

    // 添加分页和排序
    opts = append(opts,
        options.WithPage(query.Page, query.PageSize),
        options.WithOrder(query.SortField, query.SortOrder),
    )

    // 应用所有查询选项
    db := options.Apply(s.DB, opts...)

    var total int64
    db.Model(&YourModel{}).Count(&total)

    var items []YourModel
    db.Find(&items)

    return items, total, nil
}
```

## 组件访问规范

本项目使用组件模式管理各种资源，推荐以下访问方式：

```go
// 获取数据库连接
db := sqldb.Instance().DB()

// 获取Redis客户端
client := redisdb.Instance().Client()

// 获取MongoDB数据库
mongo := mongodb.Instance().DB()
```

