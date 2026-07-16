package controllers

import (
	"forum/service"
	"forum/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitJWT(secret string) {
	service.InitJWT(secret)
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

	err := service.Auth.Register(input.Username, input.QQNumber, input.DisplayName, input.Password)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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

	user, accessToken, refreshToken, err := service.Auth.Login(input.Username, input.Password)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         accessToken,
		"refresh_token": refreshToken,
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

func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := service.Auth.RefreshToken(input.RefreshToken)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"expires_in":    3600,
	})
}

func InitAdmin(c *gin.Context) {
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

	err := service.Auth.InitAdmin(input.Username, input.QQNumber, input.DisplayName, input.Password)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "初始化成功"})
}

func CheckInit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"initialized": service.Auth.CheckInit()})
}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := service.Auth.GetProfile(userID)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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

	var input struct {
		DisplayName string `json:"display_name"`
		Avatar      string `json:"avatar"`
		Signature   string `json:"signature"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.Auth.UpdateProfile(userID, input.DisplayName, input.Avatar, input.Signature)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
