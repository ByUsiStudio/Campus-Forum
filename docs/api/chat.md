# 聊天接口

## WebSocket 连接

**WS** `/ws`

建立 WebSocket 连接进行实时聊天（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**连接成功后可接收的消息类型**:

### 接收聊天消息
```json
{
  "type": "chat",
  "conversation_id": "123",
  "message_id": 456,
  "sender_id": 789,
  "content": "你好",
  "type": "text",
  "created_at": "2024-01-01T12:00:00Z"
}
```

### 发送聊天消息
```json
{
  "type": "chat",
  "data": {
    "conversation_id": "123",
    "content": "你好",
    "type": "text"
  }
}
```

---

## 获取会话列表

**GET** `/api/chat/conversations`

获取当前用户的会话列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "conversations": [
    {
      "id": 1,
      "type": "private",
      "target_user": {
        "id": 20,
        "username": "laowang",
        "avatar": "..."
      },
      "unread_count": 5,
      "last_message": {
        "id": 100,
        "content": "最后一条消息",
        "created_at": "2024-01-01T12:00:00Z"
      },
      "last_msg_time": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

## 获取消息历史

**GET** `/api/chat/messages`

获取指定会话的消息历史（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**查询参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| conversation_id | int | 会话 ID |
| limit | int | 每页数量，默认 20 |
| offset | int | 偏移量，默认 0 |

**响应**:
```json
{
  "messages": [
    {
      "id": 100,
      "conversation_id": 1,
      "sender_id": 20,
      "content": "你好",
      "type": "text",
      "status": 1,
      "created_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

## 发送消息

**POST** `/api/chat/messages`

发送消息到指定会话（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "conversation_id": 1,
  "content": "你好",
  "type": "text"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| conversation_id | int | 是 | 会话 ID |
| content | string | 是 | 消息内容 |
| type | string | 否 | 消息类型，默认 text |

**响应**:
```json
{
  "message": {
    "id": 101,
    "conversation_id": 1,
    "sender_id": 10,
    "content": "你好",
    "type": "text",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

---

## 创建私聊会话

**POST** `/api/chat/conversations/private`

创建与指定用户的私聊会话（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "target_user_id": 20
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| target_user_id | int | 是 | 目标用户 ID |

**响应**:
```json
{
  "conversation_id": 1
}
```

---

## 获取未读消息数

**GET** `/api/chat/unread`

获取当前用户的未读消息总数（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "unread_count": 15
}
```

---

## 标记会话已读

**POST** `/api/chat/conversations/:id/read`

标记指定会话为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 会话 ID |

**响应**:
```json
{
  "message": "已标记为已读"
}
```
