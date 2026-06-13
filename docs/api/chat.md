# 聊天 API 文档

## 概述

校园论坛聊天功能基于 WebSocket 实现，支持私聊和群聊。使用 `goim` 作为即时通讯后端服务。

## WebSocket 连接

### 连接地址

```
ws://host:port/ws?user_id={userId}
```

### 连接参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | int | 是 | 用户ID |

### 连接示例

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?user_id=1')
```

## 消息格式

### 发送消息

```json
{
  "type": "message",           // 消息类型：message, group, ping
  "target_id": "2",           // 目标用户ID或群ID
  "content": "你好"           // 消息内容
}
```

### 接收消息

```json
{
  "type": "message",           // 消息类型
  "message_id": "1-1234567890", // 消息ID
  "conversation_id": "private_1_2", // 会话ID
  "content": "你好",           // 消息内容
  "sender_id": 1,             // 发送者ID
  "target_id": "2",           // 接收者ID
  "timestamp": 1234567890      // 时间戳（毫秒）
}
```

### 消息确认

```json
{
  "type": "ack",
  "message_id": "1-1234567890",
  "timestamp": 1234567890
}
```

## API 接口

### 获取在线状态

检查用户是否在线。

**请求**

```
GET /api/im/online-status?user_id={userId}
```

**响应**

```json
{
  "user_id": 1,
  "online": true,
  "timestamp": 1234567890
}
```

### 获取在线用户列表

获取所有在线用户。

**请求**

```
GET /api/im/online-users
```

**响应**

```json
{
  "users": [1, 2, 3],
  "count": 3,
  "timestamp": 1234567890
}
```

### 发送消息（HTTP API）

通过 HTTP API 发送消息（可作为 WebSocket 的备份方案）。

**请求**

```
POST /api/im/send-message
Content-Type: application/json

{
  "sender_id": 1,
  "receiver_id": 2,
  "target_id": "2",
  "content": "你好"
}
```

**响应**

```json
{
  "success": true,
  "message": "消息发送成功",
  "msg_id": "1-1234567890",
  "timestamp": 1234567890
}
```

## 前端集成

### 安装依赖

无需额外安装 npm 包，使用原生 WebSocket API。

### 使用示例

```javascript
import { useGoIM } from '@/composables/useGoIM'

const {
  isConnected,
  connectWebSocket,
  disconnect,
  sendPrivateMessage,
  sendGroupMessage,
  onMessage,
} = useGoIM()

// 连接
connectWebSocket(userId)

// 监听消息
const removeListener = onMessage((message) => {
  console.log('收到消息:', message)
})

// 发送私聊消息
await sendPrivateMessage(targetUserId, '你好')

// 发送群消息
await sendGroupMessage(groupId, '大家好')

// 断开连接
disconnect()
```

### 环境配置

确保前端配置了正确的 API 基础地址：

```env
VITE_API_BASE_URL=http://localhost:8080
```

## 消息类型

| 类型 | 说明 |
|------|------|
| `message` | 私聊消息 |
| `group` | 群聊消息 |
| `ping` | 心跳检测 |
| `ack` | 消息确认 |

## 错误处理

| 错误码 | 说明 |
|--------|------|
| 1000 | 正常关闭 |
| 1001 | 服务器关闭 |
| 1006 | 异常关闭（网络问题） |

## 注意事项

1. WebSocket 连接需要保持活跃，系统会自动进行心跳检测
2. 断线后会自动重连（3秒间隔）
3. 建议在页面可见性变化时检查连接状态
4. 消息历史需要通过后端 API 单独获取
