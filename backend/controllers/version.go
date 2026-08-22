package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var (
	FrontendVersion = "unknown"
	BackendVersion  = "unknown"
)

func init() {
	if FrontendVersion != "unknown" && BackendVersion != "unknown" {
		return
	}

	version, ok := loadVersionFromFile()
	if !ok {
		return
	}

	if FrontendVersion == "unknown" && version.Frontend.Version != "" {
		FrontendVersion = version.Frontend.Version
	}
	if BackendVersion == "unknown" && version.Backend.Version != "" {
		BackendVersion = version.Backend.Version
	}
}

type fileVersion struct {
	Frontend struct {
		Version string `json:"version"`
	} `json:"frontend"`
	Backend struct {
		Version string `json:"version"`
	} `json:"backend"`
}

func loadVersionFromFile() (fileVersion, bool) {
	candidates := []string{
		"version.json",
		"../version.json",
		"build/version.json",
	}

	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "version.json"),
			filepath.Join(exeDir, "..", "version.json"),
		)
	}

	var v fileVersion
	for _, path := range candidates {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(content, &v); err != nil {
			continue
		}
		return v, true
	}
	return v, false
}

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"frontend": gin.H{"version": FrontendVersion},
		"backend":  gin.H{"version": BackendVersion},
	})
}
