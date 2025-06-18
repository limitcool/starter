# 文件存储系统使用指南

本文档详细说明了系统中本地存储和MinIO存储的完整上传下载逻辑。

## 概述

系统支持两种存储方式：
- **本地存储**：文件存储在应用服务器本地磁盘
- **MinIO存储**：文件存储在MinIO对象存储服务

两种存储方式使用**完全相同的API接口**，通过配置文件切换，对前端透明。

## 配置切换

### 本地存储配置
```yaml
Storage:
  Enabled: true
  Type: local
  Local:
    Path: storage
    URL: http://localhost:8080/static
```

### MinIO存储配置
```yaml
Storage:
  Enabled: true
  Type: s3
  S3:
    Endpoint: "http://1.1.1.1:9000"
    Region: us-east-1
    Bucket: starter-test
    AccessKey: "your-access-key"
    SecretKey: "your-secret-key"
    UseSSL: false
```

## API接口

### 管理员接口

#### 1. 获取上传URL
```http
POST /api/v1/admin/files/upload-url
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "filename": "test.jpg",
  "content_type": "image/jpeg",
  "is_public": false,
  "usage": "avatar",
  "size": 1024000
}
```

#### 2. 确认上传完成
```http
POST /api/v1/admin/files/confirm
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "file_id": "uuid",
  "size": 1024000
}
```

#### 3. 获取下载URL
```http
GET /api/v1/admin/files/{file_id}/download
Authorization: Bearer {admin_token}
```

#### 4. 删除文件
```http
DELETE /api/v1/admin/files/{file_id}
Authorization: Bearer {admin_token}
```

### 普通用户接口

#### 直接文件上传
```http
POST /api/v1/upload/file
Authorization: Bearer {user_token}
Content-Type: multipart/form-data

filename: test.jpg
usage: avatar
is_public: false
file: [binary data]
```

### 公开接口

#### 获取文件信息
```http
GET /public/files/{file_id}
```

## 本地存储完整流程

### 管理员上传流程

#### 1. 获取上传URL
**请求**：
```http
POST /api/v1/admin/files/upload-url
{
  "filename": "avatar.jpg",
  "content_type": "image/jpeg",
  "is_public": false,
  "usage": "avatar"
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "file_id": "274b5c46-0e13-4ded-b190-5cdea9c37a30",
    "upload_url": "/api/v1/upload/file?path=private/users/avatars/user_123/uuid.jpg&public=false",
    "method": "POST",
    "expires_in": 15,
    "storage_type": "local",
    "usage": "avatar",
    "path_info": {
      "category": "avatar",
      "path": "users/avatars/user_123/uuid.jpg"
    }
  }
}
```

#### 2. 上传文件
前端使用返回的 `upload_url` 上传文件：
```http
POST /api/v1/upload/file?path=private/users/avatars/user_123/uuid.jpg&public=false
Authorization: Bearer {token}
Content-Type: multipart/form-data

filename: avatar.jpg
usage: avatar
is_public: false
file: [binary data]
```

#### 3. 确认上传
```http
POST /api/v1/admin/files/confirm
{
  "file_id": "274b5c46-0e13-4ded-b190-5cdea9c37a30",
  "size": 1024000
}
```

#### 4. 获取下载URL
```http
GET /api/v1/admin/files/274b5c46-0e13-4ded-b190-5cdea9c37a30/download
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "file_id": "274b5c46-0e13-4ded-b190-5cdea9c37a30",
    "filename": "avatar.jpg",
    "download_url": "http://localhost:8080/static/private/users/avatars/user_123/uuid.jpg",
    "is_public": false,
    "size": 1024000,
    "storage_type": "local"
  }
}
```

### 普通用户直接上传流程

```http
POST /api/v1/upload/file
Authorization: Bearer {user_token}
Content-Type: multipart/form-data

filename: document.pdf
usage: document
is_public: false
file: [binary data]
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "file_id": "new-uuid",
    "filename": "document.pdf",
    "download_url": "http://localhost:8080/static/private/documents/general/2025/06/uuid.pdf",
    "size": 2048000,
    "storage_type": "local",
    "is_public": false,
    "usage": "document"
  }
}
```

## MinIO存储完整流程

### 管理员上传流程

