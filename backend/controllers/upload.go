package controllers

import (
    "fmt"
    "forum/database"
    "forum/models"
    "forum/utils"
    "net/http"
    "path/filepath"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

func UploadAvatar(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    file, err := c.FormFile("avatar")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
        return
    }
    
    // 检查文件类型
    ext := strings.ToLower(filepath.Ext(file.Filename))
    if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "只支持图片文件"})
        return
    }
    
    // 生成文件名
    filename := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().Unix(), ext)
    remotePath := fmt.Sprintf("/avatars/%s", filename)
    
    // 上传到WebDAV
    if err := utils.UploadToWebDAV(file, remotePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
        return
    }
    
    // 更新用户头像
    avatarURL := "/proxy/webdav" + remotePath
    database.DB.Model(&models.User{}).Where("id = ?", userID).Update("avatar", avatarURL)
    
    c.JSON(http.StatusOK, gin.H{
        "message": "上传成功",
        "url":     avatarURL,
    })
}

func UploadImage(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
        return
    }
    
    // 检查文件大小 (10MB)
    if file.Size > 10*1024*1024 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过10MB"})
        return
    }
    
    // 检查文件类型
    ext := strings.ToLower(filepath.Ext(file.Filename))
    allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
    allowed := false
    for _, allowedExt := range allowedExts {
        if ext == allowedExt {
            allowed = true
            break
        }
    }
    if !allowed {
        c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的图片格式"})
        return
    }
    
    // 生成文件名
    filename := fmt.Sprintf("image_%d_%d%s", userID, time.Now().UnixNano(), ext)
    remotePath := fmt.Sprintf("/images/%s", filename)
    
    // 上传到WebDAV
    if err := utils.UploadToWebDAV(file, remotePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
        return
    }
    
    imageURL := "/proxy/webdav" + remotePath
    
    c.JSON(http.StatusOK, gin.H{
        "message": "上传成功",
        "url":     imageURL,
    })
}

func UploadVideo(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    file, err := c.FormFile("video")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
        return
    }
    
    // 检查文件大小 (100MB)
    if file.Size > 100*1024*1024 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "视频大小不能超过100MB"})
        return
    }
    
    // 检查文件类型
    ext := strings.ToLower(filepath.Ext(file.Filename))
    allowedExts := []string{".mp4", ".webm", ".ogg", ".mov"}
    allowed := false
    for _, allowedExt := range allowedExts {
        if ext == allowedExt {
            allowed = true
            break
        }
    }
    if !allowed {
        c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的视频格式"})
        return
    }
    
    // 生成文件名
    filename := fmt.Sprintf("video_%d_%d%s", userID, time.Now().UnixNano(), ext)
    remotePath := fmt.Sprintf("/videos/%s", filename)
    
    // 上传到WebDAV
    if err := utils.UploadToWebDAV(file, remotePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
        return
    }
    
    videoURL := "/proxy/webdav" + remotePath
    
    c.JSON(http.StatusOK, gin.H{
        "message": "上传成功",
        "url":     videoURL,
    })
}