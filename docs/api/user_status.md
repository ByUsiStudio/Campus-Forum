# 用户状态接口

## 更新用户状态

**POST** `/api/user/status`

更新当前用户的在线状态（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "status": "online"
}
```

**响应**:
```json
{
  "message": "更新成功"
}
```

---

## 获取用户状态

**GET** `/api/user/status/{id}`

获取指定用户的在线状态（需认证）。

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
  "user_id": 1,
  "status": "online",
  "last_active_at": "2024-01-01T00:00:00Z"
}
```

---

## 获取所有用户状态（管理员）

**GET** `/api/users/status`

获取所有用户的在线状态（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "statuses": [...]
}
```

---

## 获取在线用户列表（管理员）

**GET** `/api/users/online`

获取在线用户列表（需认证，管理员）。

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

## 清理用户状态（管理员）

**POST** `/api/users/status/cleanup`

清理过期的用户状态（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "message": "清理成功"
}
```