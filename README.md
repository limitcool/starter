# Starter (Lite Version)

[![Go Reference](https://pkg.go.dev/badge/github.com/limitcool/starter.svg)](https://pkg.go.dev/github.com/limitcool/starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/limitcool/starter)](https://goreportcard.com/report/github.com/limitcool/starter)

[English](README_EN.md) | 中文

> 这是 Starter 框架的轻量级版本，采用单用户模式设计，适合快速开发和简单应用场景。如果您需要更多企业级功能，请查看 `enterprise` 分支。

## 特征
- 提供基于 Gin 框架的轻量级项目模板
- 使用简洁的手动依赖注入，实现清晰的代码结构
- 采用简化的架构设计，专注于快速开发
- 集成 GORM 进行 ORM 映射和数据库操作
  - 支持 PostgreSQL (使用 pgx 驱动)
  - 支持 MySQL
  - 支持 SQLite
- 集成 Viper 进行配置管理
- 提供常用 Gin 中间件和工具
  - 跨域中间件：处理 API 跨域请求，实现 CORS 支持
  - JWT 解析中间件：从请求中解析并验证 JWT Token，用于 API 身份认证
- 国际化 (i18n) 支持
  - 基于请求 Accept-Language 头自动选择语言
  - 错误消息多语言支持
- 使用 Cobra 命令行框架，提供清晰的子命令结构
- 支持数据库迁移与服务器启动分离，提高启动速度
- 简化的用户管理系统，使用单一用户表和 IsAdmin 字段区分管理员
- 优化的错误处理系统，支持错误码和多语言错误消息

## 架构设计

项目采用简化的架构设计，使用手动依赖注入，实现了清晰的代码结构：

### 1. 简化的分层架构

- **Model 层**：定义数据模型和数据库表结构
- **Handler 层**：处理 HTTP 请求和响应，直接与数据库交互
- **Router 层**：定义 API 路由，依赖于 Handler 层

### 2. 依赖注入

项目使用简洁的手动依赖注入，通过构造函数注入依赖：

```go
// Model 层
func NewUserRepo(db *gorm.DB) *UserRepo {
    // ...
}

// Handler 层
func NewUserHandler(db *gorm.DB, config *configs.Config) *UserHandler {
    // ...
}

// Router 层
func NewRouter(userHandler *handler.UserHandler) *gin.Engine {
    // ...
}
```

### 3. 应用容器管理

使用应用容器管理组件的生命周期，确保组件的正确初始化和清理：

```go
type App struct {
    config   *configs.Config
    db       *gorm.DB
    handlers *Handlers
    router   *gin.Engine
    server   *http.Server
}

func New(config *configs.Config) (*App, error) {
    app := &App{config: config}

    // 按顺序初始化各个组件
    if err := app.initDatabase(); err != nil {
        return nil, err
    }
    // ... 其他组件初始化

    return app, nil
}
```

## 快速开始

```bash
go install github.com/go-eagle/eagle/cmd/eagle@latest
eagle new <project name> -r https://github.com/limitcool/starter -b lite
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

## 用户管理设计

Lite 版本采用简化的用户管理设计，专注于快速开发和简单应用场景：

### 单一用户表设计

- 只使用 `user` 一张表存储所有用户信息
- 通过 `is_admin` 字段区分普通用户和管理员用户
- 简化的数据库结构，减少表之间的关联
- 使用简化的权限检查中间件，只检查用户是否为管理员

### 配置管理员用户

在配置文件中设置初始管理员用户：

```yaml
admin:
  username: "admin"
  password: "admin123"
  nickname: "超级管理员"
```

系统会在首次启动时自动创建管理员用户。

## 数据库迁移系统

Lite 版本实现了一个简洁的数据库迁移系统，用于管理数据库表结构的创建和更新。

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
        return tx.AutoMigrate(&model.User{})
    },
    Down: func(tx *gorm.DB) error { // 向下迁移函数
        return tx.Migrator().DropTable("users")
    },
})
```

### 预定义迁移

Lite 版本已预定义了基础的迁移项：

1. 用户表 (`users`)
2. 文件表 (`files`)

### 添加新迁移

要添加新的迁移，在 `internal/migration/migrations.go` 文件中：

1. 创建新的注册函数或在已有函数中添加
2. 确保版本号遵循时间戳顺序
3. 使用 `RegisterMigration` 函数注册

```go
// 示例：添加新的业务表迁移
RegisterMigration("create_products_table",
    // 向上迁移函数
    func(tx *gorm.DB) error {
        return tx.AutoMigrate(&model.Product{})
    },
    // 向下迁移函数
    func(tx *gorm.DB) error {
        return tx.Migrator().DropTable("products")
    },
)
```

### 迁移记录表

系统通过 `migrations` 表跟踪迁移的执行状态，包含以下字段：

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
# 开发环境
cp example.yaml dev.yaml

# 测试环境
cp example.yaml test.yaml

# 生产环境
cp example.yaml prod.yaml
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

## 性能分析 (Pprof) 支持

项目内置了 Go 官方的 pprof 性能分析工具，可以帮助开发者分析应用的性能瓶颈。

### 配置示例

```yaml
Pprof:
  Enabled: true    # 是否启用pprof
  Port: 0          # pprof服务端口，0表示使用主服务端口
```

### 配置说明

- `Enabled`: 是否启用pprof功能，生产环境建议设为false
- `Port`: pprof服务端口
  - `0`: 在主HTTP服务器上启用pprof路由 (推荐开发环境)
  - `6060`: 启动独立的pprof服务器 (推荐生产环境调试)

### 使用方法

#### 1. 主服务器模式 (Port: 0)

当配置 `Port: 0` 时，pprof路由会添加到主HTTP服务器：

```bash
# 启动应用
./starter server

# 访问pprof主页
curl http://localhost:8080/debug/pprof/

# 查看goroutine信息
curl http://localhost:8080/debug/pprof/goroutine?debug=1

# 获取CPU profile (30秒)
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof

# 查看内存heap信息
curl http://localhost:8080/debug/pprof/heap?debug=1
```

#### 2. 独立服务器模式 (Port: 6060)

当配置具体端口时，会启动独立的pprof服务器：

```yaml
Pprof:
  Enabled: true
  Port: 6060
```

```bash
# 启动应用
./starter server

# 访问独立pprof服务器
curl http://localhost:6060/debug/pprof/

# 使用go tool pprof分析
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
go tool pprof http://localhost:6060/debug/pprof/heap
```

### 常用pprof端点

- `/debug/pprof/` - pprof主页，显示所有可用的profile
- `/debug/pprof/goroutine` - 查看所有goroutine的堆栈信息
- `/debug/pprof/heap` - 查看内存分配信息
- `/debug/pprof/profile` - CPU性能分析 (默认30秒)
- `/debug/pprof/block` - 查看阻塞操作的堆栈信息
- `/debug/pprof/mutex` - 查看互斥锁竞争信息
- `/debug/pprof/allocs` - 查看内存分配采样信息
- `/debug/pprof/threadcreate` - 查看线程创建信息

### 安全建议

- **开发环境**: 可以使用主服务器模式 (`Port: 0`)，方便调试
- **生产环境**: 建议禁用pprof (`Enabled: false`) 或使用独立端口并限制访问
- **调试生产问题**: 临时启用独立端口模式，调试完成后立即禁用

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

   在应用容器中注册 gRPC 控制器：

   ```go
   // 在应用容器中初始化 gRPC 服务
   func (a *App) initGRPCServer() error {
       // 创建 gRPC 服务器
       grpcServer := grpc.NewServer()

       // 注册服务
       systemController := NewSystemGRPCController(a.config)
       pb.RegisterSystemServiceServer(grpcServer, systemController)

       a.grpcServer = grpcServer
       return nil
   }
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

Lite 版本采用简化的权限系统，基于用户的 `is_admin` 字段进行权限控制：

### 权限控制中间件

项目提供了三种权限控制中间件：

1. **AdminCheck**：检查用户是否为管理员
   - 基于 JWT 中的 `is_admin` 字段进行快速检查
   - 适用于管理员专属接口

2. **UserCheck**：检查用户是否已登录
   - 只验证用户是否登录，不检查用户类型
   - 适用于需要登录但不限制用户类型的接口

3. **RegularUserCheck**：检查用户是否为普通用户
   - 确保用户不是管理员
   - 适用于只允许普通用户访问的接口

### 使用方法

在路由定义中使用中间件：

```go
// 管理员接口
adminGroup := router.Group("/api/v1/admin")
adminGroup.Use(middleware.AdminCheck())
{
    adminGroup.GET("/users", handler.ListUsers)
    // 其他管理员接口...
}

// 普通用户接口
userGroup := router.Group("/api/v1/user")
userGroup.Use(middleware.UserCheck())
{
    userGroup.GET("/profile", handler.GetUserProfile)
    // 其他用户接口...
}

// 仅普通用户接口
regularUserGroup := router.Group("/api/v1/regular")
regularUserGroup.Use(middleware.RegularUserCheck())
{
    regularUserGroup.POST("/feedback", handler.SubmitFeedback)
    // 其他仅普通用户接口...
}
```

## 数据库操作

Lite 版本采用简化的数据库操作方式，直接在 Model 层提供数据库操作方法：

### Model 层设计

在 Lite 版本中，Model 层直接提供数据库操作方法，简化了代码结构：

```go
// User 用户模型
type User struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    Username  string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
    Password  string    `gorm:"size:100;not null" json:"-"`
    Nickname  string    `gorm:"size:50" json:"nickname"`
    Email     string    `gorm:"size:100" json:"email"`
    Avatar    string    `gorm:"size:255" json:"avatar"`
    IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// UserRepo 用户数据库操作
type UserRepo struct {
    DB *gorm.DB
}

// NewUserRepo 创建用户仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
    return &UserRepo{DB: db}
}

// GetByID 根据ID获取用户
func (r *UserRepo) GetByID(ctx context.Context, id uint) (*User, error) {
    var user User
    if err := r.DB.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errorx.ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
    var user User
    if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errorx.ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

// Create 创建用户
func (r *UserRepo) Create(ctx context.Context, user *User) error {
    return r.DB.Create(user).Error
}

// Update 更新用户
func (r *UserRepo) Update(ctx context.Context, user *User) error {
    return r.DB.Save(user).Error
}

// Delete 删除用户
func (r *UserRepo) Delete(ctx context.Context, id uint) error {
    return r.DB.Delete(&User{}, id).Error
}

// List 获取用户列表
func (r *UserRepo) List(ctx context.Context, page, pageSize int) ([]User, int64, error) {
    var users []User
    var total int64

    r.DB.Model(&User{}).Count(&total)

    offset := (page - 1) * pageSize
    if err := r.DB.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
}
```

### 在 Handler 中使用

在 Handler 层直接使用 Model 层提供的方法：

```go
// UserHandler 用户处理器
type UserHandler struct {
    userRepo *model.UserRepo
    config   *configs.Config
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userRepo *model.UserRepo, config *configs.Config) *UserHandler {
    return &UserHandler{
        userRepo: userRepo,
        config:   config,
    }
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
    id := cast.ToUint(c.Param("id"))

    user, err := h.userRepo.GetByID(c.Request.Context(), id)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, user)
}
```

## 错误处理

Lite 版本提供了简洁而强大的错误处理系统，支持错误码和多语言错误消息。

### 错误处理特点

- 统一的错误码定义和管理
- 错误包装，保留完整的错误链
- 多语言错误消息支持
- 区分内部错误和用户可见错误

### 错误定义

错误定义在 `internal/pkg/errorx` 包中：

```go
// 错误码定义
const (
    // 通用错误码
    ErrorSuccess       = 0    // 成功
    ErrorUnknown       = 1000 // 未知错误
    ErrorInvalidParams = 1001 // 无效参数
    ErrorNotFound      = 1002 // 资源不存在
    ErrorDatabase      = 1003 // 数据库错误

    // 用户相关错误码
    ErrorUserNotFound     = 2000 // 用户不存在
    ErrorUserAlreadyExist = 2001 // 用户已存在
    ErrorUserAuthFailed   = 2002 // 用户认证失败
    ErrorUserNoLogin      = 2003 // 用户未登录
    ErrorAccessDenied     = 2004 // 访问被拒绝
)

// Error 自定义错误类型
type Error struct {
    Code    int    // 错误码
    Message string // 错误消息
    Err     error  // 原始错误
}

// NewError 创建新错误
func NewError(code int, msg string) *Error {
    return &Error{
        Code:    code,
        Message: msg,
    }
}

// WithMsg 设置错误消息
func (e *Error) WithMsg(msg string) *Error {
    return &Error{
        Code:    e.Code,
        Message: msg,
        Err:     e.Err,
    }
}

// WithError 包装原始错误
func (e *Error) WithError(err error) *Error {
    return &Error{
        Code:    e.Code,
        Message: e.Message,
        Err:     err,
    }
}
```

### 使用示例

```go
// 在 Model 层
func (r *UserRepo) GetByID(ctx context.Context, id uint) (*User, error) {
    var user User
    if err := r.DB.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errorx.NewError(errorx.ErrorUserNotFound, "用户不存在")
        }
        return nil, errorx.NewError(errorx.ErrorDatabase, "数据库错误").WithError(err)
    }
    return &user, nil
}

// 在 Handler 层
func (h *UserHandler) GetUser(c *gin.Context) {
    id := cast.ToUint(c.Param("id"))

    user, err := h.userRepo.GetByID(c.Request.Context(), id)
    if err != nil {
        logger.Error("获取用户失败", "error", err, "id", id)
        response.Error(c, err)
        return
    }

    response.Success(c, user)
}
```

### 统一响应格式

所有 API 响应使用统一的格式：

```go
// 成功响应
{
    "code": 0,
    "message": "success",
    "data": {
        // 响应数据
    }
}

// 错误响应
{
    "code": 1001,
    "message": "无效参数",
    "data": null
}
```
