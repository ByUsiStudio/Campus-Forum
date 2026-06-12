package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// IMClient IM服务客户端
type IMClient struct {
	ApiURL     string
	WsURL      string
	AdminURL   string
	HTTPClient *http.Client
}

var imClient *IMClient

// InitIMClient 初始化IM客户端
func InitIMClient(apiURL, wsURL, adminURL string) {
	imClient = &IMClient{
		ApiURL:     apiURL,
		WsURL:      wsURL,
		AdminURL:   adminURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetIMClient 获取IM客户端实例
func GetIMClient() *IMClient {
	return imClient
}

// IMRequest IM请求结构
type IMRequest struct {
	UserID      string                 `json:"user_id"`
	TargetID    string                 `json:"target_id,omitempty"`
	ConversationID string              `json:"conversation_id,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// IMResponse IM响应结构
type IMResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// callAPI 调用IM API
func (c *IMClient) callAPI(endpoint string, method string, req interface{}) (*IMResponse, error) {
	url := c.ApiURL + endpoint

	var body io.Reader
	if req != nil {
		jsonData, err := json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("序列化请求失败: %v", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	httpReq, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var imResp IMResponse
	if err := json.Unmarshal(respBody, &imResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &imResp, nil
}

// SendMessage 发送消息
func (c *IMClient) SendMessage(userID, conversationID, content, msgType string) (*IMResponse, error) {
	req := IMRequest{
		UserID:         userID,
		ConversationID: conversationID,
		Message:        content,
		Type:           msgType,
	}
	return c.callAPI("/message/send", "POST", req)
}

// GetConversations 获取会话列表
func (c *IMClient) GetConversations(userID string) (*IMResponse, error) {
	req := IMRequest{
		UserID: userID,
	}
	return c.callAPI("/conversations/list", "GET", req)
}

// GetMessages 获取消息历史
func (c *IMClient) GetMessages(userID, conversationID string, limit, offset int) (*IMResponse, error) {
	req := IMRequest{
		UserID:         userID,
		ConversationID: conversationID,
		Data: map[string]interface{}{
			"limit":  limit,
			"offset": offset,
		},
	}
	return c.callAPI("/messages/history", "GET", req)
}

// CreatePrivateConversation 创建私聊会话
func (c *IMClient) CreatePrivateConversation(userID, targetID string) (*IMResponse, error) {
	req := IMRequest{
		UserID:   userID,
		TargetID: targetID,
	}
	return c.callAPI("/conversation/private/create", "POST", req)
}

// GetUnreadCount 获取未读消息数
func (c *IMClient) GetUnreadCount(userID string) (*IMResponse, error) {
	req := IMRequest{
		UserID: userID,
	}
	return c.callAPI("/messages/unread/count", "GET", req)
}

// MarkConversationRead 标记会话已读
func (c *IMClient) MarkConversationRead(userID, conversationID string) (*IMResponse, error) {
	req := IMRequest{
		UserID:         userID,
		ConversationID: conversationID,
	}
	return c.callAPI("/conversation/read", "POST", req)
}

// SyncUser 同步用户信息到IM系统
func (c *IMClient) SyncUser(userID, username, avatar string) (*IMResponse, error) {
	req := IMRequest{
		UserID: userID,
		Data: map[string]interface{}{
			"username": username,
			"avatar":   avatar,
		},
	}
	return c.callAPI("/user/sync", "POST", req)
}