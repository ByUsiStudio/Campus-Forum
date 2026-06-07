package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 编译时注入的版本信息
var (
	FrontendVersion = "unknown"
	BackendVersion  = "unknown"
)

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"frontend": gin.H{"version": FrontendVersion},
		"backend":  gin.H{"version": BackendVersion},
	})
}
