# Casbin 规则表说明

## 表结构

Casbin 规则表 `casbin_rule` 用于存储权限策略和角色继承关系，具有以下结构：

```sql
CREATE TABLE casbin_rule (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    ptype VARCHAR(100) DEFAULT NULL,
    v0 VARCHAR(100) DEFAULT NULL,
    v1 VARCHAR(100) DEFAULT NULL,
    v2 VARCHAR(100) DEFAULT NULL,
    v3 VARCHAR(100) DEFAULT NULL,
    v4 VARCHAR(100) DEFAULT NULL,
    v5 VARCHAR(100) DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY idx_casbin_rule (ptype,v0,v1,v2,v3,v4,v5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 字段说明

- `id`: 主键，自增长
- `ptype`: 策略类型，通常是 "p"（权限策略）或 "g"（角色继承关系）
- `v0` ~ `v5`: 策略的各个部分，具体含义如下：

### 权限策略 (ptype = "p")

对于权限策略，字段含义为：
- `v0`: 角色编码
- `v1`: 资源/对象
- `v2`: 操作/动作
- `v3`, `v4`, `v5`: 额外条件（可选）

例如：
```
p, admin, /api/users, GET
```
表示 `admin` 角色可以对 `/api/users` 资源执行 `GET` 操作。

### 角色继承关系 (ptype = "g")

对于角色继承关系，字段含义为：
- `v0`: 用户ID
- `v1`: 角色编码
- `v2`: 域（可选）

例如：
```
g, 10001, admin
```
表示用户 `10001` 拥有 `admin` 角色。

## 使用方式

在我们的系统中，Casbin 规则表由 Casbin 的 GORM 适配器自动管理。当我们调用以下方法时，适配器会自动更新表中的数据：

- `AddRoleForUser`: 为用户分配角色，添加 "g" 类型记录
- `DeleteRoleForUser`: 删除用户角色，删除 "g" 类型记录
- `AddPermissionForRole`: 为角色添加权限，添加 "p" 类型记录
- `DeletePermissionForRole`: 删除角色权限，删除 "p" 类型记录

## 查询示例

### 查询用户角色

```sql
SELECT v1 FROM casbin_rule WHERE ptype = 'g' AND v0 = '用户ID';
```

### 查询角色权限

```sql
SELECT v1, v2 FROM casbin_rule WHERE ptype = 'p' AND v0 = '角色编码';
```

### 查询所有角色

```sql
SELECT DISTINCT v1 FROM casbin_rule WHERE ptype = 'g';
```

## 注意事项

1. 不要直接修改 `casbin_rule` 表中的数据，应该通过 Casbin API 进行操作
2. 表中的 `v0` ~ `v5` 字段长度有限，不要存储过长的字符串
3. 删除用户或角色时，应该同时删除相关的 Casbin 规则
