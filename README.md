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

配置文件位于 `configs` 目录：

- `config.yaml` - 默认配置
- `config-dev.yaml` - 开发环境配置
- `config-test.yaml` - 测试环境配置
- `config-prod.yaml` - 生产环境配置
- `config.example.yaml` - 示例配置（用于版本控制）

首次使用时，请复制 `config.example.yaml` 并根据环境重命名：

```bash
# 开发环境
cp configs/config.example.yaml configs/config-dev.yaml

# 生产环境
cp configs/config.example.yaml configs/config-prod.yaml
```

配置文件加载顺序：
1. 加载 `config.yaml` 作为默认配置
2. 根据 `APP_ENV` 加载对应的环境配置文件，覆盖默认配置

## 日志配置

项目使用 [charmbracelet/log](https://github.com/charmbracelet/log) 作为日志库，支持控制台彩色输出和文件输出。

### 配置示例

```yaml
log:
  level: info                 # 日志级别: debug, info, warn, error
  output: [console, file]     # 输出方式: console, file
  file:
    path: ./logs/app.log      # 日志文件路径
    maxSize: 100              # 每个日志文件的最大大小（MB）
    maxAge: 7                 # 日志文件保留天数
    maxBackups: 10            # 保留的旧日志文件最大数量
    compress: true            # 是否压缩旧日志文件
```

### 日志级别

- `debug`: 调试信息，包含详细的开发调试信息
- `info`: 一般信息，默认级别
- `warn`: 警告信息，需要注意的信息
- `error`: 错误信息，影响程序正常运行的错误

### 输出方式

- `console`: 输出到控制台，支持彩色输出
- `file`: 输出到文件，支持按大小自动分割、自动清理和压缩

可以同时配置多个输出方式，日志会同时输出到所有配置的目标。如果不配置 output，默认只输出到控制台。

### 文件输出配置

- `path`: 日志文件路径
- `maxSize`: 单个日志文件的最大大小（MB），超过后会自动分割
- `maxAge`: 日志文件保留天数，超过后会自动删除
- `maxBackups`: 保留的旧日志文件数量
- `compress`: 是否压缩旧的日志文件