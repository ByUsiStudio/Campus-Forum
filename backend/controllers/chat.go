package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/models"
	"forum/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetConversations 获取会话列表
func GetConversations(c *gin.Context) {
	userID := c.GetUint("user_id")

	var userConvs []models.UserConversation
	database.DB.Where("user_id = ?", userID).Preload("Conversation").Find(&userConvs)

	var conversations []map[string]interface{}
	for _, uc := range userConvs {
		conv := uc.Conversation
		// 获取对方用户信息
		participants := strings.Split(conv.Participants, ",")
		var targetID uint
		for _, pid := range participants {
			id, _ := strconv.ParseUint(pid, 10, 32)
			if uint(id) != userID {
				targetID = uint(id)
				break
			}
		}

		var targetUser models.User
		database.DB.First(&targetUser, targetID)

		// 获取最后消息
		var lastMsg models.Message
		database.DB.First(&lastMsg, conv.LastMsgID)

		conversations = append(conversations, map[string]interface{}{
			"id":            conv.ID,
			"type":          conv.Type,
			"target_user":   targetUser,
			"unread_count":  uc.UnreadCount,
			"last_message":  lastMsg,
			"last_msg_time": conv.LastMsgTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

// GetMessages 获取消息历史
func GetMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	conversationID := c.Query("conversation_id")
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	// 验证用户是否属于该会话
	var userConv models.UserConversation
	result := database.DB.Where("user_id = ? AND conversation_id = ?", userID, conversationID).First(&userConv)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该会话"})
		return
	}

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	var messages []models.Message
	database.DB.Where("conversation_id = ? AND status = 1", conversationID).
		Order("created_at DESC").
		Limit(limitInt).
		Offset(offsetInt).
		Find(&messages)

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		ConversationID uint   `json:"conversation_id"`
		Content        string `json:"content"`
		Type           string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证用户是否属于该会话
	var userConv models.UserConversation
	result := database.DB.Where("user_id = ? AND conversation_id = ?", userID, req.ConversationID).First(&userConv)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该会话"})
		return
	}

	if req.Type == "" {
		req.Type = "text"
	}

	// 创建消息
	msg := models.Message{
		ConversationID: req.ConversationID,
		SenderID:       userID,
		Content:        req.Content,
		Type:           req.Type,
		Status:         1,
		CreatedAt:      time.Now(),
	}
	database.DB.Create(&msg)

	// 更新会话最后消息
	database.DB.Model(&models.Conversation{}).Where("id = ?", req.ConversationID).Updates(map[string]interface{}{
		"last_msg_id":   msg.ID,
		"last_msg_time": time.Now(),
	})

	// 更新对方未读数
	var conv models.Conversation
	database.DB.First(&conv, req.ConversationID)
	participants := strings.Split(conv.Participants, ",")
	for _, pid := range participants {
		id, _ := strconv.ParseUint(pid, 10, 32)
		if uint(id) != userID {
			database.DB.Model(&models.UserConversation{}).
				Where("user_id = ? AND conversation_id = ?", uint(id), req.ConversationID).
				Update("unread_count", database.DB.Raw("unread_count + 1"))
		}
	}

	// 通过WebSocket发送给对方
	wsServer := service.GetWSServer()
	if wsServer != nil {
		for _, pid := range participants {
			id, _ := strconv.ParseUint(pid, 10, 32)
			if uint(id) != userID {
				// 构造消息
				resp := map[string]interface{}{
					"type":            "chat",
					"conversation_id": req.ConversationID,
					"message_id":      msg.ID,
					"sender_id":       userID,
					"content":         req.Content,
					"msg_type":        req.Type,
					"created_at":      msg.CreatedAt,
				}
				data, _ := jsonMarshal(resp)
				wsServer.SendToUser(uint(id), data)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// CreatePrivateConversation 创建私聊会话
func CreatePrivateConversation(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		TargetUserID uint `json:"target_user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if userID == req.TargetUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能与自己创建会话"})
		return
	}

	// 检查是否已存在会话
	participants := strconv.FormatUint(uint64(userID), 10) + "," + strconv.FormatUint(uint64(req.TargetUserID), 10)
	participantsReverse := strconv.FormatUint(uint64(req.TargetUserID), 10) + "," + strconv.FormatUint(uint64(userID), 10)

	var existingConv models.Conversation
	result := database.DB.Where("type = 'private' AND (participants = ? OR participants = ?)", participants, participantsReverse).First(&existingConv)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{"conversation_id": existingConv.ID})
		return
	}

	// 创建新会话
	conv := models.Conversation{
		Type:         "private",
		Participants: participants,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	database.DB.Create(&conv)

	// 创建用户会话关联
	userConv1 := models.UserConversation{
		UserID:         userID,
		ConversationID: conv.ID,
		UnreadCount:    0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	userConv2 := models.UserConversation{
		UserID:         req.TargetUserID,
		ConversationID: conv.ID,
		UnreadCount:    0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	database.DB.Create(&userConv1)
	database.DB.Create(&userConv2)

	c.JSON(http.StatusOK, gin.H{"conversation_id": conv.ID})
}

// GetChatUnreadCount 获取聊天未读消息数
func GetChatUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var totalUnread int64
	database.DB.Model(&models.UserConversation{}).Where("user_id = ?", userID).Select("SUM(unread_count)").Scan(&totalUnread)

	c.JSON(http.StatusOK, gin.H{"unread_count": totalUnread})
}

// MarkConversationRead 标记会话已读
func MarkConversationRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	conversationID := c.Param("id")

	convID, _ := strconv.ParseUint(conversationID, 10, 32)

	result := database.DB.Model(&models.UserConversation{}).
		Where("user_id = ? AND conversation_id = ?", userID, uint(convID)).
		Updates(map[string]interface{}{
			"unread_count":   0,
			"last_read_time": time.Now(),
		})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已标记为已读"})
}

// HandleWebSocket 处理WebSocket连接
func HandleWebSocket(c *gin.Context) {
	userID := c.GetUint("user_id")
	service.HandleWebSocket(c.Writer, c.Request, userID)
}

func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
