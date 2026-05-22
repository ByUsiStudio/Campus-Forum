package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var webdavURL string
var webdavUsername string
var webdavPassword string

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

	// 判断文件大小，小文件使用内存上传，大文件使用流式上传
	// 超过1MB的文件使用流式上传
	const streamThreshold = 1 * 1024 * 1024

	if file.Size > streamThreshold {
		Info("使用流式上传: %s", remotePath)
		return uploadLargeFile(src, file.Size, remotePath)
	}

	Info("使用内存上传: %s", remotePath)

	// 小文件：读取到内存后上传
	data, err := io.ReadAll(src)
	if err != nil {
		Info("读取文件失败: %v", err)
		return fmt.Errorf("读取文件失败: %w", err)
	}

	Info("文件读取完成，准备上传到: %s", webdavURL+remotePath)

	req, err := http.NewRequest("PUT", webdavURL+remotePath, bytes.NewReader(data))
	if err != nil {
		Info("创建请求失败: %v", err)
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.SetBasicAuth(webdavUsername, webdavPassword)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

	// 设置超时：连接10秒，读取60秒，写入300秒（适用于大文件）
	client := &http.Client{
		Timeout: 300 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		Info("连接WebDAV失败: %v", err)
		return fmt.Errorf("连接WebDAV失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		Info("WebDAV返回错误: 状态码=%d, 响应=%s", resp.StatusCode, string(body))
		return fmt.Errorf("WebDAV上传失败: %d - %s", resp.StatusCode, string(body))
	}

	Info("WebDAV上传成功: %s, 大小: %d bytes, 状态码: %d", remotePath, file.Size, resp.StatusCode)

	return nil
}

// uploadLargeFile 流式上传大文件（超过10MB）
func uploadLargeFile(src multipart.File, size int64, remotePath string) error {
	req, err := http.NewRequest("PUT", webdavURL+remotePath, src)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.SetBasicAuth(webdavUsername, webdavPassword)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = size

	// 设置超时：连接10秒，传输超时300秒
	client := &http.Client{
		Timeout: 300 * time.Second,
	}
	resp, err := client.Do(req)
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
		proxyGET(c, targetURL)
	case "PUT":
		proxyPUT(c, targetURL)
	case "DELETE":
		proxyDELETE(c, targetURL)
	case "PROPFIND":
		proxyPROPFIND(c, targetURL)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "不支持的操作"})
	}
}

func proxyGET(c *gin.Context, targetURL string) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}
	req.SetBasicAuth(webdavUsername, webdavPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "文件不存在"})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Header("Content-Type", contentType)
	c.DataFromReader(http.StatusOK, resp.ContentLength, contentType, resp.Body, nil)
}

func proxyPUT(c *gin.Context, targetURL string) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取请求失败"})
		return
	}

	req, err := http.NewRequest("PUT", targetURL, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}
	req.SetBasicAuth(webdavUsername, webdavPassword)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("上传失败: %d", resp.StatusCode)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "上传成功"})
}

func proxyDELETE(c *gin.Context, targetURL string) {
	req, err := http.NewRequest("DELETE", targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}
	req.SetBasicAuth(webdavUsername, webdavPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func proxyPROPFIND(c *gin.Context, targetURL string) {
	body, _ := io.ReadAll(c.Request.Body)

	req, err := http.NewRequest("PROPFIND", targetURL, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}
	req.SetBasicAuth(webdavUsername, webdavPassword)
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Depth", c.Request.Header.Get("Depth"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取目录失败"})
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/xml"
	}
	c.Header("Content-Type", contentType)
	c.DataFromReader(http.StatusOK, resp.ContentLength, contentType, resp.Body, nil)
}

func GetFileURL(remotePath string) string {
	return "/proxy/webdav" + remotePath
}
