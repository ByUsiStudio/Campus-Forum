package repository

import (
	"forum/models"
	"time"
)

// SystemLogRepository 系统日志数据访问层
type SystemLogRepository struct {
	*BaseRepository[models.SystemLog]
}

// NewSystemLogRepository 创建系统日志 Repository
func NewSystemLogRepository() *SystemLogRepository {
	return &SystemLogRepository{
		BaseRepository: NewBaseRepository[models.SystemLog](),
	}
}

// CreateLog 创建系统日志
func (r *SystemLogRepository) CreateLog(userID uint, action, module, details, ip string) error {
	log := models.SystemLog{
		UserID:  userID,
		Action:  action,
		Module:  module,
		Details: details,
		IP:      ip,
	}
	return r.Create(&log)
}

// GetLogsByUser 获取用户的操作日志
func (r *SystemLogRepository) GetLogsByUser(userID uint, page, pageSize int) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64

	query := r.db.Model(&models.SystemLog{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

// GetLogsByModule 按模块获取日志
func (r *SystemLogRepository) GetLogsByModule(module string, page, pageSize int) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64

	query := r.db.Model(&models.SystemLog{})
	if module != "" {
		query = query.Where("module = ?", module)
	}
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

// DeleteOldLogs 删除旧日志
func (r *SystemLogRepository) DeleteOldLogs(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", cutoff).Delete(&models.SystemLog{}).Error
}
