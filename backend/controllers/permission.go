package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetPermissionGroups 获取所有权限组
func GetPermissionGroups(c *gin.Context) {
	var groups []models.PermissionGroup
	database.DB.Where("is_active = ?", true).Order("level DESC, created_at ASC").Find(&groups)

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// GetPermissionGroup 获取单个权限组
func GetPermissionGroup(c *gin.Context) {
	groupID := c.Param("id")

	var group models.PermissionGroup
	if err := database.DB.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "权限组不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": group})
}

// CreatePermissionGroup 创建权限组
func CreatePermissionGroup(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Level       int      `json:"level"`
		Permissions []string `json:"permissions" binding:"required"`
		IsDefault   bool     `json:"is_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 检查名称是否已存在
	var existing models.PermissionGroup
	if err := database.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "权限组名称已存在"})
		return
	}

	// 如果设置为默认组，取消其他默认组
	if req.IsDefault {
		database.DB.Model(&models.PermissionGroup{}).Where("is_default = ?", true).Update("is_default", false)
	}

	// 将权限列表转换为JSON
	permissionsJSON, err := json.Marshal(req.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "处理权限列表失败"})
		return
	}

	level := req.Level
	if level == 0 {
		level = 1
	}

	group := models.PermissionGroup{
		Name:        req.Name,
		Description: req.Description,
		Level:       level,
		Permissions: string(permissionsJSON),
		IsDefault:   req.IsDefault,
		IsActive:    true,
	}

	if err := database.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建权限组失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "group": group})
}

// UpdatePermissionGroup 更新权限组
func UpdatePermissionGroup(c *gin.Context) {
	groupID := c.Param("id")

	var group models.PermissionGroup
	if err := database.DB.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "权限组不存在"})
		return
	}

	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Level       int      `json:"level"`
		Permissions []string `json:"permissions"`
		IsDefault   *bool    `json:"is_default"`
		IsActive    *bool    `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 检查名称是否与其他组冲突
	if req.Name != "" && req.Name != group.Name {
		var existing models.PermissionGroup
		if err := database.DB.Where("name = ? AND id != ?", req.Name, groupID).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "权限组名称已存在"})
			return
		}
	}

	// 如果设置为默认组，取消其他默认组
	if req.IsDefault != nil && *req.IsDefault {
		database.DB.Model(&models.PermissionGroup{}).Where("id != ? AND is_default = ?", groupID, true).Update("is_default", false)
	}

	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Level > 0 {
		updates["level"] = req.Level
	}
	if req.Permissions != nil {
		permissionsJSON, _ := json.Marshal(req.Permissions)
		updates["permissions"] = string(permissionsJSON)
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		database.DB.Model(&group).Updates(updates)
	}

	database.DB.First(&group, groupID)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "group": group})
}

// DeletePermissionGroup 删除权限组
func DeletePermissionGroup(c *gin.Context) {
	groupID := c.Param("id")

	var group models.PermissionGroup
	if err := database.DB.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "权限组不存在"})
		return
	}

	// 检查是否有用户使用此权限组
	var count int64
	database.DB.Model(&models.UserPermissionGroup{}).Where("permission_group_id = ?", groupID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该权限组下有用户，无法删除"})
		return
	}

	database.DB.Delete(&group)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GrantPermissionGroup 授予用户权限组
func GrantPermissionGroup(c *gin.Context) {
	granterID := c.GetUint("user_id")

	var req struct {
		UserID        uint `json:"user_id" binding:"required"`
		GroupID       uint `json:"group_id" binding:"required"`
		ExpiresInDays int  `json:"expires_in_days"` // 0表示永久
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 验证用户和权限组是否存在
	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var group models.PermissionGroup
	if err := database.DB.First(&group, req.GroupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "权限组不存在"})
		return
	}

	// 检查用户是否已有该权限组
	var existing models.UserPermissionGroup
	if err := database.DB.Where("user_id = ? AND permission_group_id = ?", req.UserID, req.GroupID).First(&existing).Error; err == nil {
		// 已存在，更新过期时间
		var expiresAt *time.Time
		if req.ExpiresInDays > 0 {
			exp := time.Now().AddDate(0, 0, req.ExpiresInDays)
			expiresAt = &exp
		}
		database.DB.Model(&existing).Updates(map[string]interface{}{
			"expires_at": expiresAt,
			"granted_by": granterID,
		})
		c.JSON(http.StatusOK, gin.H{"message": "权限组已更新", "user_group": existing})
		return
	}

	// 创建新的关联
	var expiresAt *time.Time
	if req.ExpiresInDays > 0 {
		exp := time.Now().AddDate(0, 0, req.ExpiresInDays)
		expiresAt = &exp
	}

	userGroup := models.UserPermissionGroup{
		UserID:            req.UserID,
		PermissionGroupID: req.GroupID,
		GrantedBy:         granterID,
		ExpiresAt:         expiresAt,
	}

	if err := database.DB.Create(&userGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "授予权限组失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限组授予成功", "user_group": userGroup})
}

