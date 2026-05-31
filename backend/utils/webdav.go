package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var webdavURL string
var webdavUsername string
var webdavPassword string

var httpClient = &http.Client{
	Timeout: 300 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
		WriteBufferSize:     32 * 1024,
		ReadBufferSize:      32 * 1024,
	},
}

func InitWebDAV(url, username, password string) {
	webdavURL = strings.TrimSuffix(url, "/")
	webdavUsername = username
	webdavPassword = password
}

func UploadToWebDAV(file *multipart.FileHeader, remotePath string) error {
	Info("开始上传文件到 WebDAV: 路径=%s, 大小=%d bytes", remotePath, file.Size)

	src, err := file.Open()
	if err != nil {
		Info("打开文件失败: %v", err)
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	Info("使用流式上传: %s", remotePath)
	return uploadLargeFile(src, file.Size, remotePath)
}

func uploadLargeFile(src multipart.File, size int64, remotePath string) error {
	req, err := http.NewRequest("PUT", webdavURL+remotePath, src)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.SetBasicAuth(webdavUsername, webdavPassword)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = size

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("连接WebDAV失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("WebDAV上传失败: %d - %s", resp.StatusCode, string(body))
	}

	Info("WebDAV大文件上传成功: %s, 大小: %d bytes", remotePath, size)
	return nil
}

func ProxyWebDAVHandler(c *gin.Context) {
	filePath := c.Param("path")
	targetURL := webdavURL + filePath

	switch c.Request.Method {
	case "GET", "HEAD":
		// 下载使用代理转发（支持视频拖拽、断点续传）
		proxyDownload(c, targetURL)
	case "PUT":
		// 上传使用重定向（避免经过服务器，提高性能）
		redirectURL := buildURLWithAuth(targetURL)
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	case "DELETE":
		proxyDelete(c, targetURL)
	case "PROPFIND":
		proxyPropfind(c, targetURL)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "不支持的操作"})
	}
}

// proxyDownload 处理 GET 和 HEAD 请求（代理模式）
func proxyDownload(c *gin.Context, targetURL string) {
	// 创建请求
	req, err := http.NewRequest(c.Request.Method, targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}

	// 添加认证
	req.SetBasicAuth(webdavUsername, webdavPassword)

	// 复制关键请求头（支持断点续传和缓存）
	if rangeHeader := c.Request.Header.Get("Range"); rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}
	if ifMatch := c.Request.Header.Get("If-Match"); ifMatch != "" {
		req.Header.Set("If-Match", ifMatch)
	}
	if ifNoneMatch := c.Request.Header.Get("If-None-Match"); ifNoneMatch != "" {
		req.Header.Set("If-None-Match", ifNoneMatch)
	}
	if ifModifiedSince := c.Request.Header.Get("If-Modified-Since"); ifModifiedSince != "" {
		req.Header.Set("If-Modified-Since", ifModifiedSince)
	}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("连接WebDAV失败: %v", err)})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 设置状态码
	c.Status(resp.StatusCode)

	// 对于 HEAD 请求，不需要传输 body
	if c.Request.Method == "HEAD" {
		return
	}

	// 流式传输 body
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		// 客户端可能主动断开连接，不记录为错误
		if err != io.EOF && !strings.Contains(err.Error(), "broken pipe") {
			Info("下载转发失败: %v", err)
		}
	}
}

// buildURLWithAuth 将认证信息编码到 URL 中（用于上传重定向）
func buildURLWithAuth(rawURL string) string {
	if webdavUsername == "" || webdavPassword == "" {
		return rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	parsedURL.User = url.UserPassword(webdavUsername, webdavPassword)
	return parsedURL.String()
}

// proxyDelete 处理 DELETE 请求
func proxyDelete(c *gin.Context, targetURL string) {
	req, err := http.NewRequest("DELETE", targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}
	req.SetBasicAuth(webdavUsername, webdavPassword)

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("删除失败: %v", err)})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("删除失败: %d - %s", resp.StatusCode, string(body))})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// proxyPropfind 处理 PROPFIND 请求（获取目录列表）
func proxyPropfind(c *gin.Context, targetURL string) {
	req, err := http.NewRequest("PROPFIND", targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}

	req.SetBasicAuth(webdavUsername, webdavPassword)
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Depth", c.Request.Header.Get("Depth"))

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("获取目录失败: %v", err)})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	c.Status(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		Info("PROPFIND 转发失败: %v", err)
	}
}

// GetFileURL 获取文件访问 URL
func GetFileURL(remotePath string) string {
	return "/proxy/webdav" + remotePath
}
