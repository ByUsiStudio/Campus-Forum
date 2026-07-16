# 通知接口

## 获取通知列表

**GET** `/api/notifications`

获取当前用户的通知列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "notifications": [...]
}
```

---

## 获取未读通知数量

**GET** `/api/notifications/unread-count`

获取未读通知数量（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "count": 5
}
```

---

## 标记通知为已读

**POST** `/api/notifications/{id}/read`

标记单条通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 通知 ID |

**响应**:
```json
{
  "message": "标记成功"
}
```

---

## 标记所有通知为已读

**POST** `/api/notifications/read-all`

标记所有通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "message": "全部已读"
}
```

---

## 创建通知（管理员）

**POST** `/api/notifications`

创建系统通知（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "title": "string",
  "content": "string",
  "type": "string"
}
```

**响应**:
```json
{
  "message": "创建成功",
  "notification": {...}
}
```

---

## 获取评论回复通知

**GET** `/api/comment-reply-notifications`

获取评论回复通知（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "notifications": [...]
}
```

---

## 获取关注通知

**GET** `/api/follow-notifications`

获取关注用户的新内容通知（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "notifications": [...]
}
```

---

## 获取用户个人通知

**GET** `/api/user-notifications`

获取用户个人通知（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "notifications": [...]
}
```

---

## 发送用户个人通知（管理员）

**POST** `/api/user-notifications/send`

发送用户个人通知（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 1,
  "title": "string",
  "content": "string"
}
```

**响应**:
```json
{
  "message": "发送成功"
}
```

---

## 批量发送通知（管理员）

**POST** `/api/user-notifications/send-batch`

批量发送通知给多个用户（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_ids": [1, 2, 3],
  "type": "string",
  "title": "string",
  "content": "string",
  "related_type": "string",
  "related_id": 1,
  "link": "string",
  "priority": "string"
}
```

**响应**:
```json
{
  "message": "批量通知发送成功",
  "sent_count": 3
}
```

---

## 获取用户个人通知详情

**GET** `/api/user-notifications/{id}`

获取单条个人通知详情（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 通知 ID |

**响应**:
```json
{
  "notification": {
    "id": 1,
    "user_id": 1,
    "type": "string",
    "title": "string",
    "content": "string",
    "related_type": "string",
    "related_id": 1,
    "link": "string",
    "priority": "normal",
    "is_read": false,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 标记个人通知为已读

**POST** `/api/user-notifications/{id}/read`

标记单条个人通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 通知 ID |

**响应**:
```json
{
  "message": "标记成功"
}
```

---

## 标记所有个人通知为已读

**POST** `/api/user-notifications/read-all`

标记所有个人通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "message": "全部标记已读"
}
```

---

## 删除个人通知

**DELETE** `/api/user-notifications/{id}`

删除单条个人通知（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 通知 ID |

**响应**:
```json
{
  "message": "删除成功"
}
```

---

## 清空所有个人通知

**DELETE** `/api/user-notifications/clear`

清空当前用户所有个人通知（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "message": "清空成功"
}
```

---

## 管理员获取用户通知

**GET** `/api/admin/user-notifications/{user_id}`

获取指定用户的所有通知（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| user_id | int | 用户 ID |

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |

**响应**:
```json
{
  "notifications": [...],
  "total_pages": 1
}
```