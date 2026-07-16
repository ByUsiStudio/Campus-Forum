package controllers

import (
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
		Name        string `json:"name"`
		Logo        string `json:"logo"`
		Favicon     string `json:"favicon"`
		Description string `json:"description"`
		Keywords    string `json:"keywords"`
		ICP         string `json:"icp"`
		FooterText  string `json:"footer_text"`
		AllowSignup bool   `json:"allow_signup"`
		AllowPost   bool   `json:"allow_post"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.AdminConfig.UpdateSiteConfig(input.Name, input.Logo, input.Favicon, input.Description, input.Keywords, input.ICP, input.FooterText, input.AllowSignup, input.AllowPost)
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
