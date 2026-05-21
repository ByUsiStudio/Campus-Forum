package utils

import (
    "bytes"
    "fmt"
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
        return err
    }
    defer src.Close()
    
    data, err := io.ReadAll(src)
    if err != nil {
        return err
    }
    
    req, err := http.NewRequest("PUT", webdavURL+remotePath, bytes.NewReader(data))
    if err != nil {
        return err
    }
    
    req.SetBasicAuth(webdavUsername, webdavPassword)
    req.Header.Set("Content-Type", "application/octet-stream")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("webdav上传失败: %d", resp.StatusCode)
    }
    
    return nil
}

func ProxyWebDAVHandler(c *gin.Context) {
    filePath := c.Param("path")
    if filePath == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效路径"})
        return
    }
    
    targetURL := webdavURL + filePath
    
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
        c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
        return
    }
    
    // 设置响应头
    contentType := resp.Header.Get("Content-Type")
    if contentType == "" {
        contentType = "application/octet-stream"
    }
    c.Header("Content-Type", contentType)
    
    // 代理内容
    c.DataFromReader(http.StatusOK, resp.ContentLength, contentType, resp.Body, nil)
}

func GetFileURL(remotePath string) string {
    return "/proxy/webdav" + remotePath
}