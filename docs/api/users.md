# 用户接口

## 获取用户信息

**GET** `/api/users/{id}`

获取用户公开信息。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 用户 ID |

**响应**:
```json
{
  "id": 1,
  "username": "string",
  "display_name": "string",
  "avatar": "string",
  "signature": "string",
  "role": "user",
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

## 获取用户文章

**GET** `/api/users/{id}/articles`

获取指定用户的文章列表。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 用户 ID |

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |

**响应**:
```json
{
  "articles": [...],
  "total": 10,
  "page": 1,
  "page_size": 20,
  "total_pages": 1
}
```

---

## 获取用户关注列表

**GET** `/api/users/{id}/following`

获取指定用户的关注列表。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 用户 ID |

**响应**:
```json
{
  "following": [...]
}
```

---

## 获取用户粉丝列表

**GET** `/api/users/{id}/followers`

获取指定用户的粉丝列表。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 用户 ID |

**响应**:
```json
{
  "followers": [...]
}
```

---

## 获取个人资料

**GET** `/api/profile`

获取当前用户的个人资料（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "id": 1,
  "username": "string",
  "display_name": "string",
  "avatar": "string",
  "signature": "string",
  "role": "user",
  "status": "normal",
  "sign_in_days": 5,
  "total_sign_ins": 12,
  "total_coins": 100,
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

## 更新个人资料

**PUT** `/api/profile`

更新当前用户的个人资料（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

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

## 上传头像

**POST** `/api/upload/avatar`

上传用户头像（需认证）。

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求体**:
```
file: <图片文件>
```

**响应**:
```json
{
  "message": "上传成功",
  "avatar_url": "string"
}
```

---

## 上传图片

**POST** `/api/upload/image`

上传图片（需认证）。

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求体**:
```
file: <图片文件>
```

**响应**:
```json
{
  "message": "上传成功",
  "url": "string"
}
```

---

## 上传视频

**POST** `/api/upload/video`

上传视频（需认证）。

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求体**:
```
file: <视频文件>
```

**响应**:
```json
{
  "message": "上传成功",
  "url": "string"
}
```

---

## 上传语音

**POST** `/api/upload/voice`

上传语音（需认证）。

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求体**:
```
file: <语音文件>
```

**响应**:
```json
{
  "message": "上传成功",
  "url": "string"
}
```