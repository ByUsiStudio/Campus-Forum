# 权限组接口

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
  "permission_groups": [...]
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
  "permission_group": {...}
}
```

---

## 创建权限组（管理员）

**POST** `/api/permission-groups`

创建权限组（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "permissions": ["string"]
}
```

**响应**:
```json
{
  "message": "创建成功",
  "permission_group": {...}
}
```

---

## 更新权限组（管理员）

**PUT** `/api/permission-groups/{id}`

更新权限组（需认证，管理员）。

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
  "permissions": ["string"]
}
```

**响应**:
```json
{
  "message": "更新成功",
  "permission_group": {...}
}
```

---

## 删除权限组（管理员）

**DELETE** `/api/permission-groups/{id}`

删除权限组（需认证，管理员）。

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

## 授予权限组（管理员）

**POST** `/api/permission-groups/grant`

授予用户权限组（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 1,
  "permission_group_id": 1
}
```

**响应**:
```json
{
  "message": "授予成功"
}
```

---

## 撤销权限组（管理员）

**DELETE** `/api/permission-groups/{id}/revoke-user/{user_id}`

撤销用户的权限组（需认证，管理员）。

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
  "message": "撤销成功"
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
  "permission_groups": [...]
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
| permission | string | 权限名称 |

**响应**:
```json
{
  "has_permission": true
}
```

---

## 初始化默认权限组（管理员）

**POST** `/api/permission-groups/init`

初始化默认权限组（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "message": "初始化成功"
}
```