# 管理员接口

## 检查管理员权限

**GET** `/api/admin/check`

检查当前用户是否为管理员（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "is_admin": true
}
```

---

## 获取统计数据

**GET** `/api/admin/statistics`

获取系统统计数据（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "user_count": 100,
  "article_count": 500,
  "comment_count": 2000,
  "report_count": 10
}
```

---

## 获取所有用户

**GET** `/api/admin/users`

获取所有用户列表（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "users": [...]
}
```

---

## 更新用户信息

**PUT** `/api/admin/users/{id}`

更新用户信息（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 用户 ID |

**请求体**:
```json
{
  "display_name": "string",
  "signature": "string"
}
```

**响应**:
```json
{
  "message": "更新成功",
  "user": {...}
}
```

---

## 更新用户角色

**PUT** `/api/admin/users/{id}/role`

更新用户角色（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 用户 ID |

**请求体**:
```json
{
  "role": "admin"
}
```

**响应**:
```json
{
  "message": "更新成功"
}
```

---

## 封禁用户

**POST** `/api/admin/users/{id}/ban`

封禁用户（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 用户 ID |

**请求体**:
```json
{
  "reason": "string"
}
```

**响应**:
```json
{
  "message": "封禁成功"
}
```

---

## 解封用户

**POST** `/api/admin/users/{id}/unban`

解封用户（需认证，管理员）。

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
  "message": "解封成功"
}
```

---

## 删除用户

**DELETE** `/api/admin/users/{id}`

删除用户（需认证，管理员）。

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
  "message": "删除成功"
}
```

---

## 获取所有文章

**GET** `/api/admin/articles`

获取所有文章（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "articles": [...]
}
```

---

## 更新文章状态

**PUT** `/api/admin/articles/{id}/status`

更新文章状态（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**请求体**:
```json
{
  "status": "published"
}
```

**响应**:
```json
{
  "message": "更新成功"
}
```

---

## 获取所有评论

**GET** `/api/admin/comments`

获取所有评论（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "comments": [...]
}
```

---

## 删除评论（管理员）

**DELETE** `/api/admin/comments/{id}`

删除评论（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 评论 ID |

**响应**:
```json
{
  "message": "删除成功"
}
```