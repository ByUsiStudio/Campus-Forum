package repository

import (
	"forum/database"
	"forum/models"
	"time"

	"gorm.io/gorm"
)

// PermissionGroupRepository 权限组数据访问层
type PermissionGroupRepository struct {
	*BaseRepository[models.PermissionGroup]
}

// NewPermissionGroupRepository 创建权限组 Repository
func NewPermissionGroupRepository() *PermissionGroupRepository {
	return &PermissionGroupRepository{
		BaseRepository: NewBaseRepository[models.PermissionGroup](),
	}
}

// GetAllGroups 获取所有权限组
func (r *PermissionGroupRepository) GetAllGroups() ([]models.PermissionGroup, error) {
	var groups []models.PermissionGroup
	err := r.db.Where("is_active = ?", true).Order("level DESC").Find(&groups).Error
	return groups, err
}

// GetGroupWithUsers 获取权限组及其用户
func (r *PermissionGroupRepository) GetGroupWithUsers(id uint) (*models.PermissionGroup, error) {
	var group models.PermissionGroup
	err := r.db.Preload("Users").First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetDefaultGroup 获取默认权限组
func (r *PermissionGroupRepository) GetDefaultGroup() (*models.PermissionGroup, error) {
	var group models.PermissionGroup
	err := r.db.Where("is_default = ? AND is_active = ?", true, true).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// SetDefaultGroup 设置默认权限组
func (r *PermissionGroupRepository) SetDefaultGroup(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.PermissionGroup{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.PermissionGroup{}).Where("id = ?", id).Update("is_default", true).Error; err != nil {
			return err
		}
		return nil
	})
}

// GrantPermissionToUser 给用户授予权限组
func (r *PermissionGroupRepository) GrantPermissionToUser(userID, groupID uint, expiresAt *time.Time) error {
	association := models.UserPermissionGroup{
		UserID:            userID,
		PermissionGroupID: groupID,
		ExpiresAt:         expiresAt,
	}
	return database.DB.Create(&association).Error
}

// RevokePermissionFromUser 从用户撤销权限组
func (r *PermissionGroupRepository) RevokePermissionFromUser(userID, groupID uint) error {
	return database.DB.Where("user_id = ? AND permission_group_id = ?", userID, groupID).Delete(&models.UserPermissionGroup{}).Error
}

// GetUserPermissionGroups 获取用户的所有权限组
func (r *PermissionGroupRepository) GetUserPermissionGroups(userID uint) ([]models.PermissionGroup, error) {
	var groups []models.PermissionGroup
	err := database.DB.Model(&models.PermissionGroup{}).
		Joins("JOIN user_permission_groups ON user_permission_groups.permission_group_id = permission_groups.id").
		Where("user_permission_groups.user_id = ?", userID).
		Find(&groups).Error
	return groups, err
}

// CheckUserPermission 检查用户是否有指定权限
func (r *PermissionGroupRepository) CheckUserPermission(userID uint, permission string) (bool, error) {
	var hasPermission bool
	err := database.DB.Raw(`
		SELECT COUNT(*) > 0 
		FROM permission_groups pg
		JOIN user_permission_groups upg ON upg.permission_group_id = pg.id
		WHERE upg.user_id = ? 
		AND pg.is_active = 1 
		AND (upg.expires_at IS NULL OR upg.expires_at > NOW())
		AND JSON_CONTAINS(pg.permissions, ?, '$')
	`, userID, "\""+permission+"\"").Scan(&hasPermission).Error
	return hasPermission, err
}

// GetUsersByGroup 获取权限组下的用户
func (r *PermissionGroupRepository) GetUsersByGroup(groupID uint, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{}).
		Joins("JOIN user_permission_groups ON user_permission_groups.user_id = users.id").
		Where("user_permission_groups.permission_group_id = ?", groupID)

	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("users.created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}
