package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type VersionInfo struct {
	Frontend struct {
		Version string `json:"version"`
	} `json:"frontend"`
	Backend struct {
		Version string `json:"version"`
	} `json:"backend"`
	Swagger struct {
		Version string `json:"version"`
	} `json:"swagger"`
}

var versionCache *VersionInfo

func GetVersion(c *gin.Context) {
	if versionCache != nil {
		c.JSON(http.StatusOK, versionCache)
		return
	}

	file, err := os.Open("version.json")
	if err != nil {
		file, err = os.Open("../version.json")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"frontend": gin.H{"version": "unknown"},
				"backend":  gin.H{"version": "unknown"},
				"swagger":  gin.H{"version": "unknown"},
			})
			return
		}
	}
	defer file.Close()

	var version VersionInfo
	if err := json.NewDecoder(file).Decode(&version); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"frontend": gin.H{"version": "unknown"},
			"backend":  gin.H{"version": "unknown"},
			"swagger":  gin.H{"version": "unknown"},
		})
		return
	}

	versionCache = &version
	c.JSON(http.StatusOK, version)
}
