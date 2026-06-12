package service

import (
	"encoding/json"
	"forum/database"
	"forum/models"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

// WSClient WebSocket客户端
type WSClient struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

// WSMessage WebSocket消息
type WSMessage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Content string      `json:"content,omitempty"`
}

// WSServer WebSocket服务器
type WSServer struct {
	clients    map[uint]*WSClient
	register   chan *WSClient
	unregister chan *WSClient
	broadcast  chan []byte
	mu         sync.RWMutex
}

var wsServer *WSServer

// InitWebSocket 初始化WebSocket服务器
func InitWebSocket() {
	wsServer = &WSServer{
		clients:    make(map[uint]*WSClient),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		broadcast:  make(chan []byte),
	}
	go wsServer.run()
}

// GetWSServer 获取WebSocket服务器实例
func GetWSServer() *WSServer {
	return wsServer
}

func (s *WSServer) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.UserID] = client
			s.mu.Unlock()
			log.Printf("用户 %d 已连接", client.UserID)

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.UserID]; ok {
				delete(s.clients, client.UserID)
				close(client.Send)
			}
			s.mu.Unlock()
			log.Printf("用户 %d 已断开", client.UserID)

		case message := <-s.broadcast:
			s.mu.RLock()
			for _, client := range s.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(s.clients, client.UserID)
				}
			}
			s.mu.RUnlock()
		}
	}
}

// HandleWebSocket 处理WebSocket连接
func HandleWebSocket(w http.ResponseWriter, r *http.Request, userID uint) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	client := &WSClient{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	wsServer.register <- client

	// 读取消息
	go func() {
		defer func() {
			wsServer.unregister <- client
			conn.Close()
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// 解析消息
			var msg WSMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			// 处理聊天消息
			if msg.Type == "chat" {
				data, ok := msg.Data.(map[string]interface{})
				if !ok {
					continue
				}

				conversationID, _ := strconv.ParseUint(data["conversation_id"].(string), 10, 32)
				content := data["content"].(string)
				msgType := "text"
				if t, ok := data["type"].(string); ok {
					msgType = t
				}

				// 保存消息到数据库
				newMsg := models.Message{
					ConversationID: uint(conversationID),
					SenderID:       userID,
					Content:        content,
					Type:           msgType,
					Status:         1,
					CreatedAt:      time.Now(),
				}
				database.DB.Create(&newMsg)

				// 更新会话最后消息
				database.DB.Model(&models.Conversation{}).Where("id = ?", conversationID).Updates(map[string]interface{}{
					"last_msg_id":   newMsg.ID,
					"last_msg_time": time.Now(),
				})

				// 发送给对方
				var conv models.Conversation
				database.DB.First(&conv, conversationID)
				participants := parseParticipants(conv.Participants)
				for _, pid := range participants {
					if pid != userID {
						if c, ok := wsServer.clients[pid]; ok {
							resp := WSMessage{
								Type: "chat",
								Data: map[string]interface{}{
									"conversation_id": conversationID,
									"message_id":      newMsg.ID,
									"sender_id":       userID,
									"content":         content,
									"type":            msgType,
									"created_at":      newMsg.CreatedAt,
								},
							}
							data, _ := json.Marshal(resp)
							c.Send <- data
						}
					}
				}
			}
		}
	}()

	// 发送消息
	go func() {
		defer conn.Close()
		for {
			select {
			case message, ok := <-client.Send:
				if !ok {
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}()
}

// SendToUser 发送消息给指定用户
func (s *WSServer) SendToUser(userID uint, message []byte) {
	s.mu.RLock()
	if client, ok := s.clients[userID]; ok {
		client.Send <- message
	}
	s.mu.RUnlock()
}

// parseParticipants 解析参与者列表
func parseParticipants(s string) []uint {
	var ids []uint
	for _, idStr := range splitString(s, ",") {
		id, _ := strconv.ParseUint(idStr, 10, 32)
		if id > 0 {
			ids = append(ids, uint(id))
		}
	}
	return ids
}

func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}