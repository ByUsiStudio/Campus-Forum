package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket连接管理器
var clients = make(map[uint]*websocket.Conn)
var clientMutex = make(chan struct{}, 1)

// ChatMessageRequest 聊天消息请求
type ChatMessageRequest struct {
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

// ChatMessageResponse 聊天消息响应
type ChatMessageResponse struct {
	ID            uint      `json:"id"`
	SenderID      uint      `json:"sender_id"`
	ReceiverID    uint      `json:"receiver_id"`
	Content       string    `json:"content"`
	IsRead        bool      `json:"is_read"`
	CreatedAt     time.Time `json:"created_at"`
	SenderName    string    `json:"sender_name"`
	TotalCount    int       `json:"total_count"`     // 消息总数校验
	LastMessageID uint      `json:"last_message_id"` // 最后一条消息ID
}

// GetChatMessages 获取聊天消息历史
func GetChatMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	otherUserID := c.GetUint("id")

	var messages []models.ChatMessage
	database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, otherUserID, otherUserID, userID).Order("created_at ASC").Find(&messages)

	// 标记消息为已读
	database.DB.Model(&models.ChatMessage{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", otherUserID, userID, false).
		Update("is_read", true)

	var response []ChatMessageResponse
	var lastMessageID uint
	if len(messages) > 0 {
		lastMessageID = messages[len(messages)-1].ID
	}

	for _, msg := range messages {
		var sender models.User
		database.DB.First(&sender, msg.SenderID)
		response = append(response, ChatMessageResponse{
			ID:            msg.ID,
			SenderID:      msg.SenderID,
			ReceiverID:    msg.ReceiverID,
			Content:       msg.Content,
			IsRead:        msg.IsRead,
			CreatedAt:     msg.CreatedAt,
			SenderName:    sender.DisplayName,
			TotalCount:    len(messages),
			LastMessageID: lastMessageID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"messages": response})
}

// GetChatSessions 获取聊天会话列表
func GetChatSessions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var sessions []models.ChatSession
	database.DB.Where("user1_id = ? OR user2_id = ?", userID, userID).Order("updated_at DESC").Find(&sessions)

	type SessionResponse struct {
		SessionID     uint        `json:"session_id"`
		OtherUser     models.User `json:"other_user"`
		UnreadCount   int         `json:"unread_count"`
		LastMessage   string      `json:"last_message"`
		LastMessageAt time.Time   `json:"last_message_at"`
	}

	var response []SessionResponse
	for _, session := range sessions {
		otherUserID := session.User1ID
		if session.User1ID == userID {
			otherUserID = session.User2ID
		}

		var otherUser models.User
		database.DB.First(&otherUser, otherUserID)

		var lastMsg models.ChatMessage
		if session.LastMessageID != nil {
			database.DB.First(&lastMsg, *session.LastMessageID)
		}

		response = append(response, SessionResponse{
			SessionID:     session.ID,
			OtherUser:     otherUser,
			UnreadCount:   session.UnreadCount,
			LastMessage:   lastMsg.Content,
			LastMessageAt: session.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"sessions": response})
}

// SendChatMessage 发送聊天消息（HTTP接口，用于不支持WebSocket的情况）
func SendChatMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req ChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ReceiverID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能给自己发消息"})
		return
	}

	// 检查是否可以发送消息
	if !canSendMessage(userID, req.ReceiverID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "未互相关注，只能发送2条消息"})
		return
	}

	message := models.ChatMessage{
		SenderID:   userID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送失败"})
		return
	}

	updateChatSession(userID, req.ReceiverID, message.ID)

	// 通过WebSocket发送消息
	sendMessageViaWebSocket(req.ReceiverID, message)

	c.JSON(http.StatusOK, gin.H{"message": "发送成功"})
}

// canSendMessage 检查是否可以发送消息
func canSendMessage(senderID, receiverID uint) bool {
	// 检查是否互相关注
	var follow1, follow2 models.Follow
	result1 := database.DB.Where("follower_id = ? AND following_id = ?", senderID, receiverID).First(&follow1)
	result2 := database.DB.Where("follower_id = ? AND following_id = ?", receiverID, senderID).First(&follow2)

	isMutual := result1.Error == nil && result2.Error == nil
	if isMutual {
		return true
	}

	// 未互相关注，检查已发送消息数量
	var count int64
	database.DB.Model(&models.ChatMessage{}).
		Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).
		Count(&count)

	return count < 2
}

