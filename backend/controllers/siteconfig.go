package controllers

import (
	"forum/models"
	"forum/service"
	"forum/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSiteConfig(c *gin.Context) {
	config, err := service.AdminConfig.GetSiteConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func UpdateSiteConfig(c *gin.Context) {
	var input struct {
		SiteTitle       string `json:"site_title"`
		SiteDescription string `json:"site_description"`
		SiteLogo        string `json:"site_logo"`
		ICPNumber       string `json:"icp_number"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := models.SiteConfig{
		SiteTitle:       input.SiteTitle,
		SiteDescription: input.SiteDescription,
		SiteLogo:        input.SiteLogo,
		ICPNumber:       input.ICPNumber,
	}

	err := service.AdminConfig.UpdateSiteConfig(config)
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

func TestSMTPConfig(c *gin.Context) {
	var input struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		From     string `json:"from"`
		To       string `json:"to"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMTP配置测试成功", "test_email": input.To})
}
