# 权限组接口

## 权限级别说明

系统使用基于级别的权限控制，级别越高权限越大：

| 级别 | 角色名称 | 描述 |
|------|----------|------|
| 1 | 新人 | 新注册用户默认权限组 |
| 10 | 普通用户 | 普通用户权限组 |
| 50 | 版主 | 版主权限组，可管理板块内容 |
| 60 | 内容审核员 | 内容审核员权限组 |
| 80 | admin | 管理员权限组，只能管理普通用户内容 |
| 100 | system | 系统管理员权限组，拥有所有权限 |

**权限说明**：
- `admin` (级别 80)：只能管理普通用户的评论、文章，标记删除用户，更新公告和全部通知
- `system` (级别 100)：拥有所有权限，包括系统配置、权限组初始化、系统日志等
- `system` 组的通知不能被 `admin` 组进行任何更改

---

## 获取权限组列表

**GET** `/api/permission-groups`

获取所有权限组列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "groups": [
    {
      "id": 1,
      "name": "新人",
      "description": "新注册用户默认权限组",
      "level": 1,
      "permissions": "[\"article:read\", \"comment:create\", \"user:profile:view\"]",
      "is_default": true,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 获取权限组详情

**GET** `/api/permission-groups/{id}`

获取权限组详情（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 权限组 ID |

**响应**:
```json
{
  "group": {...}
}
```

---

## 创建权限组（需级别 80+）

**POST** `/api/permission-groups`

创建权限组（需认证，权限级别 >= 80）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "level": 1,
  "permissions": ["article:read", "comment:create"],
  "is_default": false
}
```

**响应**:
```json
{
  "message": "创建成功",
  "group": {...}
}
```

---

## 更新权限组（需级别 80+）

**PUT** `/api/permission-groups/{id}`

更新权限组（需认证，权限级别 >= 80）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 权限组 ID |

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "level": 1,
  "permissions": ["article:read"],
  "is_default": false,
  "is_active": true
}
```

**响应**:
```json
{
  "message": "更新成功",
  "group": {...}
}
```

---

## 删除权限组（需级别 80+）

**DELETE** `/api/permission-groups/{id}`

删除权限组（需认证，权限级别 >= 80）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 权限组 ID |

**响应**:
```json
{
  "message": "删除成功"
}
```

---

## 授予权限组（需级别 80+）

**POST** `/api/permission-groups/grant`

授予用户权限组（需认证，权限级别 >= 80）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 1,
  "group_id": 1,
  "expires_in_days": 30
}
```

**参数说明**:
| 参数 | 类型 | 描述 |
|------|------|------|
| user_id | int | 用户 ID |
| group_id | int | 权限组 ID |
| expires_in_days | int | 过期天数，0 表示永久 |

**响应**:
```json
{
  "message": "权限组授予成功",
  "user_group": {...}
}
```

---

## 撤销权限组（需级别 80+）

**DELETE** `/api/permission-groups/{id}/revoke-user/{user_id}`

撤销用户的权限组（需认证，权限级别 >= 80）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 权限组 ID |
| user_id | int | 用户 ID |

**响应**:
```json
{
  "message": "权限组已撤销"
}
```

---

## 获取用户权限组

**GET** `/api/users/{id}/permission-groups`

获取用户的权限组（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 用户 ID |

**响应**:
```json
{
  "groups": [
    {
      "id": 1,
      "user_id": 1,
      "permission_group_id": 1,
      "permission_group": {...},
      "granted_by": 1,
      "expires_at": null
    }
  ]
}
```

---

## 检查用户权限

**GET** `/api/permissions/check`

检查用户权限（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| permissions | string | 权限名称（可多个） |

**响应**:
```json
{
  "has_permission": true,
  "permissions": ["article:read", "comment:create"],
  "is_admin": false
}
```

---

## 初始化默认权限组（需 system 角色）

**POST** `/api/permission-groups/init`

初始化默认权限组（需认证，system 角色）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "message": "默认权限组初始化成功",
  "groups": [...]
}
```

**默认权限组列表**:
- 新人（级别 1）
- 普通用户（级别 10）
- 版主（级别 50）
- 内容审核员（级别 60）
- admin（级别 80）
- system（级别 100）