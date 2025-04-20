# Starter

[![Go Reference](https://pkg.go.dev/badge/github.com/limitcool/starter.svg)](https://pkg.go.dev/github.com/limitcool/starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/limitcool/starter)](https://goreportcard.com/report/github.com/limitcool/starter)

[English](README_EN.md) | 中文

## 特征
- 提供 gin 框架项目模版
- 支持 HTTP 和 gRPC 双协议服务
  - 可通过配置启用/禁用 gRPC 服务
  - 统一的 API 定义和实现
  - 支持 gRPC 健康检查和反射服务
- 使用 Uber fx 框架进行依赖注入，实现更清晰的代码结构
- 采用标准 MVC 架构，遵循关注点分离原则
- 集成 GORM 进行 ORM 映射和数据库操作
  - 支持 PostgreSQL (使用 pgx 驱动)
  - 支持 MySQL
  - 支持 SQLite
  - 提供丰富的查询选项工具函数
- 集成 Viper 进行配置管理
- 提供常用 gin 中间件和工具
  - 跨域中间件:处理 API 跨域请求,实现 CORS 支持
  - jwt 解析中间件:从请求中解析并验证 JWT Token,用于 API 身份认证
- 国际化 (i18n) 支持
  - 基于请求 Accept-Language 头自动选择语言
  - 错误消息多语言支持
  - 内置英语 (en-US) 和中文 (zh-CN) 翻译
  - 可轻松扩展支持更多语言
- 使用 Cobra 命令行框架，提供清晰的子命令结构
- 支持数据库迁移与服务器启动分离，提高启动速度
- 完善的数据库迁移系统，支持版本控制和回滚
- 内置用户、角色、权限和菜单管理系统
- 优化的错误处理系统，支持错误码和多语言错误消息

## 架构设计

项目采用标准的 MVC 架构，并结合 Uber fx 依赖注入框架，实现了清晰的层次结构：

### 1. 分层架构

- **Model 层**：定义数据模型和数据库表结构
- **Repository 层**：负责数据访问，是唯一直接与数据库交互的层
- **Service 层**：实现业务逻辑，依赖于 Repository 层
- **Controller 层**：处理 HTTP 请求和响应，依赖于 Service 层
- **Router 层**：定义 API 路由，依赖于 Controller 层

### 2. 依赖注入

项目使用 Uber fx 框架实现依赖注入，每一层都通过构造函数注入其依赖：

```go
// Repository 层
func NewUserRepo(db *gorm.DB) *UserRepo {
    // ...
}

// Service 层
func NewUserService(userRepo *repository.UserRepo) *UserService {
    // ...
}

// Controller 层
func NewUserController(userService *services.UserService) *UserController {
    // ...
}

// Router 层
func NewRouter(userController *controller.UserController) *gin.Engine {
    // ...
}
```

### 3. 生命周期管理

使用 fx.Lifecycle 管理组件的生命周期，确保组件的正确初始化和清理：

```go
func NewComponent(lc fx.Lifecycle) *Component {
    component := &Component{}

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            // 初始化逻辑
            return nil
        },
        OnStop: func(ctx context.Context) error {
            // 清理逻辑
            return nil
        },
    })

    return component
}
```

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

## 国际化 (i18n) 支持

系统内置了国际化支持，可以根据客户端请求自动切换语言。

### 配置国际化

在配置文件中设置国际化选项：

```yaml
I18n:
  Enabled: true                # 是否启用国际化
  DefaultLanguage: en-US       # 默认语言
  SupportLanguages:            # 支持的语言列表
    - zh-CN
    - en-US
  ResourcesPath: locales       # 语言资源文件路径
```

### 语言资源文件

语言资源文件位于 `locales` 目录下，采用 JSON 格式：

- `locales/en-US.json` - 英文资源
- `locales/zh-CN.json` - 中文资源

示例语言文件内容：

```json
{
  "error.success": "Success",
  "error.common.invalid_params": "Invalid request parameters",
  "error.user.user_not_found": "User not found"
}
```

### 使用方法

1. **API响应自动翻译**：
   - 系统会自动根据请求头 `Accept-Language` 选择语言
   - API错误响应会根据设置的语言返回对应的翻译文本

2. **客户端请求示例**：
   ```bash
   # 请求英文响应
   curl -X POST "http://localhost:8080/api/v1/user/login" \
        -H "Accept-Language: en-US" \
        -H "Content-Type: application/json" \
        -d '{"username": "test", "password": "wrong"}'

   # 请求中文响应
   curl -X POST "http://localhost:8080/api/v1/user/login" \
        -H "Accept-Language: zh-CN" \
        -H "Content-Type: application/json" \
        -d '{"username": "test", "password": "wrong"}'
   ```

3. **添加新的错误码翻译**：
   - 在 `tools/errorgen/error_codes.md` 中定义错误
   - 运行错误代码生成器: `go run tools/errorgen/main.go tools/errorgen/error_codes.md internal/pkg/errorx/code_gen.go`
   - 在语言文件 (`locales/en-US.json` 和 `locales/zh-CN.json`) 中添加对应的翻译

4. **添加新的语言支持**：
   - 创建新的语言文件，如 `locales/fr-FR.json`
   - 在配置中的 `SupportLanguages` 列表中添加该语言
   - 重启应用使配置生效

## 错误处理系统

项目实现了一个完整的错误处理系统，包括错误码、错误包装和多语言错误消息。

### 错误处理特点

- 统一的错误码定义和管理
- 错误包装，保留完整的错误链和堆栈信息
- 多语言错误消息支持
- 区分内部错误和用户可见错误

### 错误处理最佳实践

- Repository 层：返回具体错误，不记录日志
- Service 层：包装错误，添加业务上下文，不记录日志
- Controller 层：转换为用户友好的错误响应，记录完整错误日志

### 使用示例

```go
// Repository 层
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
    var user model.User
    if err := r.DB.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errorx.NewError(errorx.ErrorUserNotFoundCode, "用户不存在")
        }
        return nil, err
    }
    return &user, nil
}

// Service 层
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, errorx.WrapError(err, fmt.Sprintf("获取用户ID %d 失败", id))
    }
    return user, nil
}

// Controller 层
func (c *UserController) GetUser(ctx *gin.Context) {
    id := cast.ToInt64(ctx.Param("id"))
    user, err := c.userService.GetUserByID(ctx, id)
    if err != nil {
        logger.Error("获取用户失败", "error", err, "id", id)
        response.Error(ctx, err)
        return
    }
    response.Success(ctx, user)
}
```

## 日志配置

项目使用 [uber-go/zap](https://github.com/uber-go/zap) 作为日志库，提供高性能的结构化日志记录。

### 配置示例

```yaml
Log:
  Level: info                 # 日志级别: debug, info, warn, error
  Output: [console, file]     # 输出方式: console, file
  Format: text                # 日志格式: text, json
  StackTraceEnabled: true     # 是否启用堆栈跟踪
  StackTraceLevel: error      # 堆栈跟踪级别: error, warn, info, debug
  MaxStackFrames: 20          # 最大堆栈帧数
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

### 堆栈跟踪配置

- `StackTraceEnabled`: 是否启用堆栈跟踪，默认为 true
- `StackTraceLevel`: 堆栈跟踪级别，默认为 error，只有该级别及以上的日志才会记录堆栈
- `MaxStackFrames`: 最大堆栈帧数，默认为 20

### 文件输出配置

- `Path`: 日志文件路径
- `MaxSize`: 单个日志文件的最大大小（MB），超过后会自动分割
- `MaxAge`: 日志文件保留天数，超过后会自动删除
- `MaxBackups`: 保留的旧日志文件数量
- `Compress`: 是否压缩旧的日志文件

### 使用示例

日志库提供了统一的接口：

```go
// 导入日志包
import "github.com/limitcool/starter/internal/pkg/logger"

// 记录不同级别的日志
func example() {
    // 记录信息日志
    logger.Info("This is an info message", "user", "admin", "action", "login")

    // 记录警告日志
    logger.Warn("This is a warning message", "memory", "90%")

    // 记录错误日志
    err := someFunction()
    if err != nil {
        logger.Error("Operation failed", "error", err, "operation", "someFunction")
    }

    // 记录错误日志，并包含堆栈信息
    if err != nil {
        logger.LogErrorWithStack("Operation failed with stack", err, "operation", "someFunction")
    }

    // 记录调试日志
    logger.Debug("Detailed debug information", "request", req, "response", resp)
}
```

## gRPC 支持

项目集成了 gRPC 支持，可与 HTTP 服务并行运行，提供高性能的 RPC 服务。

### gRPC 特点

- 可通过配置启用/禁用 gRPC 服务
- 支持 gRPC 健康检查和反射服务
- 与 HTTP 服务共享业务逻辑
- 使用 Protocol Buffers 定义 API

### 配置示例

```yaml
GRPC:
  Enabled: true                # 是否启用 gRPC 服务
  Port: 9000                  # gRPC 服务端口
  HealthCheck: true           # 是否启用健康检查服务
  Reflection: true            # 是否启用反射服务
```

### 使用方法

1. **定义 Proto 文件**

   Proto 文件定义在 `internal/proto/v1` 目录下，生成的代码在 `internal/proto/gen/v1` 目录下。

   ```protobuf
   // internal/proto/v1/system.proto
   syntax = "proto3";

   package internal.proto.v1;

   option go_package = "internal/proto/gen/v1;protov1";

   // SystemService 系统服务
   service SystemService {
     // GetSystemInfo 获取系统信息
     rpc GetSystemInfo(SystemInfoRequest) returns (SystemInfoResponse) {}
   }

   // SystemInfoRequest 系统信息请求
   message SystemInfoRequest {
     // 请求ID
     string request_id = 1;
   }

   // SystemInfoResponse 系统信息响应
   message SystemInfoResponse {
     // 应用名称
     string app_name = 1;
     // 应用版本
     string version = 2;
     // 运行模式
     string mode = 3;
     // 服务器时间
     int64 server_time = 4;
   }
   ```

2. **生成 gRPC 代码**

   使用 Makefile 中的 proto 命令生成 gRPC 代码：

   ```bash
   make proto
   ```

3. **实现 gRPC 控制器**

   在 `internal/controller` 目录下创建 gRPC 控制器，使用 `_grpc` 后缀区分：

   ```go
   // internal/controller/system_grpc.go
   package controller

   import (
       "context"
       "time"

       pb "github.com/limitcool/starter/internal/proto/gen/v1"
       // ...
   )

   // SystemGRPCController gRPC系统控制器
   type SystemGRPCController struct {
       pb.UnimplementedSystemServiceServer
       // ...
   }

   // GetSystemInfo 获取系统信息
   func (c *SystemGRPCController) GetSystemInfo(ctx context.Context, req *pb.SystemInfoRequest) (*pb.SystemInfoResponse, error) {
       // 实现业务逻辑
       return &pb.SystemInfoResponse{
           AppName:    c.config.App.Name,
           Version:    "1.0.0",
           Mode:       c.config.App.Mode,
           ServerTime: time.Now().Unix(),
       }, nil
   }
   ```

4. **注册 gRPC 服务**

   在 `internal/controller/module.go` 中注册 gRPC 控制器：

   ```go
   // Module 控制器模块
   var Module = fx.Options(
       // 提供 HTTP 控制器
       fx.Provide(NewUserController),
       // ...

       // 提供 gRPC 控制器
       fx.Provide(NewSystemGRPCController),

       // 注册 gRPC 控制器
       fx.Invoke(RegisterSystemGRPCController),
   )
   ```

5. **使用 gRPC 客户端**

   ```go
   // 创建 gRPC 连接
   conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
   if err != nil {
       log.Fatalf("Failed to connect: %v", err)
   }
   defer conn.Close()

   // 创建客户端
   client := pb.NewSystemServiceClient(conn)

   // 调用服务
   resp, err := client.GetSystemInfo(context.Background(), &pb.SystemInfoRequest{
       RequestId: "test-request",
   })
   ```

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

## 泛型仓库系统

项目利用 Go 1.18+ 的泛型特性，实现了完整的泛型仓库系统，显著减少了重复代码，提高了开发效率。

### 泛型仓库接口

泛型仓库接口定义了所有仓库实现必须提供的方法：

```go
// Repository 通用仓库接口
type Repository[T Entity] interface {
	// 基本 CRUD 操作
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id any) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id any) error
	List(ctx context.Context, page, pageSize int) ([]T, int64, error)
	Count(ctx context.Context) (int64, error)
	UpdateFields(ctx context.Context, id any, fields map[string]any) error

	// 批量操作
	BatchCreate(ctx context.Context, entities []T) error
	BatchDelete(ctx context.Context, ids []any) error

	// 查询方法
	FindByField(ctx context.Context, field string, value any) (*T, error)
	FindAllByField(ctx context.Context, field string, value any) ([]T, error)
	FindByCondition(ctx context.Context, condition string, args ...any) ([]T, error)
	FindOneByCondition(ctx context.Context, condition string, args ...any) (*T, error)
	GetPage(ctx context.Context, page, pageSize int, condition string, args ...any) ([]T, int64, error)

	// 高级查询
	FindWithLike(ctx context.Context, field string, value string) ([]T, error)
	FindWithIn(ctx context.Context, field string, values []any) ([]T, error)
	FindWithBetween(ctx context.Context, field string, min, max any) ([]T, error)
	CountWithCondition(ctx context.Context, condition string, args ...any) (int64, error)
	AggregateField(ctx context.Context, aggregate Aggregate, field string, condition string, args ...any) (float64, error)
	GroupBy(ctx context.Context, groupFields []string, selectFields []string, condition string, args ...any) ([]map[string]any, error)
	Join(ctx context.Context, joinType string, table string, on string, selectFields []string, condition string, args ...any) ([]map[string]any, error)
	Exists(ctx context.Context, condition string, args ...any) (bool, error)
	Raw(ctx context.Context, sql string, values ...any) ([]map[string]any, error)

	// 事务相关
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	WithTx(tx *gorm.DB) Repository[T]
}
```

### 泛型仓库实现

项目提供了两种泛型仓库实现：

1. **GenericRepo**：基本泛型仓库实现，直接与数据库交互
2. **CachedRepo**：带缓存的泛型仓库实现，封装了基本仓库并添加缓存功能

#### 创建泛型仓库

```go
// 创建基本泛型仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
    // 创建泛型仓库
    genericRepo := NewGenericRepo[model.User](db)
    // 设置错误码
    genericRepo.SetErrorCode(errorx.ErrorUserNotFoundCode)

    return &UserRepo{
        DB:          db,
        GenericRepo: genericRepo,
    }
}

// 创建带缓存的泛型仓库
func NewCachedUserRepo(repo Repository[model.User]) (Repository[model.User], error) {
    return WithCache(repo, "user", "user", 30*time.Minute)
}
```

#### 使用泛型仓库

```go
// 在服务层使用泛型仓库
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
    // 直接调用泛型仓库方法
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// 使用高级查询功能
func (s *UserService) SearchUsers(ctx context.Context, keyword string, page, pageSize int) ([]model.User, int64, error) {
    condition := ""
    var args []any

    if keyword != "" {
        condition = "username LIKE ? OR email LIKE ?"
        args = append(args, "%"+keyword+"%", "%"+keyword+"%")
    }

    return s.userRepo.GetPage(ctx, page, pageSize, condition, args...)
}
```

### 泛型仓库的优势

1. **代码复用**：所有实体共享相同的仓库实现，显著减少重复代码
2. **类型安全**：利用Go泛型特性确保类型安全，编译期就能发现类型错误
3. **功能丰富**：提供了完整的CRUD操作和高级查询功能
4. **缓存支持**：通过装饰器模式轻松添加缓存功能，提高性能
5. **事务支持**：内置事务支持，确保数据一致性

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
