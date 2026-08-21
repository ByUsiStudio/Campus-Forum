package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"
	"net/smtp"
	"strconv"

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
	
	// 隐藏敏感信息，不返回SMTP密码
	response := gin.H{
		"site_title":            config.SiteTitle,
		"icp_number":            config.ICPNumber,
		"public_security_number": config.PublicSecurityNumber,
		"smtp_host":             config.SMTPHost,
		"smtp_port":             config.SMTPPort,
		"smtp_username":         config.SMTPUsername,
		"smtp_from":             config.SMTPFrom,
		"smtp_from_name":        config.SMTPFromName,
		"smtp_password_set":     config.SMTPPassword != "", // 仅返回是否已设置密码的标记
	}
	c.JSON(http.StatusOK, response)
}

func UpdateSiteConfig(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var config models.SiteConfig
	database.DB.First(&config)

	// 只有在提供了值时才更新（包括空值表示删除）
	if val, exists := updates["site_title"]; exists {
		if str, ok := val.(string); ok {
			config.SiteTitle = str
		}
	}
	if val, exists := updates["icp_number"]; exists {
		if str, ok := val.(string); ok {
			config.ICPNumber = str
		}
	}
	if val, exists := updates["public_security_number"]; exists {
		if str, ok := val.(string); ok {
			config.PublicSecurityNumber = str
		}
	}
	if val, exists := updates["smtp_host"]; exists {
		if str, ok := val.(string); ok {
			config.SMTPHost = str
		}
	}
	if val, exists := updates["smtp_port"]; exists {
		if num, ok := val.(float64); ok {
			config.SMTPPort = int(num)
		}
	}
	if val, exists := updates["smtp_username"]; exists {
		if str, ok := val.(string); ok {
			config.SMTPUsername = str
		}
	}
	if val, exists := updates["smtp_password"]; exists {
		if str, ok := val.(string); ok {
			// 只有在密码不为空时才更新（空字符串表示不修改密码）
			if str != "" {
				config.SMTPPassword = str
			}
		}
	}
	if val, exists := updates["smtp_from"]; exists {
		if str, ok := val.(string); ok {
			config.SMTPFrom = str
		}
	}
	if val, exists := updates["smtp_from_name"]; exists {
		if str, ok := val.(string); ok {
			config.SMTPFromName = str
		}
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
		input.SMTPHost+":"+strconv.Itoa(input.SMTPPort),
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