#### 1. 获取上传URL
**请求**：
```http
POST /api/v1/admin/files/upload-url
{
  "filename": "avatar.jpg",
  "content_type": "image/jpeg",
  "is_public": false,
  "usage": "avatar"
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "file_id": "274b5c46-0e13-4ded-b190-5cdea9c37a30",
    "upload_url": "https://minio.example.com/bucket/private/users/avatars/user_123/uuid.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...",
    "method": "PUT",
    "expires_in": 15,
    "storage_type": "minio",
    "usage": "avatar",
    "path_info": {
      "category": "avatar",
      "path": "users/avatars/user_123/uuid.jpg"
    }
  }
}
```

#### 2. 直接上传到MinIO
前端使用预签名URL直接上传到MinIO：
```http
PUT https://minio.example.com/bucket/private/users/avatars/user_123/uuid.jpg?X-Amz-Algorithm=...
Content-Type: image/jpeg

[binary data]
```

#### 3. 确认上传
```http
POST /api/v1/admin/files/confirm
{
  "file_id": "274b5c46-0e13-4ded-b190-5cdea9c37a30",
  "size": 1024000
}
```

#### 4. 获取下载URL
```http
GET /api/v1/admin/files/274b5c46-0e13-4ded-b190-5cdea9c37a30/download
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "file_id": "274b5c46-0e13-4ded-b190-5cdea9c37a30",
    "filename": "avatar.jpg",
    "download_url": "https://minio.example.com/bucket/private/users/avatars/user_123/uuid.jpg?X-Amz-Algorithm=...",
    "is_public": false,
    "size": 1024000,
    "storage_type": "minio"
  }
}
```

### 普通用户直接上传流程

对于MinIO存储，普通用户的直接上传仍然通过应用服务器：

```http
POST /api/v1/upload/file
Authorization: Bearer {user_token}
Content-Type: multipart/form-data

filename: document.pdf
usage: document
is_public: false
file: [binary data]
```

服务器接收文件后上传到MinIO，然后返回：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "file_id": "new-uuid",
    "filename": "document.pdf",
    "download_url": "https://minio.example.com/bucket/private/documents/general/2025/06/uuid.pdf?X-Amz-Algorithm=...",
    "size": 2048000,
    "storage_type": "minio",
    "is_public": false,
    "usage": "document"
  }
}
```

## 文件路径规则

### 路径结构
```
{public|private}/{category}/{subcategory}/{date}/{filename}
```

### 示例路径
- **头像**: `private/users/avatars/user_123/uuid.jpg`
- **横幅**: `public/content/banners/2025/06/uuid.png`
- **文档**: `private/documents/general/2025/06/uuid.pdf`
- **临时文件**: `private/temp/2025/06/18/uuid.tmp`

### 文件用途类型
- `avatar`: 用户头像
- `profile`: 用户资料图片
- `banner`: 横幅图片
- `document`: 文档文件
- `image`: 一般图片
- `video`: 视频文件
- `audio`: 音频文件
- `temp`: 临时文件

## 权限控制

### 管理员权限
- 可以获取任何文件的上传/下载URL
- 可以删除任何文件
- 可以访问所有文件管理接口

### 普通用户权限
- 只能直接上传文件
- 只能访问自己上传的文件
- 无法删除文件（需要管理员权限）

### 公开访问
- 任何人都可以访问公开文件的信息
- 公开文件可以直接通过URL访问

## 错误处理

### 常见错误码
- `5000000`: 无效的参数
- `4010000`: 用户未登录
- `4030000`: 权限不足
- `4040000`: 文件不存在
- `5001000`: 文件存储错误
- `5002000`: 数据库错误

### 错误响应格式
```json
{
  "code": 5000000,
  "message": "无效的参数",
  "request_id": "req-xxx",
  "timestamp": 1750233346
}
```

## 最佳实践

### 前端开发建议
1. **统一处理**：无论后端使用哪种存储，前端代码保持一致
2. **错误处理**：根据响应的 `storage_type` 字段调整上传逻辑
3. **进度显示**：MinIO直传可以显示真实上传进度
4. **重试机制**：预签名URL过期时重新获取

### 后端配置建议
1. **生产环境**：推荐使用MinIO存储，提高性能和可扩展性
2. **开发环境**：可以使用本地存储，简化部署
3. **混合部署**：可以根据文件类型选择不同存储方式

### 安全建议
1. **私有文件**：使用预签名URL，设置合理的过期时间
2. **公开文件**：设置适当的缓存策略
3. **文件验证**：上传前验证文件类型和大小
4. **权限控制**：严格控制文件访问权限