// RevokePermissionGroup 撤销用户权限组
func RevokePermissionGroup(c *gin.Context) {
	userID := c.Param("user_id")
	groupID := c.Param("id")

	var userGroup models.UserPermissionGroup
	if err := database.DB.Where("user_id = ? AND permission_group_id = ?", userID, groupID).First(&userGroup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户权限组关联不存在"})
		return
	}

	database.DB.Delete(&userGroup)
	c.JSON(http.StatusOK, gin.H{"message": "权限组已撤销"})
}

// GetUserPermissionGroups 获取用户的所有权限组
func GetUserPermissionGroups(c *gin.Context) {
	userID := c.Param("user_id")

	var userGroups []models.UserPermissionGroup
	database.DB.Where("user_id = ?", userID).
		Preload("PermissionGroup").
		Find(&userGroups)

	// 过滤掉已过期的权限组
	var activeGroups []models.UserPermissionGroup
	now := time.Now()
	for _, ug := range userGroups {
		if ug.ExpiresAt == nil || ug.ExpiresAt.After(now) {
			activeGroups = append(activeGroups, ug)
		}
	}

	c.JSON(http.StatusOK, gin.H{"groups": activeGroups})
}

// CheckUserPermissions 检查用户权限
func CheckUserPermissions(c *gin.Context) {
	userID := c.GetUint("user_id")
	requiredPermissions := c.QueryArray("permissions")

	// 获取用户角色
	role := c.GetString("role")
	isAdmin := role == "admin" || role == "system"

	// 管理员拥有所有权限
	if isAdmin {
		c.JSON(http.StatusOK, gin.H{
			"has_permission": true,
			"is_admin":       true,
		})
		return
	}

	// 获取用户的所有权限组
	var userGroups []models.UserPermissionGroup
	database.DB.Where("user_id = ?", userID).
		Preload("PermissionGroup").
		Find(&userGroups)

	var allPermissions []string
	now := time.Now()

	for _, ug := range userGroups {
		// 检查是否过期
		if ug.ExpiresAt != nil && ug.ExpiresAt.Before(now) {
			continue
		}

		var perms []string
		if err := json.Unmarshal([]byte(ug.PermissionGroup.Permissions), &perms); err == nil {
			allPermissions = append(allPermissions, perms...)
		}
	}

	// 检查是否有所需权限
	hasAll := true
	for _, required := range requiredPermissions {
		found := false
		for _, perm := range allPermissions {
			if perm == required || perm == "*" {
				found = true
				break
			}
		}
		if !found {
			hasAll = false
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"has_permission": hasAll,
		"permissions":    allPermissions,
		"is_admin":       false,
	})
}

// InitializeDefaultPermissionGroups 初始化默认权限组
func InitializeDefaultPermissionGroups(c *gin.Context) {
	// 检查是否已有权限组
	var count int64
	database.DB.Model(&models.PermissionGroup{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "权限组已存在，无需初始化"})
		return
	}

	defaultGroups := []models.PermissionGroup{
		{
			Name:        "新人",
			Description: "新注册用户默认权限组",
			Level:       1,
			IsDefault:   true,
			IsActive:    true,
			Permissions: `["article:read", "comment:create", "user:profile:view"]`,
		},
		{
			Name:        "普通用户",
			Description: "普通用户权限组",
			Level:       10,
			IsDefault:   false,
			IsActive:    true,
			Permissions: `["article:read", "article:create", "article:edit:own", "article:delete:own", "comment:create", "comment:edit:own", "comment:delete:own", "user:profile:view", "user:profile:edit"]`,
		},
		{
			Name:        "版主",
			Description: "版主权限组，可管理板块内容",
			Level:       50,
			IsDefault:   false,
			IsActive:    true,
			Permissions: `["article:read", "article:create", "article:edit:own", "article:delete:own", "article:pin", "article:delete:any", "comment:create", "comment:edit:own", "comment:delete:own", "comment:delete:any", "user:profile:view", "user:profile:edit", "category:manage"]`,
		},
		{
			Name:        "内容审核员",
			Description: "内容审核员权限组",
			Level:       60,
			IsDefault:   false,
			IsActive:    true,
			Permissions: `["article:read", "article:edit:any", "article:delete:any", "article:pin", "article:featured", "comment:read", "comment:delete:any", "report:view", "report:handle", "user:profile:view", "user:ban"]`,
		},
		{
			Name:        "admin",
			Description: "管理员权限组，只能管理普通用户内容",
			Level:       80,
			IsDefault:   false,
			IsActive:    true,
			Permissions: `["article:read", "article:create", "article:edit:own", "article:delete:own", "article:delete:user", "comment:create", "comment:edit:own", "comment:delete:own", "comment:delete:user", "user:profile:view", "user:profile:edit", "user:mark_delete", "announcement:update", "notification:manage", "report:view", "report:handle"]`,
		},
		{
			Name:        "system",
			Description: "系统管理员权限组，拥有所有权限",
			Level:       100,
			IsDefault:   false,
			IsActive:    true,
			Permissions: `["*"]`,
		},
	}

	for _, group := range defaultGroups {
		if err := database.DB.Create(&group).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "初始化权限组失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "默认权限组初始化成功", "groups": defaultGroups})
}
