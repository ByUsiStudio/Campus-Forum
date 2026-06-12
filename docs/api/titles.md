# 头衔接口

## 获取所有头衔

**GET** `/api/titles`

获取所有头衔列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "titles": [...]
}
```

---

## 创建头衔（管理员）

**POST** `/api/titles`

创建头衔（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "icon": "string",
  "color": "string"
}
```

**响应**:
```json
{
  "message": "创建成功",
  "title": {...}
}
```

---

## 更新头衔（管理员）

**PUT** `/api/titles/{id}`

更新头衔（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 头衔 ID |

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "icon": "string",
  "color": "string"
}
```

**响应**:
```json
{
  "message": "更新成功",
  "title": {...}
}
```

---

## 删除头衔（管理员）

**DELETE** `/api/titles/{id}`

删除头衔（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 头衔 ID |

**响应**:
```json
{
  "message": "删除成功"
}
```

---

## 授予头衔（管理员）

**POST** `/api/titles/grant`

授予用户头衔（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 1,
  "title_id": 1
}
```

**响应**:
```json
{
  "message": "授予成功"
}
```

---

## 撤销头衔（管理员）

**POST** `/api/titles/revoke`

撤销用户头衔（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 1,
  "title_id": 1
}
```

**响应**:
```json
{
  "message": "撤销成功"
}
```

---

## 获取用户头衔

**GET** `/api/users/{id}/titles`

获取用户的头衔列表（需认证）。

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
  "titles": [...]
}
```