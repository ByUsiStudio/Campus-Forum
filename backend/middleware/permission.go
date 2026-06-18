package middleware

import (
	"forum/database"
	"forum/models"
	"forum/repository"
	"forum/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var permissionGroupRepo = repository.NewPermissionGroupRepository()

// GetUserPermissionLevel 获取用户权限级别
func GetUserPermissionLevel(userID uint) int {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return 0
	}

	// system 角色级别为 100
	if user.Role == "system" {
		return 100
	}

	// admin 角色级别为 80
	if user.Role == "admin" {
		return 80
	}

	// 获取用户权限组的最高级别
	var userGroups []models.UserPermissionGroup
	database.DB.Where("user_id = ?", userID).Preload("PermissionGroup").Find(&userGroups)

	maxLevel := 0
	for _, ug := range userGroups {
		if ug.PermissionGroup.Level > maxLevel {
			maxLevel = ug.PermissionGroup.Level
		}
	}

	return maxLevel
}

// IsSystemAdmin 检查用户是否是系统管理员
func IsSystemAdmin(userID uint) bool {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return false
	}
	return user.Role == "system"
}

// IsAdmin 检查用户是否是管理员（admin 或 system）
func IsAdmin(userID uint) bool {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return false
	}
	return user.Role == "admin" || user.Role == "system"
}

// RequireSystemAdmin 要求系统管理员权限
func RequireSystemAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		if !IsSystemAdmin(userID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "需要系统管理员权限"})
			return
		}

		c.Next()
	}
}

// RequireMinLevel 要求最低权限级别
func RequireMinLevel(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		level := GetUserPermissionLevel(userID)
		if level < minLevel {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限级别不足"})
			return
		}

		c.Next()
	}
}

// RequirePermission 检查权限中间件
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		if userID == 0 {
			utils.Error("用户未登录")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		// 检查用户是否是超级管理员（system 角色拥有所有权限）
		var user models.User
		err := database.DB.First(&user, userID).Error
		if err == nil && user.Role == "system" {
			utils.Info("系统管理员 %s 跳过权限检查", user.Username)
			c.Next()
			return
		}

		// admin 角色需要检查具体权限
		if err == nil && user.Role == "admin" {
			// admin 只能管理普通用户内容，不能管理 system 相关内容
			hasPermission, checkErr := permissionGroupRepo.CheckUserPermission(userID, permission)
			if checkErr != nil {
				utils.Error("权限检查失败: %v", checkErr)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
				return
			}

			if !hasPermission {
				utils.Warn("管理员 %d 无权限: %s", userID, permission)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足"})
				return
			}

			utils.Info("管理员 %d 权限验证通过: %s", userID, permission)
			c.Next()
			return
		}

		// 检查用户是否有指定权限
		hasPermission, err := permissionGroupRepo.CheckUserPermission(userID, permission)
		if err != nil {
			utils.Error("权限检查失败: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
			return
		}

		if !hasPermission {
			utils.Warn("用户 %d 无权限: %s", userID, permission)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			return
		}

		utils.Info("用户 %d 权限验证通过: %s", userID, permission)
		c.Next()
	}
}

// RequireAnyPermission 检查是否有任一权限
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		// system 角色拥有所有权限
		var user models.User
		if err := database.DB.First(&user, userID).Error; err == nil && user.Role == "system" {
			c.Next()
			return
		}

		// 检查是否有任一权限
		for _, perm := range permissions {
			hasPermission, _ := permissionGroupRepo.CheckUserPermission(userID, perm)
			if hasPermission {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足"})
	}
}
