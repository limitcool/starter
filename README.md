# Starter
## 特征
- 提供 gin 框架项目模版
- 集成 GORM 进行 ORM 映射和数据库操作
- 集成 Viper 进行配置管理
- 提供常用 gin 中间件和工具
  - 跨域中间件:处理 API 跨域请求,实现 CORS 支持
  - jwt 解析中间件:从请求中解析并验证 JWT Token,用于 API 身份认证

## Quick Start

```bash
go install github.com/go-eagle/eagle/cmd/eagle@latest
eagle new <project name> -r https://github.com/limitcool/starter
```