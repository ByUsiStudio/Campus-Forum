package utils

import (
	"bytes"
	"fmt"
	"forum/utils"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

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
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	req, err := http.NewRequest("PUT", webdavURL+remotePath, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.SetBasicAuth(webdavUsername, webdavPassword)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("连接WebDAV失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("WebDAV上传失败: %d - %s", resp.StatusCode, string(body))
	}

	utils.Info("WebDAV上传成功: %s, 状态码: %d", remotePath, resp.StatusCode)

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
