package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 简单的基于IP的速率限制器
type RateLimiter struct {
	clients map[string]*ClientInfo
	mu      sync.Mutex
	limit   int
	window  time.Duration
}

// ClientInfo 客户端请求信息
type ClientInfo struct {
	count    int
	resetAt  time.Time
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*ClientInfo),
		limit:   limit,
		window:  window,
	}

	// 启动清理协程
	go rl.cleanup()

	return rl
}

// cleanup 定期清理过期的客户端记录
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, client := range rl.clients {
			if now.After(client.resetAt) {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[ip]

	if !exists || now.After(client.resetAt) {
		// 新客户端或窗口已重置
		rl.clients[ip] = &ClientInfo{
			count:   1,
			resetAt: now.Add(rl.window),
		}
		return true
	}

	client.count++
	if client.count > rl.limit {
		return false
	}

	return true
}

// RateLimit 速率限制中间件
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PerEndpointRateLimit 针对特定端点的速率限制
func PerEndpointRateLimit(limits map[string]struct {
	Limit  int
	Window time.Duration
}) gin.HandlerFunc {
	limiters := make(map[string]*RateLimiter)
	for path, config := range limits {
		limiters[path] = NewRateLimiter(config.Limit, config.Window)
	}

	return func(c *gin.Context) {
		path := c.FullPath()
		if limiter, exists := limiters[path]; exists {
			ip := c.ClientIP()
			if !limiter.Allow(ip) {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "请求过于频繁，请稍后再试",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
