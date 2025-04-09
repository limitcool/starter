# 应用错误码定义

本文档定义了应用中使用的所有错误码，包括错误码值、错误消息和HTTP状态码。

## 基础错误码

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 0 | SuccessCode | 成功 | 200 |

## 通用错误码 (10000-19999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 10000 | InvalidParamsCode | 请求参数错误 | 400 |
| 10001 | ErrorUnknownCode | 服务器开小差啦，稍后再来试一试 | 500 |
| 10002 | ErrorNotExistCertCode | 不存在的认证类型 | 400 |
| 10003 | ErrorNotFoundCode | 资源不存在 | 404 |
| 10004 | ErrorDatabaseCode | 数据库操作失败 | 500 |
| 10005 | ErrorInternalCode | 服务器内部错误 | 500 |
| 10006 | ErrorCode | 错误 | 500 |
| 10007 | ErrorParamCode | 参数错误 | 400 |

## 数据库错误码 (20000-29999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 20000 | DatabaseInsertErrorCode | 数据库插入失败 | 500 |
| 20001 | DatabaseDeleteErrorCode | 数据库删除失败 | 500 |
| 20002 | DatabaseQueryErrorCode | 数据库查询失败 | 500 |

## 用户相关错误码 (30000-39999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 30000 | UserNoLoginCode | 用户未登录 | 401 |
| 30001 | UserNotFoundCode | 用户不存在 | 404 |
| 30002 | UserPasswordErrorCode | 密码错误 | 401 |
| 30003 | UserNotVerifyCode | 用户未验证 | 401 |
| 30004 | UserLockedCode | 用户已锁定 | 401 |
| 30005 | UserDisabledCode | 用户已禁用 | 401 |
| 30006 | UserExpiredCode | 用户已过期 | 401 |
| 30007 | UserAlreadyExistsCode | 用户已存在 | 401 |
| 30008 | UserNameOrPasswordErrorCode | 用户名或密码错误 | 401 |
| 30009 | UserAuthFailedCode | 认证失败 | 401 |
| 30010 | UserNoPermissionCode | 没有权限 | 401 |
| 30011 | UserPasswordErrCode | 密码错误 | 401 |
| 30012 | UserNotExistCode | 用户不存在 | 401 |
| 30013 | UserTokenErrorCode | 登录凭证无效 | 401 |
| 30014 | UserTokenExpiredCode | 登录已过期，请重新登录 | 401 |

## 权限相关错误码 (40000-49999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 40000 | AccessDeniedCode | 访问被拒绝 | 403 |
| 40001 | CasbinServiceCode | Casbin服务错误 | 500 |

## 缓存相关错误码 (50000-59999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 50000 | ErrorCacheCode | 缓存操作失败 | 500 |
| 50001 | ErrorCacheTimeoutCode | 缓存操作超时 | 500 |
| 50002 | ErrorCacheKeyCode | 缓存键不存在 | 500 |
| 50003 | ErrorCacheValueCode | 缓存值错误 | 500 |

## 文件相关错误码 (60000-69999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 60000 | FileStroageCode | 文件存储失败 | 500 |