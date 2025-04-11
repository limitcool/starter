# 应用错误码定义

本文档定义了应用中使用的所有错误码，包括错误码值、错误消息和HTTP状态码。

## 基础错误码

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 0 | SuccessCode | Success | 200 |

## 通用错误码 (10000-19999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 10000 | InvalidParamsCode | Invalid request parameters | 400 |
| 10001 | ErrorUnknownCode | Server is busy, please try again later | 500 |
| 10002 | ErrorNotExistCertCode | Authentication type does not exist | 400 |
| 10003 | ErrorNotFoundCode | Resource not found | 404 |
| 10004 | ErrorDatabaseCode | Database operation failed | 500 |
| 10005 | ErrorInternalCode | Internal server error | 500 |
| 10006 | ErrorCode | Error | 500 |
| 10007 | ErrorParamCode | Parameter error | 400 |

## 数据库错误码 (20000-29999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 20000 | DatabaseInsertErrorCode | Database insert failed | 500 |
| 20001 | DatabaseDeleteErrorCode | Database delete failed | 500 |
| 20002 | DatabaseQueryErrorCode | Database query failed | 500 |

## 用户相关错误码 (30000-39999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 30000 | UserNoLoginCode | User not logged in | 401 |
| 30001 | UserNotFoundCode | User not found | 404 |
| 30002 | UserPasswordErrorCode | Incorrect password | 401 |
| 30003 | UserNotVerifyCode | User not verified | 401 |
| 30004 | UserLockedCode | User is locked | 401 |
| 30005 | UserDisabledCode | User is disabled | 401 |
| 30006 | UserExpiredCode | User account expired | 401 |
| 30007 | UserAlreadyExistsCode | User already exists | 401 |
| 30008 | UserNameOrPasswordErrorCode | Incorrect username or password | 401 |
| 30009 | UserAuthFailedCode | Authentication failed | 401 |
| 30010 | UserNoPermissionCode | No permission | 401 |
| 30011 | UserPasswordErrCode | Password error | 401 |
| 30012 | UserNotExistCode | User does not exist | 401 |
| 30013 | UserTokenErrorCode | Invalid login credentials | 401 |
| 30014 | UserTokenExpiredCode | Login expired, please login again | 401 |

## 权限相关错误码 (40000-49999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 40000 | AccessDeniedCode | Access denied | 403 |
| 40001 | CasbinServiceCode | Casbin service error | 500 |

## 缓存相关错误码 (50000-59999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 50000 | ErrorCacheCode | Cache operation failed | 500 |
| 50001 | ErrorCacheTimeoutCode | Cache operation timeout | 500 |
| 50002 | ErrorCacheKeyCode | Cache key does not exist | 500 |
| 50003 | ErrorCacheValueCode | Cache value error | 500 |

## 文件相关错误码 (60000-69999)

| 错误码 | 名称 | 错误消息 | HTTP状态码 |
|------|------|--------|----------|
| 60000 | FileStroageCode | File storage failed | 500 |