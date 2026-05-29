package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"
)

func GetSiteConfig(c *gin.Context) {
	var config models.SiteConfig
	result := database.DB.First(&config)
	if result.Error != nil {
		// 如果不存在，创建一个默认配置
		config = models.SiteConfig{
			SiteTitle: "校园论坛 - 分享与交流",
		}
		database.DB.Create(&config)
	}
	c.JSON(http.StatusOK, config)
}

func UpdateSiteConfig(c *gin.Context) {
	var input struct {
		SiteTitle     string `json:"site_title"`
		SMTPHost      string `json:"smtp_host"`
		SMTPPort      int    `json:"smtp_port"`
		SMTPUsername  string `json:"smtp_username"`
		SMTPPassword  string `json:"smtp_password"`
		SMTPFrom      string `json:"smtp_from"`
		SMTPFromName  string `json:"smtp_from_name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var config models.SiteConfig
	database.DB.First(&config)

	if input.SiteTitle != "" {
		config.SiteTitle = input.SiteTitle
	}
	if input.SMTPHost != "" {
		config.SMTPHost = input.SMTPHost
	}
	if input.SMTPPort != 0 {
		config.SMTPPort = input.SMTPPort
	}
	if input.SMTPUsername != "" {
		config.SMTPUsername = input.SMTPUsername
	}
	if input.SMTPPassword != "" {
		config.SMTPPassword = input.SMTPPassword
	}
	if input.SMTPFrom != "" {
		config.SMTPFrom = input.SMTPFrom
	}
	if input.SMTPFromName != "" {
		config.SMTPFromName = input.SMTPFromName
	}

	database.DB.Save(&config)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "config": config})
}

func TestSMTPConfig(c *gin.Context) {
	var input struct {
		SMTPHost     string `json:"smtp_host" binding:"required"`
		SMTPPort     int    `json:"smtp_port" binding:"required"`
		SMTPUsername string `json:"smtp_username" binding:"required"`
		SMTPPassword string `json:"smtp_password" binding:"required"`
		SMTPFrom     string `json:"smtp_from" binding:"required"`
		SMTPTo       string `json:"smtp_to" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建邮件内容
	subject := "SMTP配置测试"
	body := "SMTP配置测试成功！"
	msg := []byte("To: " + input.SMTPTo + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// 发送测试邮件
	auth := smtp.PlainAuth("", input.SMTPUsername, input.SMTPPassword, input.SMTPHost)
	err := smtp.SendMail(
		input.SMTPHost+":"+string(rune(input.SMTPPort+'0')),
		auth,
		input.SMTPFrom,
		[]string{input.SMTPTo},
		msg,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "邮件发送失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "邮件发送成功"})
}
