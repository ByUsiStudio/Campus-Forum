package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

type AuthService struct{}

var Auth = &AuthService{}

func (s *AuthService) Register(username, qqNumber, displayName, password string) error {
	var existingUser models.User
	result := database.DB.Where("username = ? OR qq_number = ?", username, qqNumber).First(&existingUser)
	if result.Error == nil {
		return utils.NewError("用户名或QQ号已存在", 400)
	}

	avatarURL := "https://q1.qlogo.cn/g?b=qq&nk=" + qqNumber + "&s=640"
	hashedPassword, err := database.HashPassword(password)
	if err != nil {
		return utils.NewError("密码加密失败", 500)
	}

	user := models.User{
		Username:    username,
		QQNumber:    qqNumber,
		DisplayName: displayName,
		Password:    hashedPassword,
		Avatar:      avatarURL,
		Role:        "user",
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return utils.NewError("注册失败", 500)
	}
	return nil
}

func (s *AuthService) Login(username, password string) (*models.User, string, string, error) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, "", "", utils.NewError("用户名或密码错误", 401)
	}

	if user.Status == "banned" {
		return nil, "", "", utils.NewError("该账号已被封禁，原因："+user.BanReason, 403)
	}

	if !database.CheckPasswordHash(password, user.Password) {
		return nil, "", "", utils.NewError("用户名或密码错误", 401)
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, "", "", utils.NewError("生成令牌失败", 500)
	}

	refreshToken, err := s.generateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, "", "", utils.NewError("生成刷新令牌失败", 500)
	}

	return &user, accessToken, refreshToken, nil
}

func (s *AuthService) generateAccessToken(userID uint, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"type":    "access",
	})
	return token.SignedString(jwtSecret)
}

func (s *AuthService) generateRefreshToken(userID uint, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"type":    "refresh",
	})
	return token.SignedString(jwtSecret)
}

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", utils.NewError("刷新令牌无效", 401)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", utils.NewError("令牌解析失败", 401)
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", utils.NewError("令牌类型错误", 401)
	}

	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return "", "", utils.NewError("用户不存在", 401)
	}

	if user.Status == "banned" {
		return "", "", utils.NewError("该账号已被封禁", 403)
	}

	accessToken, err := s.generateAccessToken(userID, role)
	if err != nil {
		return "", "", utils.NewError("生成令牌失败", 500)
	}

	newRefreshToken, err := s.generateRefreshToken(userID, role)
	if err != nil {
		return "", "", utils.NewError("生成刷新令牌失败", 500)
	}

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) InitAdmin(username, qqNumber, displayName, password string) error {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return utils.NewError("系统已初始化", 400)
	}

	avatarURL := "https://q1.qlogo.cn/g?b=qq&nk=" + qqNumber + "&s=640"
	hashedPassword, err := database.HashPassword(password)
	if err != nil {
		return utils.NewError("密码加密失败", 500)
	}

	admin := models.User{
		Username:    username,
		QQNumber:    qqNumber,
		DisplayName: displayName,
		Password:    hashedPassword,
		Avatar:      avatarURL,
		Role:        "admin",
	}

	if result := database.DB.Create(&admin); result.Error != nil {
		return utils.NewError("初始化失败", 500)
	}

	categories := []models.Category{
		{Name: "表白墙", Description: "不止局限于表白墙", SortOrder: 1},
	}
	database.DB.Create(&categories)

	return nil
}

func (s *AuthService) CheckInit() bool {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	return count > 0
}

func (s *AuthService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return nil, utils.NewError("用户不存在", 404)
	}
	return &user, nil
}

func (s *AuthService) UpdateProfile(userID uint, displayName, avatar, signature string) error {
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.NewError("用户不存在", 404)
	}

	if displayName != "" {
		user.DisplayName = displayName
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if signature != "" {
		user.Signature = utils.SanitizeHTML(signature)
	}

	database.DB.Save(&user)
	return nil
}
