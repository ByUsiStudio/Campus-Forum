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