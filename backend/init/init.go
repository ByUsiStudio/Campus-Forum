package init

import (
	"encoding/json"
	"forum/database"
	"forum/models"
	"forum/repository"
	"forum/utils"
)

// InitializeDefaultPermissionGroups 初始化默认权限组
func InitializeDefaultPermissionGroups() {
	utils.Info("检查并初始化默认权限组...")

	// 定义默认权限组
	defaultGroups := []struct {
		Name        string
		Description string
		Level       int
		IsDefault   bool
		Permissions []string
	}{
		{
			Name:        "新人",
			Description: "新注册用户的默认权限组",
			Level:       1,
			IsDefault:   true,
			Permissions: []string{
				"article.view",
				"comment.view",
				"comment.create",
				"user.profile",
			},
		},
		{
			Name:        "普通用户",
			Description: "通过审核的普通用户",
			Level:       10,
			IsDefault:   false,
			Permissions: []string{
				"article.view",
				"article.create",
				"article.edit.own",
				"article.delete.own",
				"comment.view",
				"comment.create",
				"comment.edit.own",
				"comment.delete.own",
				"user.profile",
				"user.follow",
				"like.create",
				"favorite.create",
			},
		},
		{
			Name:        "版主",
			Description: "版块版主，可管理相关内容",
			Level:       50,
			IsDefault:   false,
			Permissions: []string{
				"article.view",
				"article.create",
				"article.edit.all",
				"article.delete.all",
				"article.pin",
				"comment.view",
				"comment.create",
				"comment.edit.all",
				"comment.delete.all",
				"user.profile",
				"user.follow",
				"user.view.all",
				"like.create",
				"favorite.create",
				"report.view",
				"report.handle",
			},
		},
		{
			Name:        "内容审核员",
			Description: "负责内容审核和举报处理",
			Level:       60,
			IsDefault:   false,
			Permissions: []string{
				"article.view",
				"article.pin",
				"article.delete.all",
				"comment.view",
				"comment.delete.all",
				"user.view.all",
				"user.ban",
				"report.view",
				"report.handle",
				"deletion.view",
				"deletion.handle",
			},
		},
	}

	// 创建或更新权限组
	for _, groupData := range defaultGroups {
		var existing models.PermissionGroup
		err := database.DB.Where("name = ?", groupData.Name).First(&existing).Error

		permissionsJSON, _ := json.Marshal(groupData.Permissions)

		if err != nil {
			// 不存在，创建新的
			group := models.PermissionGroup{
				Name:        groupData.Name,
				Description: groupData.Description,
				Level:       groupData.Level,
				Permissions: string(permissionsJSON),
				IsDefault:   groupData.IsDefault,
				IsActive:    true,
			}
			if err := database.DB.Create(&group).Error; err != nil {
				utils.Error("创建权限组 %s 失败: %v", groupData.Name, err)
			} else {
				utils.Success("创建权限组: %s", groupData.Name)
			}
		} else {
			// 存在，更新
			updates := map[string]interface{}{
				"description": groupData.Description,
				"level":       groupData.Level,
				"permissions": string(permissionsJSON),
			}
			// 如果是默认组，处理默认标志
			if groupData.IsDefault {
				updates["is_default"] = true
			}
			if err := database.DB.Model(&existing).Updates(updates).Error; err != nil {
				utils.Error("更新权限组 %s 失败: %v", groupData.Name, err)
			} else {
				utils.Success("更新权限组: %s", groupData.Name)
			}
		}
	}

	// 确保只有一个默认组
	var defaultCount int64
	database.DB.Model(&models.PermissionGroup{}).Where("is_default = ? AND is_active = ?", true, true).Count(&defaultCount)
	if defaultCount == 0 {
		// 如果没有默认组，设置"新人"为默认
		database.DB.Model(&models.PermissionGroup{}).Where("name = ?", "新人").Update("is_default", true)
	} else if defaultCount > 1 {
		// 如果有多个默认组，只保留一个
		database.DB.Model(&models.PermissionGroup{}).
			Where("is_default = ? AND name != ?", true, "新人").
			Update("is_default", false)
	}

	utils.Success("权限组初始化完成")
}

// InitializeDefaultUserTitle 初始化默认用户头衔
func InitializeDefaultUserTitle() {
	utils.Info("检查并初始化默认用户头衔...")

	defaultTitles := []struct {
		Name        string
		Description string
		Color       string
	}{
		{
			Name:        "社区成员",
			Description: "普通社区成员",
			Color:       "#6b7280",
		},
		{
			Name:        "活跃用户",
			Description: "积极参与社区的用户",
			Color:       "#3b82f6",
		},
		{
			Name:        "资深用户",
			Description: "社区资深用户",
			Color:       "#8b5cf6",
		},
		{
			Name:        "荣誉版主",
			Description: "曾担任过版主的用户",
			Color:       "#f59e0b",
		},
	}

	for _, titleData := range defaultTitles {
		var existing models.Title
		err := database.DB.Where("name = ?", titleData.Name).First(&existing).Error

		if err != nil {
			title := models.Title{
				Name:        titleData.Name,
				Description: titleData.Description,
				Color:       titleData.Color,
			}
			if err := database.DB.Create(&title).Error; err != nil {
				utils.Error("创建头衔 %s 失败: %v", titleData.Name, err)
			} else {
				utils.Success("创建头衔: %s", titleData.Name)
			}
		}
	}

	utils.Success("用户头衔初始化完成")
}

// AssignDefaultGroupToUsers 为没有权限组的用户分配默认权限组
func AssignDefaultGroupToUsers() {
	utils.Info("检查并为用户分配默认权限组...")

	// 在函数内部创建 Repository，此时 database.DB 已初始化
	permissionGroupRepo := repository.NewPermissionGroupRepository()

	defaultGroup, err := permissionGroupRepo.GetDefaultGroup()
	if err != nil {
		utils.Error("获取默认权限组失败: %v", err)
		return
	}

	// 使用批量插入，避免逐个处理导致的性能问题
	result := database.DB.Exec(`
		INSERT INTO user_permission_groups (user_id, permission_group_id, created_at)
		SELECT u.id, ?, NOW()
		FROM users u
		LEFT JOIN user_permission_groups upg ON upg.user_id = u.id
		WHERE upg.id IS NULL
	`, defaultGroup.ID)

	if result.Error != nil {
		utils.Error("批量分配权限组失败: %v", result.Error)
	} else if result.RowsAffected > 0 {
		utils.Success("已为 %d 个用户分配默认权限组", result.RowsAffected)
	} else {
		utils.Info("所有用户都已分配权限组")
	}
}

// SystemInit 完整系统初始化
func SystemInit() {
	InitializeDefaultPermissionGroups()
	InitializeDefaultUserTitle()
	AssignDefaultGroupToUsers()
	utils.Success("系统初始化完成!")
}
