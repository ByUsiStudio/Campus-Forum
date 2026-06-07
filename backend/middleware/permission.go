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

// RequirePermission 检查权限中间件
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		if userID == 0 {
			utils.Error("用户未登录")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		// 检查用户是否是超级管理员
		var user models.User
		if err := database.DB.First(&user, userID).Error; err == nil && (user.Role == "admin" || user.Role == "system") {
			utils.Info("超级管理员 %s 跳过权限检查", user.Username)
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

		// 检查用户是否是超级管理员
		var user models.User
		if err := database.DB.First(&user, userID).Error; err == nil && (user.Role == "admin" || user.Role == "system") {
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
