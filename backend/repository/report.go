package repository

import (
	"forum/models"
)

type ReportRepository struct {
	*BaseRepository[models.Report]
}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{
		BaseRepository: NewBaseRepository[models.Report](),
	}
}

func (r *ReportRepository) GetPendingReports(page, pageSize int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := r.db.Model(&models.Report{}).Where("status = ?", "pending")
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Reporter").Preload("Article").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reports).Error

	return reports, total, err
}

func (r *ReportRepository) GetReportsByUser(userID uint, page, pageSize int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := r.db.Model(&models.Report{}).Where("reporter_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Article").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reports).Error

	return reports, total, err
}