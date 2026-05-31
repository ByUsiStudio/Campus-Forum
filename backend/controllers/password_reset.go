package controllers

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"forum/database"
	"forum/models"
	"math/big"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
)

func generateCode(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code
}

func sendEmail(smtpHost string, smtpPort int, username, password, from, to, subject, body string) error {
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// 构建带 UTF-8 编码的邮件内容
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: =?UTF-8?B?%s?=\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, to, base64.StdEncoding.EncodeToString([]byte(subject)), body))

	// 尝试使用 TLS 连接
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpHost, smtpPort), &tls.Config{
		ServerName:         smtpHost,
		InsecureSkipVerify: true, // 允许自签名证书
	})
	if err != nil {
		// TLS 连接失败，尝试普通连接
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
	err = w.Close()
	return err
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
	database.DB.First(&siteConfig, 1)

	if siteConfig.SMTPHost == "" || siteConfig.SMTPUsername == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "邮件服务未配置"})
		return
	}

	code := generateCode(6)

	expiry := time.Now().Add(15 * time.Minute)
	user.ResetToken = code
	user.ResetExpiry = &expiry
	database.DB.Save(&user)

	email := input.QQNumber + "@qq.com"
	subject := "密码重置验证码"
	body := fmt.Sprintf("您的密码重置验证码是：%s，15分钟内有效。\n\n如果不是您本人操作，请忽略此邮件。", code)

	err := sendEmail(
		siteConfig.SMTPHost,
		siteConfig.SMTPPort,
		siteConfig.SMTPUsername,
		siteConfig.SMTPPassword,
		siteConfig.SMTPFrom,
		email,
		subject,
		body,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送邮件失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送到您的QQ邮箱"})
}

func ResetPassword(c *gin.Context) {
	var input struct {
		QQNumber string `json:"qq_number" binding:"required"`
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required;min=6"`
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

	if user.ResetToken != input.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	if user.ResetExpiry == nil || time.Now().After(*user.ResetExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码已过期"})
		return
	}

	user.Password = database.HashPassword(input.Password)
	user.ResetToken = ""
	user.ResetExpiry = nil
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}
