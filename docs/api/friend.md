# 好友接口

## 发送好友请求

**POST** `/api/friends/request`

发送好友请求（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 123,
  "message": "你好，我想加你为好友"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| user_id | int | 是 | 目标用户 ID |
| message | string | 否 | 请求消息 |

**响应**:
```json
{
  "message": "好友请求已发送"
}
```

---

## 同意好友请求

**POST** `/api/friends/request/:id/accept`

同意好友请求（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 好友请求 ID |

**响应**:
```json
{
  "message": "已同意好友请求"
}
```

---

## 拒绝好友请求

**POST** `/api/friends/request/:id/reject`

拒绝好友请求（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 好友请求 ID |

**响应**:
```json
{
  "message": "已拒绝好友请求"
}
```

---

## 删除好友

**DELETE** `/api/friends/:id`

删除好友（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 好友关系 ID |

**响应**:
```json
{
  "message": "已删除好友"
}
```

---

## 获取好友列表

**GET** `/api/friends`

获取当前用户的好友列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "friends": [
    {
      "id": 1,
      "user_id": 10,
      "friend_id": 20,
      "display_name": "老王",
      "status": 1,
      "friend": {
        "id": 20,
        "username": "laowang",
        "avatar": "..."
      }
    }
  ]
}
```

---

## 获取收到的好友请求

**GET** `/api/friends/requests`

获取收到的好友请求列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "requests": [
    {
      "id": 1,
      "sender_id": 20,
      "receiver_id": 10,
      "message": "你好，我想加你为好友",
      "status": 0,
      "sender": {
        "id": 20,
        "username": "laowang",
        "avatar": "..."
      }
    }
  ]
}
```

**状态说明**:
- `0`: 待处理
- `1`: 已同意
- `2`: 已拒绝

---

## 获取发送的好友请求

**GET** `/api/friends/sent-requests`

获取发送的好友请求列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "requests": [...]
}
```

---

## 更新好友备注名

**PUT** `/api/friends/:id/display-name`

更新好友备注名（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 好友关系 ID |

**请求体**:
```json
{
  "display_name": "同事小王"
}
```

**响应**:
```json
{
  "message": "备注名已更新"
}
```

---

## 检查好友状态

**GET** `/api/friends/status/:id`

检查与指定用户的好友关系状态（需认证）。

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
  "status": "friends",
  "friend_id": 1
}
```

**状态返回值**:
- `none`: 不是好友，无请求
- `pending_sent`: 已发送请求，等待对方同意
- `pending_received`: 收到请求，等待处理
- `friends`: 已是好友

---

## 获取共同好友

**GET** `/api/friends/mutual/:id`

获取与指定用户的共同好友（需认证）。

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
  "mutual_friends": [
    {
      "id": 30,
      "username": "common_friend",
      "avatar": "..."
    }
  ]
}
```
