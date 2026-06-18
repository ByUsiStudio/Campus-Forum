package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"net/http"
	_ "strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

func Register(c *gin.Context) {
	var input struct {
		Username    string `json:"username" binding:"required"`
		QQNumber    string `json:"qq_number" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
		Password    string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否存在
	var existingUser models.User
	if result := database.DB.Where("username = ? OR qq_number = ?", input.Username, input.QQNumber).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或QQ号已存在"})
		return
	}

	// 使用QQ头像作为默认头像
	avatarURL := "https://q1.qlogo.cn/g?b=qq&nk=" + input.QQNumber + "&s=640"

	hashedPassword, err := database.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		Username:    input.Username,
		QQNumber:    input.QQNumber,
		DisplayName: input.DisplayName,
		Password:    hashedPassword,
		Avatar:      avatarURL,
		Role:        "user",
	}

	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := database.DB.Where("username = ?", input.Username).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查用户是否被封禁
	if user.Status == "banned" {
		c.JSON(http.StatusForbidden, gin.H{"error": "该账号已被封禁，原因：" + user.BanReason})
		return
	}

	if !database.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成访问令牌（有效期1小时）
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"type":    "access",
	})

	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 生成刷新令牌（有效期7天）
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"type":    "refresh",
	})

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成刷新令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         accessTokenString,
		"refresh_token": refreshTokenString,
		"expires_in":    3600,
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"display_name": user.DisplayName,
			"avatar":       user.Avatar,
			"role":         user.Role,
			"status":       user.Status,
		},
	})
}

// RefreshToken 刷新访问令牌
func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析刷新令牌
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "刷新令牌无效"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌解析失败"})
		return
	}

	// 检查令牌类型
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌类型错误"})
		return
	}

	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	// 检查用户是否存在且状态正常
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	if user.Status == "banned" {
		c.JSON(http.StatusForbidden, gin.H{"error": "该账号已被封禁"})
		return
	}

	// 生成新的访问令牌
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"type":    "access",
	})

	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 生成新的刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"type":    "refresh",
	})

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成刷新令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         accessTokenString,
		"refresh_token": refreshTokenString,
		"expires_in":    3600,
	})
}

func InitAdmin(c *gin.Context) {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统已初始化"})
		return
	}

	var input struct {
		Username    string `json:"username" binding:"required"`
		QQNumber    string `json:"qq_number" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
		Password    string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	avatarURL := "https://q1.qlogo.cn/g?b=qq&nk=" + input.QQNumber + "&s=640"

	hashedPassword, err := database.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	admin := models.User{
		Username:    input.Username,
		QQNumber:    input.QQNumber,
		DisplayName: input.DisplayName,
		Password:    hashedPassword,
		Avatar:      avatarURL,
		Role:        "admin",
	}

	if result := database.DB.Create(&admin); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "初始化失败"})
		return
	}

	// 创建默认分区
	categories := []models.Category{
		{Name: "表白墙", Description: "不止局限于表白墙", SortOrder: 1},
	}
	database.DB.Create(&categories)

	c.JSON(http.StatusOK, gin.H{"message": "初始化成功"})
}

func CheckInit(c *gin.Context) {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"initialized": count > 0})
}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"qq_number":    user.QQNumber,
		"display_name": user.DisplayName,
		"avatar":       user.Avatar,
		"role":         user.Role,
		"signature":    user.Signature,
		"status":       user.Status,
		"ban_reason":   user.BanReason,
		"created_at":   user.CreatedAt,
	})
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User

	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var input struct {
		DisplayName string `json:"display_name"`
		Avatar      string `json:"avatar"`
		Signature   string `json:"signature"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.DisplayName != "" {
		user.DisplayName = input.DisplayName
	}
	if input.Avatar != "" {
		user.Avatar = input.Avatar
	}
	if input.Signature != "" {
		// XSS过滤：清理用户签名
		user.Signature = utils.SanitizeHTML(input.Signature)
	}

	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
