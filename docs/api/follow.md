# 关注接口

## 关注用户

**POST** `/api/follow/{id}`

关注指定用户（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 被关注用户 ID |

**响应**:
```json
{
  "message": "关注成功"
}
```

---

## 取消关注

**DELETE** `/api/follow/{id}`

取消关注指定用户（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 被取消关注用户 ID |

**响应**:
```json
{
  "message": "取消关注成功"
}
```

---

## 获取关注列表

**GET** `/api/following`

获取当前用户的关注列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "following": [...]
}
```

---

## 获取粉丝列表

**GET** `/api/followers`

获取当前用户的粉丝列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "followers": [...]
}
```

---

## 检查关注状态

**GET** `/api/follow/status/{id}`

检查是否关注了指定用户（需认证）。

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
  "following": true
}
```

---

## 获取互相关注好友

**GET** `/api/mutual`

获取当前用户的互相关注好友列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "mutual_friends": [...]
}
```