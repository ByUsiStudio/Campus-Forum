package im

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Controller IM控制器
type Controller struct{}

// NewController 创建IM控制器
func NewController() *Controller {
	return &Controller{}
}

// ChatMessage 聊天消息结构（与service.go保持一致）
type ChatMessage struct {
	Type           string `json:"type"`            // message, join, leave, ack
	MessageID      string `json:"message_id"`      // 消息ID
	ConversationID string `json:"conversation_id"` // 会话ID
	Content        string `json:"content"`         // 消息内容
	SenderID       uint   `json:"sender_id"`       // 发送者ID
	TargetID       string `json:"target_id"`       // 接收者ID
	Timestamp      int64  `json:"timestamp"`       // 时间戳
}

// MessageAck 消息确认
type MessageAck struct {
	Type      string `json:"type"`       // ack
	MessageID string `json:"message_id"` // 消息ID
	Timestamp int64  `json:"timestamp"`  // 时间戳
}

// Upgrader WebSocket升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// ClientManager 客户端管理器
type ClientManager struct {
	clients    map[uint]*websocket.Conn // userID -> conn
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mutex      sync.RWMutex
}

// Client 客户端
type Client struct {
	UserID uint
	Conn   *websocket.Conn
}

// GlobalClientManager 全局客户端管理器
var GlobalClientManager *ClientManager

func init() {
	GlobalClientManager = &ClientManager{
		clients:    make(map[uint]*websocket.Conn),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
	go GlobalClientManager.start()
}

// start 启动管理器
func (m *ClientManager) start() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client.UserID] = client.Conn
			m.mutex.Unlock()
			log.Printf("用户 %d WebSocket连接建立", client.UserID)

		case client := <-m.unregister:
			m.mutex.Lock()
			if conn, ok := m.clients[client.UserID]; ok && conn == client.Conn {
				delete(m.clients, client.UserID)
				client.Conn.Close()
				log.Printf("用户 %d WebSocket连接关闭", client.UserID)
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.mutex.RLock()
			for _, conn := range m.clients {
				conn.WriteMessage(websocket.TextMessage, message)
			}
			m.mutex.RUnlock()
		}
	}
}

// Register 注册客户端
func (m *ClientManager) Register(userID uint, conn *websocket.Conn) {
	m.register <- &Client{UserID: userID, Conn: conn}
}

// Unregister 注销客户端
func (m *ClientManager) Unregister(userID uint, conn *websocket.Conn) {
	m.unregister <- &Client{UserID: userID, Conn: conn}
}

// SendToUser 发送消息给指定用户
func (m *ClientManager) SendToUser(userID uint, message []byte) error {
	m.mutex.RLock()
	conn, ok := m.clients[userID]
	m.mutex.RUnlock()

	if !ok {
		return fmt.Errorf("用户 %d 不在线", userID)
	}

	return conn.WriteMessage(websocket.TextMessage, message)
}

// Broadcast 广播消息
func (m *ClientManager) Broadcast(message []byte) {
	m.broadcast <- message
}

// IsOnline 检查用户是否在线
func (m *ClientManager) IsOnline(userID uint) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, ok := m.clients[userID]
	return ok
}

// GetOnlineUsers 获取所有在线用户
func (m *ClientManager) GetOnlineUsers() []uint {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	users := make([]uint, 0, len(m.clients))
	for userID := range m.clients {
		users = append(users, userID)
	}
	return users
}

// WebSocketHandle WebSocket处理函数
func (c *Controller) WebSocketHandle(ctx *gin.Context) {
	// 获取用户ID
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// 升级为WebSocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 注册客户端
	GlobalClientManager.Register(uint(userID), conn)

	// 处理消息
	defer GlobalClientManager.Unregister(uint(userID), conn)
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息错误: %v", err)
			}
			break
		}

		// 解析消息
		var msg ChatMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("解析消息失败: %v", err)
			continue
		}

		msg.SenderID = uint(userID)
		msg.Timestamp = time.Now().UnixMilli()

		// 生成消息ID
		msg.MessageID = fmt.Sprintf("%d-%d", userID, msg.Timestamp)

		// 构造完整消息
		fullMsg, _ := json.Marshal(msg)

		// 发送确认
		ack := MessageAck{
			Type:      "ack",
			MessageID: msg.MessageID,
			Timestamp: msg.Timestamp,
		}
		ackData, _ := json.Marshal(ack)
		conn.WriteMessage(websocket.TextMessage, ackData)

		// 私聊：发送给目标用户
		if msg.TargetID != "" {
			if targetID, err := strconv.ParseUint(msg.TargetID, 10, 32); err == nil {
				GlobalClientManager.SendToUser(uint(targetID), fullMsg)
			}
		}

		// 群聊：广播给所有在线用户
		if msg.Type == "group" {
			GlobalClientManager.Broadcast(fullMsg)
		}
	}
}

// SendMessage 发送消息API
func (c *Controller) SendMessage(ctx *gin.Context) {
	var req struct {
		SenderID   uint   `json:"sender_id" binding:"required"`
		TargetID   string `json:"target_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
		ReceiverID uint   `json:"receiver_id"` // 接收者ID（用户ID）
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := ChatMessage{
		Type:      "message",
		Content:   req.Content,
		SenderID:  req.SenderID,
		TargetID:  req.TargetID,
		Timestamp: time.Now().UnixMilli(),
	}
	msg.MessageID = fmt.Sprintf("%d-%d", req.SenderID, msg.Timestamp)

	msgData, _ := json.Marshal(msg)

	// 发送给目标用户
	if req.ReceiverID > 0 {
		if err := GlobalClientManager.SendToUser(req.ReceiverID, msgData); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "消息已发送，但目标用户不在线",
				"msg_id":  msg.MessageID,
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "消息发送成功",
		"msg_id":    msg.MessageID,
		"timestamp": msg.Timestamp,
	})
}

// GetOnlineStatus 获取用户在线状态
func (c *Controller) GetOnlineStatus(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"online":    GlobalClientManager.IsOnline(uint(userID)),
		"timestamp": time.Now().UnixMilli(),
	})
}

// GetOnlineUsers 获取在线用户列表
func (c *Controller) GetOnlineUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"users":     GlobalClientManager.GetOnlineUsers(),
		"count":     len(GlobalClientManager.GetOnlineUsers()),
		"timestamp": time.Now().UnixMilli(),
	})
}
