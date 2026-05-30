package utils

import (
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

	resp, err := httpClient.Do(req)
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
	c.Header("Content-Length", fmt.Sprintf("%d", resp.ContentLength))

	c.Status(http.StatusOK)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		Info("下载转发失败: %v", err)
	}
}

func proxyPUT(c *gin.Context, targetURL string) {
	req, err := http.NewRequest("PUT", targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}

	req.SetBasicAuth(webdavUsername, webdavPassword)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = c.Request.ContentLength

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("上传失败: %d - %s", resp.StatusCode, string(body))})
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

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func proxyPROPFIND(c *gin.Context, targetURL string) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取目录失败"})
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/xml"
	}

	c.Header("Content-Type", contentType)
	c.Status(http.StatusOK)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		Info("PROPFIND 转发失败: %v", err)
	}
}

func GetFileURL(remotePath string) string {
	return "/proxy/webdav" + remotePath
}