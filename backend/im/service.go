package im

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Type           string `json:"type"`            // message, join, leave, ack
	MessageID      string `json:"message_id"`      // 消息ID
	ConversationID string `json:"conversation_id"` // 会话ID
	Content        string `json:"content"`         // 消息内容
	SenderID       uint   `json:"sender_id"`       // 发送者ID
	TargetID       string `json:"target_id"`       // 接收者ID（用户ID或群ID）
	Timestamp      int64  `json:"timestamp"`       // 时间戳
}

// IMConfig IM配置
type IMConfig struct {
	ListenAddr string // WebSocket监听地址
}

// IMService IM服务
type IMService struct {
	config     *IMConfig
	clients    map[uint]*websocket.Conn // userID -> conn
	register   chan *Client
	unregister chan *Client
	broadcast  chan *ChatMessage
	mutex      sync.RWMutex
}

// Client 客户端
type Client struct {
	UserID uint
	Conn   *websocket.Conn
}

// 单例
var imService *IMService
var once sync.Once

// GetIMService 获取IM服务单例
func GetIMService() *IMService {
	once.Do(func() {
		imService = &IMService{
			clients:    make(map[uint]*websocket.Conn),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			broadcast:  make(chan *ChatMessage, 256),
			config: &IMConfig{
				ListenAddr: ":8081",
			},
		}
	})
	return imService
}

// SetConfig 设置配置
func (s *IMService) SetConfig(addr string) {
	s.config.ListenAddr = addr
}

// Start 启动服务
func (s *IMService) Start() error {
	log.Printf("IM服务初始化完成，WebSocket地址: %s", s.config.ListenAddr)

	// 启动注册/注销处理循环
	go s.handleLoop()

	return nil
}

// handleLoop 处理注册/注销
func (s *IMService) handleLoop() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client.UserID] = client.Conn
			s.mutex.Unlock()
			log.Printf("用户 %d 连接建立", client.UserID)

		case client := <-s.unregister:
			s.mutex.Lock()
			if conn, ok := s.clients[client.UserID]; ok && conn == client.Conn {
				delete(s.clients, client.UserID)
				log.Printf("用户 %d 连接关闭", client.UserID)
			}
			s.mutex.Unlock()
		}
	}
}

// Register 注册客户端
func (s *IMService) Register(client *Client) {
	s.register <- client
}

// Unregister 注销客户端
func (s *IMService) Unregister(client *Client) {
	s.unregister <- client
}

// Broadcast 广播消息
func (s *IMService) Broadcast(msg *ChatMessage) {
	s.broadcast <- msg
}

// SendToUser 发送消息给指定用户
func (s *IMService) SendToUser(userID uint, msg *ChatMessage) error {
	s.mutex.RLock()
	conn, ok := s.clients[userID]
	s.mutex.RUnlock()

	if !ok {
		return fmt.Errorf("用户 %d 不在线", userID)
	}

	data, _ := json.Marshal(msg)
	return conn.WriteMessage(websocket.TextMessage, data)
}

// IsOnline 检查用户是否在线
func (s *IMService) IsOnline(userID uint) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.clients[userID]
	return ok
}

// GetOnlineUsers 获取所有在线用户
func (s *IMService) GetOnlineUsers() []uint {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]uint, 0, len(s.clients))
	for userID := range s.clients {
		users = append(users, userID)
	}
	return users
}

// CreateMessage 创建消息
func CreateMessage(senderID uint, targetID, content string) *ChatMessage {
	return &ChatMessage{
		Type:      "message",
		Content:   content,
		SenderID:  senderID,
		TargetID:  targetID,
		Timestamp: time.Now().UnixMilli(),
	}
}
