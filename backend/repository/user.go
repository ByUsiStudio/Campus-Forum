package repository

import (
	"forum/models"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	*BaseRepository[models.User]
}

// NewUserRepository 创建用户 Repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository[models.User](),
	}
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SearchUsers 搜索用户
func (r *UserRepository) SearchUsers(keyword string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

// UpdateOnlineStatus 更新在线状态
func (r *UserRepository) UpdateOnlineStatus(userID uint, onlineStatus string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("online_status", onlineStatus).Error
}

// UpdateLastActive 更新最后活跃时间
func (r *UserRepository) UpdateLastActive(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_active_at", gorm.Expr("NOW()")).Error
}

// GetOnlineUsers 获取在线用户
func (r *UserRepository) GetOnlineUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("online_status = ?", "online").Find(&users).Error
	return users, err
}

// GetUserWithPermissionGroups 获取用户及其权限组
func (r *UserRepository) GetUserWithPermissionGroups(userID uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("PermissionGroups").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