// updateChatSession 更新聊天会话
func updateChatSession(userID1, userID2, messageID uint) {
	var session models.ChatSession
	result := database.DB.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		userID1, userID2, userID2, userID1).First(&session)

	if result.Error != nil {
		// 创建新会话
		session = models.ChatSession{
			User1ID:       userID1,
			User2ID:       userID2,
			LastMessageID: &messageID,
		}
		database.DB.Create(&session)
	} else {
		// 更新会话
		database.DB.Model(&session).Updates(models.ChatSession{
			LastMessageID: &messageID,
			UpdatedAt:     time.Now(),
		})
		// 增加对方的未读计数
		database.DB.Model(&session).UpdateColumn("unread_count", session.UnreadCount+1)
	}
}

// sendMessageViaWebSocket 通过WebSocket发送消息
func sendMessageViaWebSocket(userID uint, message models.ChatMessage) {
	clientMutex <- struct{}{}
	conn, exists := clients[userID]
	<-clientMutex

	if exists {
		var sender models.User
		database.DB.First(&sender, message.SenderID)

		// 获取消息总数
		var totalCount int64
		database.DB.Model(&models.ChatMessage{}).
			Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
				message.SenderID, userID, userID, message.SenderID).
			Count(&totalCount)

		response := ChatMessageResponse{
			ID:            message.ID,
			SenderID:      message.SenderID,
			ReceiverID:    message.ReceiverID,
			Content:       message.Content,
			IsRead:        message.IsRead,
			CreatedAt:     message.CreatedAt,
			SenderName:    sender.DisplayName,
			TotalCount:    int(totalCount),
			LastMessageID: message.ID,
		}

		data, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

// WebSocketHandler WebSocket连接处理
func WebSocketHandler(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	clientMutex <- struct{}{}
	clients[uint(userID)] = conn
	<-clientMutex

	defer func() {
		clientMutex <- struct{}{}
		delete(clients, uint(userID))
		<-clientMutex
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var req ChatMessageRequest
		if err := json.Unmarshal(message, &req); err != nil {
			continue
		}

		if !canSendMessage(uint(userID), req.ReceiverID) {
			response := map[string]string{"error": "未互相关注，只能发送2条消息"}
			data, _ := json.Marshal(response)
			conn.WriteMessage(websocket.TextMessage, data)
			continue
		}

		msg := models.ChatMessage{
			SenderID:   uint(userID),
			ReceiverID: req.ReceiverID,
			Content:    req.Content,
		}

		if err := database.DB.Create(&msg).Error; err != nil {
			continue
		}

		updateChatSession(uint(userID), req.ReceiverID, msg.ID)
		sendMessageViaWebSocket(req.ReceiverID, msg)

		var sender models.User
		database.DB.First(&sender, uint(userID))

		// 获取消息总数
		var totalCount int64
		database.DB.Model(&models.ChatMessage{}).
			Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
				msg.SenderID, req.ReceiverID, req.ReceiverID, msg.SenderID).
			Count(&totalCount)

		response := ChatMessageResponse{
			ID:            msg.ID,
			SenderID:      msg.SenderID,
			ReceiverID:    msg.ReceiverID,
			Content:       msg.Content,
			IsRead:        msg.IsRead,
			CreatedAt:     msg.CreatedAt,
			SenderName:    sender.DisplayName,
			TotalCount:    int(totalCount),
			LastMessageID: msg.ID,
		}

		data, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

// GetChatUnreadCount 获取未读消息数量
func GetChatUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var sessions []models.ChatSession
	database.DB.Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&sessions)

	totalUnread := 0
	for _, session := range sessions {
		// 只计算对方发过来的消息未读数量
		if session.User1ID == userID {
			// 用户是user1，未读消息来自user2
			var count int64
			database.DB.Model(&models.ChatMessage{}).
				Where("sender_id = ? AND receiver_id = ? AND is_read = ?", session.User2ID, userID, false).
				Count(&count)
			totalUnread += int(count)
		} else {
			var count int64
			database.DB.Model(&models.ChatMessage{}).
				Where("sender_id = ? AND receiver_id = ? AND is_read = ?", session.User1ID, userID, false).
				Count(&count)
			totalUnread += int(count)
		}
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": totalUnread})
}
