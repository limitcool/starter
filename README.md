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