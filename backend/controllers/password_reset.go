package controllers

import (
	"crypto/rand"
	"crypto/tls"
	_ "embed"
	"encoding/base64"
	"fmt"
	"forum/database"
	"forum/models"
	"math/big"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed templates/reset_password.html
var resetPasswordTemplate string

func generateCode(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code
}

func generateIdentifier() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func renderEmailTemplate(subject, code string) string {
	htmlBody := strings.ReplaceAll(resetPasswordTemplate, "{{.Subject}}", subject)
	htmlBody = strings.ReplaceAll(htmlBody, "{{.Code}}", code)
	return htmlBody
}

func sendEmail(smtpHost string, smtpPort int, username, password, from, to, subject, code string) error {
	auth := smtp.PlainAuth("", username, password, smtpHost)

	htmlBody := renderEmailTemplate(subject, code)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: =?UTF-8?B?%s?=\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, to, base64.StdEncoding.EncodeToString([]byte(subject)), htmlBody))

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpHost, smtpPort), &tls.Config{
		ServerName:         smtpHost,
		InsecureSkipVerify: true,
	})
	if err != nil {
		err = smtp.SendMail(
			fmt.Sprintf("%s:%d", smtpHost, smtpPort),
			auth,
			from,
			[]string{to},
			msg,
		)
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	return w.Close()
}

func SendResetCode(c *gin.Context) {
	var input struct {
		QQNumber string `json:"qq_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供QQ号码"})
		return
	}

	var user models.User
	if result := database.DB.Where("qq_number = ?", input.QQNumber).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "该QQ号码未注册"})
		return
	}

	var siteConfig models.SiteConfig
	result := database.DB.First(&siteConfig)
	if result.Error != nil {
		siteConfig = models.SiteConfig{SiteTitle: "校园论坛 - 分享与交流"}
		database.DB.Create(&siteConfig)
	}

	if siteConfig.SMTPHost == "" || siteConfig.SMTPUsername == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "邮件服务未配置"})
		return
	}

	code := generateCode(6)
	identifier := generateIdentifier()

	expiry := time.Now().Add(15 * time.Minute)
	user.ResetToken = code
	user.ResetIdentifier = identifier
	user.ResetExpiry = &expiry
	database.DB.Save(&user)

	email := input.QQNumber + "@qq.com"
	subject := "校园论坛 - 密码重置验证码"

	err := sendEmail(
		siteConfig.SMTPHost,
		siteConfig.SMTPPort,
		siteConfig.SMTPUsername,
		siteConfig.SMTPPassword,
		siteConfig.SMTPFrom,
		email,
		subject,
		code,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送邮件失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "验证码已发送到您的QQ邮箱",
		"identifier": identifier,
	})
}

func ResetPassword(c *gin.Context) {
	var input struct {
		QQNumber   string `json:"qq_number" binding:"required"`
		Code       string `json:"code" binding:"required"`
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required;min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供完整的参数"})
		return
	}

	var user models.User
	if result := database.DB.Where("qq_number = ?", input.QQNumber).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "该QQ号码未注册"})
		return
	}

	if user.ResetIdentifier != input.Identifier {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标识token"})
		return
	}

	if user.ResetToken != input.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	if user.ResetExpiry == nil || time.Now().After(*user.ResetExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码已过期"})
		return
	}

	hashedPassword, err := database.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	user.Password = hashedPassword
	user.ResetToken = ""
	user.ResetIdentifier = ""
	user.ResetExpiry = nil
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}
