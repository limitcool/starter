# Starter

[![Go Reference](https://pkg.go.dev/badge/github.com/limitcool/starter.svg)](https://pkg.go.dev/github.com/limitcool/starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/limitcool/starter)](https://goreportcard.com/report/github.com/limitcool/starter)

## 特征
- 提供 gin 框架项目模版
- 集成 GORM 进行 ORM 映射和数据库操作
  - 支持 PostgreSQL (使用 pgx 驱动)
- 集成 Viper 进行配置管理
- 提供常用 gin 中间件和工具
  - 跨域中间件:处理 API 跨域请求,实现 CORS 支持
  - jwt 解析中间件:从请求中解析并验证 JWT Token,用于 API 身份认证

## 快速开始

```bash
go install github.com/go-eagle/eagle/cmd/eagle@latest
eagle new <project name> -r https://github.com/limitcool/starter -b main
```

## 环境配置

通过环境变量 `APP_ENV` 来指定运行环境：

- `APP_ENV=dev` 或 `APP_ENV=development` - 开发环境（默认）
- `APP_ENV=test` 或 `APP_ENV=testing` - 测试环境
- `APP_ENV=prod` 或 `APP_ENV=production` - 生产环境

示例：
```bash
# 开发环境运行
APP_ENV=dev go run main.go

# 生产环境运行
APP_ENV=prod ./starter
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