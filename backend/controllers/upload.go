package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"forum/database"
	"forum/models"
	"forum/utils"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// generateSecureFilename 生成安全的随机文件名
func generateSecureFilename(ext string) string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%s%s", hex.EncodeToString(bytes), ext)
}

// validateFileType 验证文件类型（通过magic number）
func validateFileType(file multipart.File, allowedTypes []string) (string, error) {
	// 读取文件头部512字节用于检测MIME类型
	header := make([]byte, 512)
	_, err := file.Read(header)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	
	// 重置文件指针
	file.Seek(0, 0)
	
	// 检测实际的文件类型
	contentType := http.DetectContentType(header)
	
	// 检查是否在允许的列表中
	for _, allowed := range allowedTypes {
		if strings.Contains(contentType, allowed) {
			return contentType, nil
		}
	}
	
	return "", fmt.Errorf("不支持的文件类型: %s", contentType)
}

func UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	// 检查文件大小（2MB）
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "头像大小不能超过2MB"})
		return
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持JPG、PNG、GIF格式的图片"})
		return
	}

	// 打开文件并验证实际类型
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
		return
	}
	defer src.Close()

	allowedMIMETypes := []string{"image/jpeg", "image/png", "image/gif"}
	if _, err := validateFileType(src, allowedMIMETypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成安全的文件名
	filename := generateSecureFilename(ext)
	remotePath := fmt.Sprintf("/avatars/%s", filename)

	// 上传到WebDAV
	if err := utils.UploadToWebDAV(file, remotePath); err != nil {
		utils.Error("文件上传失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
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
	file, err := c.FormFile("image")
	if err != nil {
		utils.Error("获取上传文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	utils.Info("开始上传图片: %s, 大小: %d", file.Filename, file.Size)

	// 检查文件大小 (10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过10MB"})
		return
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的图片格式，仅支持JPG、PNG、GIF、WebP"})
		return
	}

	// 打开文件并验证实际类型
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
		return
	}
	defer src.Close()

	allowedMIMETypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	if _, err := validateFileType(src, allowedMIMETypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成安全的文件名
	filename := generateSecureFilename(ext)
	remotePath := fmt.Sprintf("/images/%s", filename)

	// 上传到WebDAV
	if err := utils.UploadToWebDAV(file, remotePath); err != nil {
		utils.Error("图片上传失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	imageURL := "/proxy/webdav" + remotePath
	utils.Info("图片上传成功: %s -> %s", file.Filename, imageURL)

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"url":     imageURL,
	})
}

func UploadVideo(c *gin.Context) {
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

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".mp4": true, ".webm": true, ".ogg": true, ".mov": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的视频格式，仅支持MP4、WebM、OGG、MOV"})
		return
	}

	// 打开文件并验证实际类型
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
		return
	}
	defer src.Close()

	allowedMIMETypes := []string{"video/mp4", "video/webm", "video/ogg", "video/quicktime"}
	if _, err := validateFileType(src, allowedMIMETypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成安全的文件名
	filename := generateSecureFilename(ext)
	remotePath := fmt.Sprintf("/videos/%s", filename)

	// 上传到WebDAV
	if err := utils.UploadToWebDAV(file, remotePath); err != nil {
		utils.Error("文件上传失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	videoURL := "/proxy/webdav" + remotePath

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"url":     videoURL,
	})
}

func UploadVoice(c *gin.Context) {
	file, err := c.FormFile("voice")
	if err != nil {
		utils.Error("获取上传文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	utils.Info("开始上传语音: %s, 大小: %d", file.Filename, file.Size)

	// 检查文件大小 (50MB)
	if file.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "语音文件大小不能超过50MB"})
		return
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	
	// 检查文件是否有扩展名
	if ext == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件没有扩展名，请确保文件格式正确"})
		return
	}
	
	allowedExts := map[string]bool{
		".mp3": true, ".webm": true, ".ogg": true, ".wav": true, ".m4a": true, ".aac": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的语音格式，支持的格式: MP3, WAV, AAC, OGG, M4A, WebM"})
		return
	}

	// 打开文件并验证实际类型
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
		return
	}
	defer src.Close()

	allowedMIMETypes := []string{"audio/mpeg", "audio/webm", "audio/ogg", "audio/wav", "audio/mp4", "audio/aac"}
	if _, err := validateFileType(src, allowedMIMETypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成安全的文件名
	filename := generateSecureFilename(ext)
	remotePath := fmt.Sprintf("/voices/%s", filename)

	// 上传到WebDAV
	if err := utils.UploadToWebDAV(file, remotePath); err != nil {
		utils.Error("语音上传失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	voiceURL := "/proxy/webdav" + remotePath
	utils.Info("语音上传成功: %s -> %s", file.Filename, voiceURL)

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"url":     voiceURL,
	})
}