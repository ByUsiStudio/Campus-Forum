import { ref, onMounted, onUnmounted } from 'vue'
import axios from 'axios'

// WebSocket连接状态
const WS_STATE = {
  CONNECTING: 0,
  OPEN: 1,
  CLOSING: 2,
  CLOSED: 3
}

export function useGoIM() {
  const ws = ref(null)
  const isConnected = ref(false)
  const error = ref(null)
  const messages = ref([])
  const conversations = ref([])
  const onlineUsers = ref([])
  
  let reconnectTimer = null
  let heartbeatTimer = null
  let currentUserId = null
  let messageHandlers = []
  let connectionHandlers = []

  // 通知连接状态变化
  const notifyConnectionStatus = (connected) => {
    connectionHandlers.forEach(handler => handler(connected))
  }

  // 通知新消息
  const notifyNewMessage = (message) => {
    messageHandlers.forEach(handler => handler(message))
  }

  // WebSocket连接
  const connectWebSocket = (userId) => {
    currentUserId = userId
    
    // 清除之前的连接
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }

    // 清除定时器
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }

    // 获取API基础URL
    const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
    const wsUrl = `${apiBase.replace('http', 'ws')}/ws?user_id=${userId}`

    try {
      ws.value = new WebSocket(wsUrl)

      ws.value.onopen = () => {
        console.log('IM WebSocket连接成功')
        isConnected.value = true
        error.value = null
        startHeartbeat()
        notifyConnectionStatus(true)
      }

      ws.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          handleMessage(data)
        } catch (err) {
          console.error('解析WebSocket消息失败:', err)
        }
      }

      ws.value.onerror = (err) => {
        console.error('WebSocket错误:', err)
        error.value = '连接错误'
        notifyConnectionStatus(false)
      }

      ws.value.onclose = (event) => {
        console.log('WebSocket连接关闭:', event.code, event.reason)
        isConnected.value = false
        notifyConnectionStatus(false)
        
        // 尝试重连
        if (event.code !== 1000 && currentUserId) {
          scheduleReconnect(currentUserId)
        }
      }
    } catch (err) {
      console.error('建立WebSocket连接失败:', err)
      error.value = err.message || '连接失败'
    }
  }

  // 处理消息
  const handleMessage = (data) => {
    switch (data.type) {
      case 'ack':
        // 消息确认
        console.log('消息已送达:', data.message_id)
        break
      case 'message':
        // 新消息
        notifyNewMessage(data)
        messages.value.push(data)
        break
      case 'system':
        // 系统消息
        console.log('系统消息:', data.content)
        break
      default:
        console.log('未知消息类型:', data.type)
    }
  }

  // 发送消息
  const sendMessage = (targetId, content, type = 'private') => {
    if (!ws.value || ws.value.readyState !== WS_STATE.OPEN) {
      console.error('WebSocket未连接')
      return Promise.reject(new Error('WebSocket未连接'))
    }

    const message = {
      type: type === 'group' ? 'group' : 'message',
      target_id: String(targetId),
      content: content
    }

    ws.value.send(JSON.stringify(message))
    
    return Promise.resolve()
  }

  // 发送私聊消息
  const sendPrivateMessage = (targetId, content) => {
    return sendMessage(targetId, content, 'private')
  }

  // 发送群消息
  const sendGroupMessage = (groupId, content) => {
    return sendMessage(groupId, content, 'group')
  }

  // 开始心跳
  const startHeartbeat = () => {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
    }
    
    heartbeatTimer = setInterval(() => {
      if (ws.value && ws.value.readyState === WS_STATE.OPEN) {
        ws.value.send(JSON.stringify({ type: 'ping' }))
      }
    }, 30000) // 30秒心跳
  }

  // 定时重连
  const scheduleReconnect = (userId) => {
    if (reconnectTimer) return
    
    reconnectTimer = setTimeout(() => {
      console.log('尝试重新连接...')
      connectWebSocket(userId)
      reconnectTimer = null
    }, 3000) // 3秒后重连
  }

  // 断开连接
  const disconnect = () => {
    currentUserId = null
    
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
    
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    
    if (ws.value) {
      ws.value.close(1000, '用户主动断开')
      ws.value = null
    }
    
    isConnected.value = false
  }

  // 获取在线状态
  const getOnlineStatus = async (userId) => {
    try {
      const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
      const response = await axios.get(`${apiBase}/api/im/online-status`, {
        params: { user_id: userId }
      })
      return response.data
    } catch (err) {
      console.error('获取在线状态失败:', err)
      throw err
    }
  }

  // 获取在线用户列表
  const getOnlineUsers = async () => {
    try {
      const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
      const response = await axios.get(`${apiBase}/api/im/online-users`)
      onlineUsers.value = response.data.users || []
      return response.data
    } catch (err) {
      console.error('获取在线用户失败:', err)
      throw err
    }
  }

  // 通过API发送消息（可选的后备方式）
  const sendMessageViaAPI = async (senderId, receiverId, content) => {
    try {
      const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
      const response = await axios.post(`${apiBase}/api/im/send-message`, {
        sender_id: senderId,
        receiver_id: receiverId,
        target_id: String(receiverId),
        content: content
      })
      return response.data
    } catch (err) {
      console.error('API发送消息失败:', err)
      throw err
    }
  }

  // 订阅消息
  const onMessage = (handler) => {
    messageHandlers.push(handler)
    return () => {
      messageHandlers = messageHandlers.filter(h => h !== handler)
    }
  }

  // 订阅连接状态变化
  const onConnectionChange = (handler) => {
    connectionHandlers.push(handler)
    return () => {
      connectionHandlers = connectionHandlers.filter(h => h !== handler)
    }
  }

  // 组件卸载时清理
  onUnmounted(() => {
    disconnect()
  })

  return {
    // 状态
    isConnected,
    error,
    messages,
    conversations,
    onlineUsers,
    
    // 方法
    connectWebSocket,
    disconnect,
    sendMessage,
    sendPrivateMessage,
    sendGroupMessage,
    getOnlineStatus,
    getOnlineUsers,
    sendMessageViaAPI,
    
    // 事件订阅
    onMessage,
    onConnectionChange
  }
}
